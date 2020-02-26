package auth

import (
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"

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
	authCacheMtx.RLock()
	auth, ok := authCache[accessKey]
	authCacheMtx.RUnlock()
	if !ok {
		res := authCol.FindOne(nil, bson.M{"_id": accessKey}, options.FindOne())
		if err := res.Decode(&auth); err != nil {
			log.Debug("Error On GetAccessKey",
				zap.Error(err),
				zap.String("AccessKey", accessKey),
			)
			msg.WriteError(ctx, http.StatusForbidden, msg.ErrAccessTokenInvalid)
			return
		}
	}
	if auth.ExpiredOn > 0 && time.Now().Unix() > auth.ExpiredOn {
		msg.WriteError(ctx, http.StatusForbidden, msg.ErrAccessTokenExpired)
		return
	}

	authCacheMtx.Lock()
	authCache[accessKey] = auth
	authCacheMtx.Unlock()

	ctx.Values().Save(CtxAuth, auth, true)
	ctx.Values().Save(CtxClientName, auth.AppName, true)
	ctx.Next()
}

func MustAdmin(ctx iris.Context) {
	if !hasAdminAccess(ctx) {
		msg.WriteError(ctx, http.StatusForbidden, msg.ErrNoPermission)
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
		msg.WriteError(ctx, http.StatusForbidden, msg.ErrNoPermission)
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
		msg.WriteError(ctx, http.StatusForbidden, msg.ErrNoPermission)
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
	accessToken := tools.RandomID(64)

	req := &CreateAccessToken{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	authPerms := make([]Permission, 0, 3)

	for _, p := range req.Permissions {
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
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPermissionIsNotSet)
		return
	}

	authCreatedOn := time.Now().Unix()
	authExpireOn := int64(0)
	if req.Period > 0 {
		authExpireOn = authCreatedOn + req.Period*86400
	}

	_, err = authCol.InsertOne(nil, Auth{
		ID:          accessToken,
		Permissions: authPerms,
		CreatedOn:   authCreatedOn,
		ExpiredOn:   authExpireOn,
		AppName:     req.AppName,
	})
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, CAccessTokenCreated, AccessTokenCreated{
		AccessToken: accessToken,
		ExpireOn:    authExpireOn,
	})
}

// LoginHandler is API handler
// API: /auth/send_code
// Http Method: POST
// Inputs: JSON
//	phone: string
// Returns: PhoneCodeSent (PHONE_CODE_SENT)
// Possible Errors:
//	1. 500: READ_FROM_CACHE
//	2. 400: CANNOT_UNMARSHAL_JSON
//	2. 400: PHONE_NOT_VALID
func SendCodeHandler(ctx iris.Context) {
	req := &SendCodeReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	if config.GetBool(config.TestMode) && strings.HasPrefix(req.Phone, config.GetString(config.MagicPhone)) {
		sendCodeMagicNumber(ctx, req.Phone)
		return
	}

	if len(req.Phone) < 5 {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneNotValid)
		return
	}

	v, err := redisCache.GetString(fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone))
	if err != nil {
		log.Warn("Error On ReadFromCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromCache)
		return
	}
	if v != "" {
		u, _ := user.GetByPhone(req.Phone)
		verifyParams := strings.Split(v, "|")
		// TODO:: fix this
		_, _ = smsProvider.SendInBackground(req.Phone, fmt.Sprintf("MusicChi Code: %s", verifyParams[2]))
		msg.WriteResponse(ctx, CPhoneCodeSent, PhoneCodeSent{
			PhoneCodeHash: verifyParams[0],
			Registered:    u != nil,
		})
		return
	}

	switch ctx.Values().GetString(CtxClientName) {
	case AppNameMusicChi:
		sendMusicChi(ctx, req.Phone)
	default:
		sendCode(ctx, req.Phone)
	}

}
func sendCode(ctx iris.Context, phone string) {
	phoneCodeHash := tools.RandomID(12)
	phoneCode := tools.RandomDigit(4)
	if config.GetBool(config.TestMode) {
		phoneCode = "2374"
	}
	u, _ := user.GetByPhone(phone)
	err := redisCache.Do(radix.FlatCmd(nil, "SETEX",
		fmt.Sprintf("%s.%s", config.RkPhoneCode, phone),
		600,
		fmt.Sprintf("%s|%s|%s", phoneCodeHash, "", phoneCode),
	))
	if err != nil {
		log.Warn("Error On WriteToCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToCache)
		return
	}

	msg.WriteResponse(ctx, CPhoneCodeSent, PhoneCodeSent{
		PhoneCodeHash: phoneCodeHash,
		Registered:    u != nil,
	})
}
func sendMusicChi(ctx iris.Context, phone string) {
	var phoneCodeHash, otpID, phoneCode string
	phoneCodeHash = tools.RandomID(12)
	u, _ := user.GetByPhone(phone)
	if u != nil && u.VasPaid {
		phoneCode = tools.RandomDigit(4)
		// User our internal sms provider
		if ce := log.Check(log.DebugLevel, "Send Code (MusicChi)"); ce != nil {
			ce.Write(
				zap.String("Phone", phone),
				zap.String("PhoneCode", phoneCode),
			)
		}
		_, err := smsProvider.SendInBackground(phone, fmt.Sprintf("MusicChi Code: %s", phoneCode))
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrNoResponseFromSmsServer)
			return
		}
	} else {
		if _, ok := supportedCarriers[phone[:5]]; !ok {
			msg.WriteError(ctx, http.StatusNotAcceptable, msg.ErrUnsupportedCarrier)
			return
		}
		res, err := saba.Subscribe(phone)
		if err != nil {
			log.Warn("Error On Saba Subscribe", zap.Error(err))
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Err3rdParty)
			return
		}
		otpID = res.OtpID
		switch res.StatusCode {
		case "SC111", "SC000":
		default:
			// If we are here, then it means VAS did not send the sms
			msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrNoResponseFromVAS)
		}
	}

	err := redisCache.Do(radix.FlatCmd(nil, "SETEX",
		fmt.Sprintf("%s.%s", config.RkPhoneCode, phone),
		600,
		fmt.Sprintf("%s|%s|%s", phoneCodeHash, otpID, phoneCode),
	))
	if err != nil {
		log.Warn("Error On WriteToCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToCache)
		return
	}

	msg.WriteResponse(ctx, CPhoneCodeSent, PhoneCodeSent{
		PhoneCodeHash: phoneCodeHash,
		Registered:    u != nil,
	})
}
func sendCodeMagicNumber(ctx iris.Context, magicPhone string) {
	phoneCode := config.GetString(config.MagicPhoneCode)
	phoneCodeHash := tools.RandomID(12)
	err := redisCache.Do(radix.FlatCmd(nil, "SETEX",
		fmt.Sprintf("%s.%s", config.RkPhoneCode, magicPhone),
		600,
		fmt.Sprintf("%s|%s|%s", phoneCodeHash, "", phoneCode),
	))
	if err != nil {
		log.Warn("Error On WriteToCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToCache)
		return
	}
	msg.WriteResponse(ctx, CPhoneCodeSent, PhoneCodeSent{
		PhoneCodeHash: phoneCodeHash,
		Registered:    true,
	})
}

// LoginHandler is API handler
// API: /auth/login
// Http Method: POST
// Inputs: JSON
//	phone_code: string
//	phone_code_hash: string
//	phone: string
// Returns: Authorization (AUTHORIZATION)
// Possible Errors:
//	1. 500: with error text
//	2. 400: PHONE_NOT_VALID
//	3. 400: PHONE_CODE_NOT_VALID
//	4. 400: PHONE_CODE_HASH_NOT_VALID
func LoginHandler(ctx iris.Context) {
	req := &LoginReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	var otpID, phoneCode, phoneCodeHash string
	if v, err := redisCache.GetString(fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone)); err != nil {
		log.Warn("Error On ReadFromCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromCache)
		return
	} else {
		verifyParams := strings.Split(v, "|")
		phoneCodeHash = verifyParams[0]
		otpID = verifyParams[1]
		phoneCode = verifyParams[2]
	}
	if req.PhoneCodeHash != phoneCodeHash {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneCodeHashNotValid)
		return
	}
	if otpID != "" {
		vasCode, err := saba.Confirm(req.Phone, req.PhoneCode, otpID)
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
			return
		}
		if vasCode != saba.SuccessfulCode {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(saba.Codes[vasCode]))
			return
		}
	} else if req.PhoneCode != phoneCode {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneCodeNotValid)
		return
	}

	appName := ctx.Values().GetString(CtxClientName)
	u, err := user.GetByPhone(req.Phone)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneNotValid)
		return
	}
	err = session.Remove(u.ID, appName)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromDb)
		return
	}
	sessionID := tools.RandomID(64)
	timeNow := time.Now().Unix()
	err = session.Save(&session.Session{
		ID:         sessionID,
		UserID:     u.ID,
		CreatedOn:  timeNow,
		LastAccess: timeNow,
		App:        appName,
	})
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	_ = redisCache.Del(fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone))
	msg.WriteResponse(ctx, CAuthorization, Authorization{
		UserID:    u.ID,
		Phone:     u.Phone,
		Username:  u.Username,
		SessionID: sessionID,
	})

}

// RegisterHandler is API handler
// API: /auth/register
// Http Method: POST
// Inputs: JSON
//	phone_code: string
//	phone_code_hash: string
//	phone: string
//	username: string
// Returns: Authorization (AUTHORIZATION)
// Possible Errors:
//	1. 500: with error text
//	2. 400: PHONE_NOT_VALID
//	3. 400: PHONE_CODE_NOT_VALID
//	4. 400: PHONE_CODE_HASH_NOT_VALID
func RegisterHandler(ctx iris.Context) {
	req := &RegisterReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	var otpID, phoneCode, phoneCodeHash string
	if v, err := redisCache.GetString(fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone)); err != nil {
		log.Warn("Error On ReadFromCache", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromCache)
		return
	} else {
		verifyParams := strings.Split(v, "|")
		if len(verifyParams) != 3 {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrCorruptData)
			return
		}
		phoneCodeHash = verifyParams[0]
		otpID = verifyParams[1]
		phoneCode = verifyParams[2]
	}

	if req.PhoneCodeHash != phoneCodeHash {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneCodeHashNotValid)
		return
	}
	if otpID != "" {
		vasCode, err := saba.Confirm(req.Phone, req.PhoneCode, otpID)
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
			return
		}
		if vasCode != saba.SuccessfulCode {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(saba.Codes[vasCode]))
			return
		}
	} else if req.PhoneCode != phoneCode {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneCodeNotValid)
		return
	}

	_, err = user.GetByPhone(req.Phone)
	if err == nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrAlreadyRegistered)
		return
	}

	if req.Username == "" {
		req.Username = fmt.Sprintf("USER%s", strings.ToUpper(tools.RandomID(12)))
	} else if !usernameREGX.Match(tools.StrToByte(req.Username)) {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrUsernameFormat)
		return
	}

	userID := fmt.Sprintf("U%s", tools.RandomID(32))
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
		u, _ := user.GetByPhone(req.Phone)
		if u == nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
			return
		}
	}

	sessionID := tools.RandomID(64)
	err = session.Save(&session.Session{
		ID:         sessionID,
		UserID:     userID,
		CreatedOn:  timeNow,
		LastAccess: timeNow,
		App:        ctx.Values().GetString(CtxClientName),
	})
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	_ = redisCache.Del(fmt.Sprintf("%s.%s", config.RkPhoneCode, req.Phone))
	msg.WriteResponse(ctx, CAuthorization, Authorization{
		UserID:    userID,
		Phone:     req.Phone,
		Username:  req.Username,
		SessionID: sessionID,
	})

}

func LogoutHandler(ctx iris.Context) {
	req := &LogoutReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	s, ok := ctx.Values().Get(session.CtxSession).(session.Session)
	if !ok {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrSessionInvalid)
		return
	}
	if req.Unsubscribe {
		u, err := user.Get(s.UserID)
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrUserNotFound)
			return
		}
		_, err = saba.Unsubscribe(u.Phone)
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrNoResponseFromVAS)
			return
		}
	}
	err = session.Remove(s.UserID, ctx.Values().GetString(CtxClientName))
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{
		Success: true,
	})
}
