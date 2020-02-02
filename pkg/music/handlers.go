package music

import (
	"encoding/base64"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"net/http"
	"time"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func AddStore(ctx iris.Context) {
	// TODO:: implement it
}

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
	sound := ctx.PostValue("sound")
	soundBytes, err := base64.StdEncoding.DecodeString(sound)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
		return
	}

	foundMusic, err := acr.IdentifyByByteString(soundBytes)
	if err != nil {
		log.Warn("Error On SearchBySound",
			zap.Error(err),
			zap.String("SessionID", ctx.GetHeader(session.HdrSessionID)),
		)
		msg.WriteError(ctx, http.StatusNotAcceptable, msg.Item(err.Error()))
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

func SearchByText(ctx iris.Context) {
	req := &SearchReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	songChan := StartSearch(ctx.GetHeader(session.HdrSessionID), req.Keyword)

	songIDs, err := SearchLocalIndex(req.Keyword)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrLocalIndexFailure)
		return
	}

	songs, err := GetManySongs(songIDs)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromDb)
		return
	}

	if len(songs) == 0 {
		songX, ok := <-songChan
		if ok {
			songs = append(songs, songX)
		}
	MainLoop:
		for {
			select {
			case songX, ok := <-songChan:
				if !ok {
					break MainLoop
				}
				songs = append(songs, songX)
			default:

			}
		}
	}

	msg.WriteResponse(ctx, CSearchResult, &SearchResult{
		Songs: songs,
	})
}

func SearchByCursor(ctx iris.Context) {
	songChan := ResumeSearch(ctx.GetHeader(session.HdrSessionID))
	t := time.NewTimer(time.Minute)
	t.Stop()
	songs := make([]*Song, 0, 100)
MainLoop:
	for {
		select {
		case songX, ok := <-songChan:
			if !ok {
				break MainLoop
			}
			songs = append(songs, songX)
			t.Reset(time.Second)
		case <-t.C:
			break MainLoop
		}
	}
	msg.WriteResponse(ctx, CSearchResult, &SearchResult{
		Songs: songs,
	})
}

func Download(ctx iris.Context) {
	// TODO:: implement it
}

func Upload(ctx iris.Context) {
	// TODO:: implement it
}
