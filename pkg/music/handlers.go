package music

import (
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func SearchByProxy(ctx iris.Context) {
	if ce := log.Check(log.DebugLevel, "SearchByProxy"); ce != nil {
		s, ok := ctx.Values().Get(session.CtxSession).(session.Session)
		if ok {
			ce.Write(zap.String("UserID", s.UserID))
		} else {
			ce.Write(zap.String("UserID", "Not Set"))
		}

	}
	reverseProxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}



func SearchBySound(ctx iris.Context) {
	// sound := ctx.PostValue("sound")
	// soundBytes, err := base64.StdEncoding.DecodeString(sound)
	// if err != nil {
	// 	msg.Error(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
	// 	return
	// }

	// music, err := acr.IdentifyByByteString(soundBytes)
	// if err != nil {
	// 	msg.Error(ctx, http.StatusNotAcceptable, msg.Item(err.Error()))
	// 	return
	// }

	// music.Status.Code

}
