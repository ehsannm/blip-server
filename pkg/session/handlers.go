package session

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris/v12"
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

// MustHaveSession is a middleware which make sure context has valid sessionID otherwise it stops
// the context and response the client with appropriate message
// On Error: http status is StatusForbidden (403) and the error message is SESSION_INVALID
func MustHaveSession(ctx iris.Context) {
	sessionID := ctx.GetHeader(HdrSessionID)
	sessionCacheMtx.RLock()
	session, ok := sessionCache[sessionID]
	sessionCacheMtx.RUnlock()
	if !ok {
		if s, err := Get(sessionID); err != nil {
			if ce := log.Check(log.DebugLevel, "Error On GetSession"); ce != nil {
				ce.Write(
					zap.Error(err),
					zap.String("SessionID", sessionID),
				)
			}
			msg.WriteError(ctx, http.StatusForbidden, msg.ErrSessionInvalid)
			return
		} else {
			session = s
		}
	}

	sessionCacheMtx.Lock()
	sessionCache[sessionID] = session
	sessionCacheMtx.Unlock()
	ctx.Values().Save(CtxSession, session, true)
	ctx.Next()
}
