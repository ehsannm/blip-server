package session

import (
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
)

/*
   Creation Time: 2019 - Sep - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func MustHaveSession(ctx iris.Context) {
	sessionID := ctx.GetHeader(HdrSessionID)
	mtxLock.RLock()
	session, ok := sessionCache[sessionID]
	mtxLock.RUnlock()
	if !ok {
		res := sessionCol.FindOne(nil, bson.M{"_id": sessionID}, options.FindOne())
		if err := res.Decode(&session); err != nil {
			if ce := log.Check(log.DebugLevel, "Error On GetSession"); ce != nil {
				ce.Write(
					zap.Error(err),
					zap.String("SessionID", sessionID),
				)
			}
			msg.Error(ctx, http.StatusForbidden, msg.ErrSessionInvalid)
			return
		}
	}

	mtxLock.Lock()
	sessionCache[sessionID] = session
	mtxLock.Unlock()
	log.Debug("Session Detected", zap.String("UserID", session.UserID))
	ctx.Values().Save(CtxSession, session, true)
	ctx.Next()
}
