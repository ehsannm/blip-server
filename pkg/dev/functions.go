package dev

import (
	"database/sql"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
)

/*
   Creation Time: 2020 - Feb - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func MigrateLegacyDB() {
	db, err := sqlx.Connect("mysql", "ehsan:ZOAPcQf7rs8hRV02@(139.59.191.4:3306)/blip")
	if err != nil {
		log.Warn("Error On Connect MySql", zap.Error(err))
		return
	}
	rows, err := db.Query("SELECT artist, title, uri_local, cover FROM archives WHERE uri_local != ''")
	if err != nil {
		log.Warn("Error On Query", zap.Error(err))
		return
	}
	waitGroup := sync.WaitGroup{}
	rateLimit := make(chan struct{}, 10)
	cnt := 0
	for rows.Next() {
		if cnt++; cnt > 10 {
			break
		}
		var artist, title, uriLocal, cover sql.NullString
		err = rows.Scan(&artist, &title, &uriLocal, &cover)
		if err != nil {
			log.Warn("Error On Scan Legacy DB", zap.Error(err))
			continue
		}
		if !artist.Valid || !title.Valid || !uriLocal.Valid || !cover.Valid {
			continue
		}
		waitGroup.Add(1)
		rateLimit <- struct{}{}
		go func(artist, title, songUrl, coverUrl string) {
			defer waitGroup.Done()
			defer func() {
				<-rateLimit
			}()
			uniqueKey := music.GenerateUniqueKey(title, artist)
			songX, err := music.GetSongByUniqueKey(uniqueKey)
			if err != nil {
				songX = &music.Song{
					ID:             primitive.NilObjectID,
					UniqueKey:      uniqueKey,
					Title:          title,
					Genre:          "",
					Lyrics:         "",
					Artists:        artist,
					StoreID:        0,
					OriginCoverUrl: coverUrl,
					OriginSongUrl:  songUrl,
					Source:         "Archive",
				}
				songID, err := music.SaveSong(songX)
				if err != nil {
					log.Warn("Error On Save Search Result",
						zap.Error(err),
						zap.String("Title", title),
					)
					return
				}
				downloadFromSource("songs", songID, songUrl)
				downloadFromSource("covers", songID, coverUrl)
				return
			}
			songX.Artists = artist
			songX.Title = title
			songX.OriginSongUrl = songUrl
			songX.OriginCoverUrl = coverUrl
			songX.Source = "Archive"
			songID, err := music.SaveSong(songX)
			if err != nil {
				log.Warn("Error On Save Search Result",
					zap.Error(err),
				)
				return
			}
			downloadFromSource("songs", songID, songUrl)
			downloadFromSource("covers", songID, coverUrl)

		}(artist.String, title.String, uriLocal.String, cover.String)
	}
	waitGroup.Wait()

}
func downloadFromSource(bucketName string, songID primitive.ObjectID, url string) int64 {
	// download from source url
	storeID, dbWriter, err := store.GetUploadStream(bucketName, songID)
	if err != nil {
		log.Warn("Error On GetUploadStream", zap.Error(err))
		return 0
	}
	defer dbWriter.Close()

	res, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Warn("Error On Read From Source", zap.Error(err), zap.String("Url", url))
		return 0
	}
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		_, err = io.Copy(dbWriter, res.Body)
		if err != nil {
			log.Warn("Error On Copy", zap.Error(err))
			return 0
		}

	default:
		log.Warn("Invalid HTTP Status", zap.String("Status", res.Status))
		return 0
	}

	log.Info("Download Successfully", zap.String("Url", url))

	return storeID
}
