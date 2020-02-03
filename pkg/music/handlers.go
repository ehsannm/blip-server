package music

import (
	"crypto/tls"
	"encoding/base64"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/store"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
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
	req.Keyword = strings.Trim(req.Keyword, "\" ")
	songChan := StartSearch(ctx.GetHeader(session.HdrSessionID), req.Keyword)
	songIDs, err := SearchLocalIndex(req.Keyword)
	if err != nil {
		log.Warn("Error On LocalIndex", zap.Error(err))
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
	}
	msg.WriteResponse(ctx, CSearchResult, &SearchResult{
		Songs: songs,
	})
}

func SearchByCursor(ctx iris.Context) {
	songChan := ResumeSearch(ctx.GetHeader(session.HdrSessionID))
	if songChan == nil {
		msg.WriteError(ctx, http.StatusAlreadyReported, msg.ErrAlreadyServed)
		return
	}
	t := time.NewTimer(time.Second * 5) // Wait
	songs := make([]*Song, 0, 100)
MainLoop:
	for {
		select {
		case songX, ok := <-songChan:
			if !ok {
				break MainLoop
			}
			songs = append(songs, songX)
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second) // WaitAfter
		case <-t.C:
			break MainLoop
		}
	}
	msg.WriteResponse(ctx, CSearchResult, &SearchResult{
		Songs: songs,
	})
}

func Download(ctx iris.Context) {
	downloadID := ctx.Params().GetString("downloadID")

	songID, err := primitive.ObjectIDFromHex(downloadID)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrInvalidDownloadID)
		return
	}

	songX, err := GetSongByID(songID)
	if err != nil {
		msg.WriteError(ctx, http.StatusNotFound, msg.ErrInvalidDownloadID)
		return
	}

	if songX.StoreID != 0 {
		err = store.DownloadSong(songX.StoreID, songX.ID, ctx.ResponseWriter())
		if err != nil {
			log.Warn("Error On Download Song", zap.Error(err))
		}
	} else {
		downloadFromSource(ctx, songX)
	}
	return
}
func downloadFromSource(ctx iris.Context, songX *Song) {
	// download from source url
	storeID, dbWriter, err := store.GetUploadStreamForSong(songX.ID)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}
	writer := io.MultiWriter(dbWriter, ctx.ResponseWriter())
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := http.Client{}
	res, err := httpClient.Get(songX.OriginSongUrl)
	if err != nil {
		log.Warn("Error On Read From Source", zap.Error(err), zap.String("Url", songX.OriginSongUrl))
		msg.WriteError(ctx, http.StatusFailedDependency, msg.ErrReadFromSource)
		return
	}
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		_, err = io.Copy(writer, res.Body)
		if err != nil {
			msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
			return
		}
		_ = dbWriter.Close()
	default:
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}
	songX.StoreID = storeID
	_, _ = SaveSong(songX)
}

func Upload(ctx iris.Context) {
	// TODO:: implement it
}
