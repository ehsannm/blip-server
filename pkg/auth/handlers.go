package auth

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/kataras/iris"
	"github.com/mediocregopher/radix/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func MustHaveAccessKey(ctx iris.Context) {
	accessKey := ctx.GetHeader(HdrAccessKey)
	mtxLock.RLock()
	auth, ok := authCache[accessKey]
	mtxLock.RUnlock()
	if !ok {
		res := authCol.FindOne(nil, bson.M{"_id": accessKey}, options.FindOne())
		if err := res.Decode(&auth); err != nil {
			log.Debug("Error On GetAuthorization",
				zap.Error(err),
				zap.String("AccessKey", accessKey),
			)
			msg.Error(ctx, http.StatusForbidden, msg.ErrAccessTokenInvalid)
			return
		}
	}
	if auth.ExpiredOn > 0 && time.Now().Unix() > auth.ExpiredOn {
		msg.Error(ctx, http.StatusForbidden, msg.ErrAccessTokenExpired)
		return
	}

	mtxLock.Lock()
	authCache[accessKey] = auth
	mtxLock.Unlock()

	ctx.Values().Save(CtxAuth, auth, true)
	ctx.Next()
}

func MustAdmin(ctx iris.Context) {
	if !hasAdminAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}
func hasAdminAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Admin {
			return true
		}
	}
	return false
}

func MustWriteAccess(ctx iris.Context) {
	if !hasWriteAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}
func hasWriteAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Write || p == Admin {
			return true
		}
	}
	return false
}

func MustReadAccess(ctx iris.Context) {
	if !hasReadAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}
func hasReadAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Read || p == Admin {
			return true
		}
	}
	return false
}

func CreateAccessKeyHandler(ctx iris.Context) {
	accessToken := ronak.RandomID(64)
	authPerms := make([]Permission, 0, 3)
	perms := ctx.PostValues("Permissions")
	period := ctx.PostValueInt64Default("Period", 0)

	for _, p := range perms {
		switch strings.ToLower(p) {
		case "admin":
			authPerms = append(authPerms, Admin)
		case "read":
			authPerms = append(authPerms, Read)
		case "write":
			authPerms = append(authPerms, Write)
		}
	}

	if len(authPerms) == 0 {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrPermissionIsNotSet)
		return
	}

	authCreatedOn := time.Now().Unix()
	authExpireOn := int64(0)
	if period > 0 {
		authExpireOn = authCreatedOn + period*86400
	}

	_, err := authCol.InsertOne(nil, Auth{
		ID:          accessToken,
		Permissions: authPerms,
		CreatedOn:   authCreatedOn,
		ExpiredOn:   authExpireOn,
	})
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, CAccessTokenCreated, AccessTokenCreated{
		AccessToken: accessToken,
		ExpireOn:    authExpireOn,
	})
}

func SendCodeHandler(ctx iris.Context) {
	req := new(SendCodeReq)
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	if _, ok := supportedCarriers[req.Phone[:5]]; !ok {
		msg.Error(ctx, http.StatusNotAcceptable, msg.ErrUnsupportedCarrier)
		return
	}

	optID, err := saba.Subscribe(req.Phone)
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	if optID == "" {
		// If we are here, then it means VAS did not send the sms
	}

	err = redisCache.Do(radix.FlatCmd(nil, "SET", fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone), optID))
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.ErrWriteToCache)
		return
	}

	registered := false
	u, _ := user.GetByPhone(req.Phone)
	if u != nil {
		registered = true
	}

	msg.WriteResponse(ctx, CPhoneCodeSent, PhoneCodeSent{
		PhoneCodeHash: ronak.RandomID(12),
		OperationID:   optID,
		Registered:    registered,
	})
}

func LoginHandler(ctx iris.Context) {
	req := new(LoginReq)
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	vasCode, err := saba.Confirm(req.Phone, req.PhoneCode, req.OperationID)
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	if vasCode != saba.SuccessfulCode {
		errText, _ := saba.Codes[vasCode]
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(errText))
		return
	}

	u, err := user.GetByPhone(req.Phone)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrPhoneNotValid)
		return
	}

	sessionID := ronak.RandomID(64)
	timeNow := time.Now().Unix()
	err = session.Save(session.Session{
		ID:         sessionID,
		UserID:     u.ID,
		CreatedOn:  timeNow,
		LastAccess: timeNow,
	})
	if err != nil {
		errText, _ := saba.Codes[vasCode]
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(errText))
		return
	}

	msg.WriteResponse(ctx, CAuthorization, Authorization{
		UserID:    u.ID,
		Phone:     u.Phone,
		Username:  u.Username,
		SessionID: sessionID,
	})

}

func RegisterHandler(ctx iris.Context) {
	req := new(RegisterReq)
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	vasCode, err := saba.Confirm(req.Phone, req.PhoneCode, req.OperationID)
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	if vasCode != saba.SuccessfulCode {
		errText, _ := saba.Codes[vasCode]
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(errText))
		return
	}

	_, err = user.GetByPhone(req.Phone)
	if err == nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrAlreadyRegistered)
		return
	}

	if req.Username == "" {
		req.Username = fmt.Sprintf("USER%s", strings.ToUpper(ronak.RandomID(12)))
	} else if !usernameREGX.Match(ronak.StrToByte(req.Username)) {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrUsernameFormat)
		return
	}

	userID := fmt.Sprintf("U%s", ronak.RandomID(32))
	timeNow := time.Now().Unix()
	err = user.Save(&user.User{
		ID:        userID,
		Username:  req.Username,
		Phone:     req.Phone,
		Email:     "",
		CreatedOn: timeNow,
		Disabled:  false,
	})
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	sessionID := ronak.RandomID(64)
	err = session.Save(session.Session{
		ID:         sessionID,
		UserID:     userID,
		CreatedOn:  timeNow,
		LastAccess: timeNow,
	})
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	msg.WriteResponse(ctx, CAuthorization, Authorization{
		UserID:    userID,
		Phone:     req.Phone,
		Username:  req.Username,
		SessionID: sessionID,
	})

}
