package token

import (
	"encoding/json"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

/*
   Creation Time: 2019 - Sep - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func CreateHandler(ctx iris.Context) {
	phone := ctx.PostValue("Phone")
	period := ctx.PostValueInt64Default("Period", 0) // Period is the number of days

	tokenStartDate := time.Now().Unix()
	tokenExpireDate := time.Now().Unix() + period*84600
	token := ronak.RandomID(64)
	res := make([][]string, 0, 1)
	_, err := tokenCol.InsertOne(nil, Token{
		ID:        token,
		Phone:     phone,
		Period:    period,
		CreatedOn: tokenStartDate,
		ExpiredOn: tokenExpireDate,
	})
	if err != nil {
		res = append(res, []string{"", "Error"})
	} else {
		res = append(res, []string{token, "Created."})
	}

	resBytes, _ := json.Marshal(res)
	_, _ = ctx.Write(resBytes)
}

func ValidateHandler(ctx iris.Context) {
	deviceID := ctx.PostValue("DeviceID")
	tokenID := ctx.PostValue("Token")

	mtxLock.RLock()
	token, ok := tokenCache[tokenID]
	mtxLock.RUnlock()
	if !ok {
		res := tokenCol.FindOne(nil, bson.M{"_id": tokenID}, options.FindOne())
		if err := res.Decode(&token); err != nil {
			msg.Error(ctx, http.StatusForbidden, msg.ErrAccessTokenInvalid)
			return
		}
		mtxLock.Lock()
		tokenCache[tokenID] = token
		mtxLock.Unlock()
	}

	if time.Now().Unix() > token.ExpiredOn {
		msg.Error(ctx, http.StatusForbidden, msg.ErrAccessTokenExpired)
		return
	}
	remainingDays := int64(time.Now().Sub(time.Unix(token.ExpiredOn, 0)).Hours() / 24)

	if token.DeviceID == "" {
		_, err := tokenCol.UpdateOne(nil, bson.M{"_id": tokenID}, bson.M{"$set": bson.M{"device_id": deviceID}})
		if err != nil {
			msg.Error(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
			return
		}
		token.DeviceID = deviceID
		mtxLock.Lock()
		tokenCache[tokenID] = token
		mtxLock.Unlock()
	}

	msg.WriteResponse(ctx, CValidated, Validated{
		RemainingDays: remainingDays,
	})
}
