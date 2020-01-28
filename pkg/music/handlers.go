package music

import (
	"encoding/base64"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"net/http"
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

func SearchByText(ctx iris.Context) {
	// TODO:: implement it
}

func SearchBySound(ctx iris.Context) {
	sound := ctx.PostValue("sound")
	soundBytes, err := base64.StdEncoding.DecodeString(sound)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
		return
	}

	foundMusic, err := acr.IdentifyByByteString(soundBytes)
	if err != nil {
		log.Warn("Error On SearchBySound",
			zap.Error(err),
			zap.String("SessionID", ctx.GetHeader(session.HdrSessionID)),
		)
		msg.Error(ctx, http.StatusNotAcceptable, msg.Item(err.Error()))
		return
	}

	// TODO: We must do the following steps
	// #1. Search Local Database for musics and return a result with with a cursorID to the client
	// #2. For each crawler send search request
	// #3. Create an in memory object holding information about pending request
	for _, m := range foundMusic.Metadata.Music {
		_ = m
	}
}

func SearchByCursor(ctx iris.Context) {
	// TODO:: implement it
}

func Download(ctx iris.Context) {
	// TODO:: implement it
}
