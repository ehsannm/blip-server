package auth

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetAuthorization(ctx iris.Context) {
	accessKey := ctx.GetHeader(HdrAccessKey)
	mtxLock.RLock()
	auth, ok := authCache[accessKey]
	mtxLock.RUnlock()
	if !ok {
		res := authCol.FindOne(nil, bson.M{"_id": accessKey}, options.FindOne())
		if err := res.Decode(&auth); err != nil {
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
	if !HasAdminAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}

func MustWriteAccess(ctx iris.Context) {
	if !HasWriteAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}

func MustReadAccess(ctx iris.Context) {
	if !HasReadAccess(ctx) {
		msg.Error(ctx, http.StatusForbidden, msg.ErrNoPermission)
		return
	}
	ctx.Next()
}

func CreateAccessKey(ctx iris.Context) {
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
