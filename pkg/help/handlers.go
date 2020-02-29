package help

import (
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"github.com/kataras/iris"
	"github.com/rogpeppe/go-internal/semver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

/*
   Creation Time: 2020 - Feb - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// SetHandler is API Handler
// Http Method: POST  /help/config
// Inputs: POST VALUES:
// Returns: Bool (BOOL)
// Possible Errors:
//	1. 400: CANNOT_UNMARSHAL_JSON
//	2. 500: WRITE_TO_DB
func SetHandler(ctx iris.Context) {
	req := &SetDefaultConfig{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	_, err = helpCol.UpdateOne(nil,
		bson.M{"_id": "defaults"},
		bson.M{"$set": bson.M{
			req.Key: req.Value,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	// reload default configs from the server
	// FIXME:: this is not multi server
	loadDefaultConfig()

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}

// UnsetHandler is API Handler
// Http Method: DELETE /help/config
// Inputs:
// Possible Errors:
//	1.
func UnsetHandler(ctx iris.Context) {
	req := &UnsetDefaultConfig{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	_, err = helpCol.UpdateOne(nil,
		bson.M{"_id": "defaults"},
		bson.M{"$unset": bson.M{
			req.Key: "",
		}},
		options.Update(),
	)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	// reload default configs from the server
	// FIXME:: this is not multi server
	loadDefaultConfig()

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}

// SetHandler is API Handler
// Http Method: GET  /help/config
// Returns: Config (CONFIG)
// Possible Errors:
//	1. 400: CANNOT_UNMARSHAL_JSON
//	2. 500: WRITE_TO_DB
func GetHandler(ctx iris.Context) {
	clientAppVer := ctx.GetHeader(HdrAppVersion)
	clientPlatform := strings.ToLower(ctx.GetHeader(HdrPlatform))
	currAppVersion := getConfig(fmt.Sprintf("%s_%s_cur",
		ctx.Values().GetString(auth.CtxClientName),
		clientPlatform,
	))
	minAppVersion := getConfig(fmt.Sprintf("%s_%s_min",
		ctx.Values().GetString(auth.CtxClientName),
		clientPlatform,
	))

	updateAvailable := false
	updateForce := false
	if currAppVersion != "" {
		updateAvailable = semver.Compare(currAppVersion, clientAppVer) > 0
	}
	if minAppVersion != "" {
		updateForce = semver.Compare(minAppVersion, clientAppVer) >= 0
	}

	if ce := log.Check(log.DebugLevel, "Get Config"); ce != nil {
		ce.Write(
			zap.String("ClientVer", clientAppVer),
			zap.String("ClientPlatform", clientPlatform),
			zap.String("CurVer", currAppVersion),
			zap.String("MinVer", minAppVersion),
		)
	}

	res := &Config{
		UpdateAvailable: updateAvailable,
		UpdateForce:     updateForce,
		StoreLink:       getConfig(StoreLink),
		ShowBlipLink:    getConfig(ShowBlipLink) != "",
		ShowShareLink:   getConfig(ShowShareLink) != "",
	}
	s, ok := ctx.Values().Get(session.CtxSession).(*session.Session)
	if ok {
		res.Authorized = true
		u, _ := user.Get(s.UserID)
		if u != nil {
			res.VasEnabled = u.VasPaid
		}
	}

	msg.WriteResponse(ctx, CConfig, res)
}

// FeedbackHandler is API Handler
// Http Method: POST /help/feedback
// Returns: Bool (BOOL)
// Possible Errors:
//	1. 500: WRITE_TO_DB
//	2. 400: CANNOT_UNMARSHAL_JSON
func FeedbackHandler(ctx iris.Context) {
	req := &Feedback{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	s, _ := ctx.Values().Get(session.CtxSession).(session.Session)
	if ce := log.Check(log.DebugLevel, "Feedback Received"); ce != nil {
		ce.Write(
			zap.String("UserID", s.UserID))
	}


	_, err = feedbackCol.InsertOne(nil,
		bson.M{
			"txt":     req.Text,
			"rate":    req.Rate,
			"user_id": s.UserID,
		},
	)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}
