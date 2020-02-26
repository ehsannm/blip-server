package music

import (
	"encoding/base64"
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/store"
	"git.ronaksoftware.com/blip/server/pkg/store/gridfs"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"io/ioutil"
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

// SearchBySoundHandler is API Handler
// This is a reverse proxy end point
// @Deprecated
func SearchByProxyHandler(ctx iris.Context) {
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

// SearchBySoundHandler is API Handler
// Http Method: POST  /music/search/sound
// Inputs: POST VALUES:
//	1. "sound" ->  Based64 Standard Encoded
// Returns: SoundSearchResult (SOUND_SEARCH_RESULT)
// Possible Errors:
//	1. 400: BAD_SOUND_FILE
//	2. 406: SEARCH_ENGINE
//	3. 500: LOCAL_INDEX_FAILED
//  4. 403: SONG_NOT_FOUND
func SearchBySoundHandler(ctx iris.Context) {
	soundFile, _, err := ctx.FormFile("sound")
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
		return
	}
	soundBytes, err := ioutil.ReadAll(soundFile)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
		return
	}

	log.Debug("Received Sound",
		zap.Int("Len", len(soundBytes)),
	)

	foundMusic, err := acr.IdentifyByByteString(soundBytes)
	if err != nil {
		log.Warn("Error On SearchBySound",
			zap.Error(err),
			zap.String("SessionID", ctx.GetHeader(session.HdrSessionID)),
		)
		msg.WriteError(ctx, http.StatusNotAcceptable, msg.ErrSearchEngine)
		return
	}

	if len(foundMusic.Metadata.Music) == 0 {
		msg.WriteError(ctx, http.StatusNotFound, msg.ErrSongNotFound)
		return
	}
	var keyword string
	if len(foundMusic.Metadata.Music[0].Artists) > 0 {
		keyword = foundMusic.Metadata.Music[0].Artists[0].Name
	}
	keyword = fmt.Sprintf("%s+%s", keyword, foundMusic.Metadata.Music[0].Title)
	songChan := StartSearch(ctx.GetHeader(session.HdrSessionID), keyword)
	indexedSongs, err := SearchLocalIndex(keyword)
	if err != nil {
		log.Warn("Error On LocalIndex", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrLocalIndexFailure)
		return
	}
	songs := make([]*Song, 0, len(indexedSongs))
	for idx := range indexedSongs {
		songs = append(songs, indexedSongs[idx].song)
	}

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
	res := &SoundSearchResult{
		Songs: songs,
	}
	res.Info.Title = foundMusic.Metadata.Music[0].Title
	res.Info.ReleaseDate = foundMusic.Metadata.Music[0].ReleaseDate
	for _, artist := range foundMusic.Metadata.Music[0].Artists {
		res.Info.Artists = append(res.Info.Artists, artist.Name)
	}
	msg.WriteResponse(ctx, CSoundSearchResult, res)
}

// SearchByFingerprintHandler is API Handler
// Http Method: POST  /music/search/fingerprint
// Inputs: POST VALUES:
//	1. "fingerprint" ->  byte slice
// Returns: SoundSearchResult (SOUND_SEARCH_RESULT)
// Possible Errors:
//	1. 500: LOCAL_INDEX_FAILED
//	2. 406: SEARCH_ENGINE
//  3. 403: SONG_NOT_FOUND
func SearchByFingerprintHandler(ctx iris.Context) {
	fingerprint := ctx.FormValue("fingerprint")
	if ce := log.Check(log.DebugLevel, "Received Fingerprint"); ce != nil {
		ce.Write(
			zap.Int("Len", len(fingerprint)),
		)
	}
	decodeFP, err := base64.StdEncoding.DecodeString(fingerprint)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCorruptData)
		return
	}

	foundMusic, err := acr.IdentifyByFingerprint(decodeFP)
	if err != nil {
		log.Warn("Error On SearchBySound",
			zap.Error(err),
			zap.String("SessionID", ctx.GetHeader(session.HdrSessionID)),
		)
		msg.WriteError(ctx, http.StatusNotAcceptable, msg.ErrSearchEngine)
		return
	}

	if ce := log.Check(log.DebugLevel, "ACR Result"); ce != nil {
		ce.Write(
			zap.String("Status.Msg", foundMusic.Status.Message),
			zap.Int("Status.Code", foundMusic.Status.Code),
			zap.Int("ResultType", foundMusic.ResultType),
			zap.Int("Total", len(foundMusic.Metadata.Music)),
		)
	}

	if len(foundMusic.Metadata.Music) == 0 {
		msg.WriteError(ctx, http.StatusNotFound, msg.ErrSongNotFound)
		return
	}
	var keyword string
	if len(foundMusic.Metadata.Music[0].Artists) > 0 {
		keyword = foundMusic.Metadata.Music[0].Artists[0].Name
	}
	keyword = fmt.Sprintf("%s+%s", keyword, foundMusic.Metadata.Music[0].Title)
	_ = StartSearch(ctx.GetHeader(session.HdrSessionID), keyword)
	indexedSongs, err := SearchLocalIndex(keyword)
	if err != nil {
		log.Warn("Error On LocalIndex", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrLocalIndexFailure)
		return
	}
	songs := make([]*Song, 0, len(indexedSongs))
	for idx := range indexedSongs {
		songs = append(songs, indexedSongs[idx].song)
	}

	res := &SoundSearchResult{
		Songs: songs,
	}
	res.Info.Title = foundMusic.Metadata.Music[0].Title
	res.Info.ReleaseDate = foundMusic.Metadata.Music[0].ReleaseDate
	for _, artist := range foundMusic.Metadata.Music[0].Artists {
		res.Info.Artists = append(res.Info.Artists, artist.Name)
	}
	msg.WriteResponse(ctx, CSoundSearchResult, res)
}

// SearchByTextHandler is API Handler
// Http Method: POST /music/search/text
// Inputs: JSON - SearchReq
// Returns: SearchResult (SEARCH_RESULT)
// Possible Errors:
//	1. 400: CANNOT_MARSHAL_JSON
//	2. 500: LOCAL_INDEX_FAILED
//	3. 500: READ_FROM_DB
func SearchByTextHandler(ctx iris.Context) {
	req := &SearchReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	req.Keyword = strings.Trim(req.Keyword, "\"")
	songChan := StartSearch(ctx.GetHeader(session.HdrSessionID), req.Keyword)
	indexedSongs, err := SearchLocalIndex(req.Keyword)
	if err != nil {
		log.Warn("Error On LocalIndex", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrLocalIndexFailure)
		return
	}
	songs := make([]*Song, 0, len(indexedSongs))
	for idx := range indexedSongs {
		songs = append(songs, indexedSongs[idx].song)
	}
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

// SearchByCursorHandler is API handler
// API: /music/search
// Http Method: GET
// Inputs: N/A
// Returns: SearchResult (SEARCH_RESULT)
// Possible Errors:
//	1. 208: ALREADY_SERVED
func SearchByCursorHandler(ctx iris.Context) {
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

// DownloadHandler is API handler
// Http Method: GET /music/download/{bucket}/{downloadID}
// Inputs:	URL
// Possible Bucket Values: songs, covers
func DownloadHandler(ctx iris.Context) {
	downloadID := ctx.Params().GetString("downloadID")
	bucketName := strings.ToLower(ctx.Params().GetString("bucket"))
	switch bucketName {
	case store.BucketCovers, store.BucketSongs:
	default:
		msg.WriteError(ctx, http.StatusNotFound, msg.ErrInvalidUrl)
		return
	}

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

	startTime := time.Now()
	switch bucketName {
	case store.BucketSongs:
		if songX.SongStoreID == 0 {
			songStoreID := DownloadFromSource(store.BucketSongs, songX.ID, songX.OriginSongUrl)
			if songStoreID > 0 {
				songX.SongStoreID = songStoreID
				_, err := SaveSong(songX)
				if err != nil {
					log.Warn("Error On Save Search Result",
						zap.Error(err),
						zap.String("Title", songX.ID.Hex()),
					)
					msg.WriteError(ctx, http.StatusNotFound, msg.ErrReadFromSource)
					return
				}
			} else {
				_ = DeleteSong(songX.ID)
				msg.WriteError(ctx, http.StatusNotFound, msg.ErrReadFromSource)
				return
			}
		}
		fallthrough
	case store.BucketCovers:
		if songX.CoverStoreID == 0 {
			coverStoreID := DownloadFromSource(store.BucketCovers, songX.ID, songX.OriginCoverUrl)
			if coverStoreID > 0 {
				songX.CoverStoreID = coverStoreID
				_, err := SaveSong(songX)
				if err != nil {
					log.Warn("Error On Save Search Result",
						zap.Error(err),
						zap.String("Title", songX.ID.Hex()),
					)
				}
			}
		}
	}

	var dbReader *gridfs.DownloadStream
	switch bucketName {
	case store.BucketCovers:
		dbReader, err = store.GetDownloadStream(bucketName, songX.ID, songX.CoverStoreID)
	case store.BucketSongs:
		dbReader, err = store.GetDownloadStream(bucketName, songX.ID, songX.SongStoreID)
	}
	if err != nil {
		log.Warn("Error On Download Song (GetDownloadStream)", zap.Error(err))
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrReadFromDb)
		return
	}
	defer dbReader.Close()
	_, err = store.Copy(ctx.ResponseWriter(), dbReader, ctx.ResponseWriter().Flush)
	if err != nil {
		log.Warn("Error On Download Song (Copy)", zap.Error(err))
		ctx.StatusCode(http.StatusServiceUnavailable)
		return
	}

	if ce := log.Check(log.DebugLevel, "Song Download"); ce != nil {
		ce.Write(
			zap.String("Bucket", bucketName),
			zap.String("SongID", songX.ID.Hex()),
			zap.Duration("Time", time.Now().Sub(startTime)),
		)

	}
	return
}
