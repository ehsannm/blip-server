package admin

import (
	"database/sql"
	"fmt"
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
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2020 - Feb - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	migrateRunning           bool
	migrateScanned           int32
	migrateDownloaded        int32
	migrateDownloadFailed    int32
	migrateAlreadyDownloaded int32
)

func MigrateLegacyDB() {
	migrateScanned = 0
	migrateDownloaded = 0
	migrateDownloadFailed = 0
	migrateAlreadyDownloaded = 0
	db, err := sqlx.Connect("mysql", "ehsan:ZOAPcQf7rs8hRV02@(139.59.191.4:3306)/blip")
	if err != nil {
		log.Warn("Error On Connect MySql", zap.Error(err))
		return
	}

	waitGroup := sync.WaitGroup{}
	rateLimit := make(chan struct{}, 50)
	rows, err := db.Query("SELECT uid, artist, title, uri_local, cover FROM archives WHERE uri_local != '' ORDER by uid ASC")
	if err != nil {
		log.Warn("Error On Query", zap.Error(err))
		return
	}
	var uid, artist, title, uriLocal, cover sql.NullString
	for {
		for rows.Next() {
			err = rows.Scan(&uid, &artist, &title, &uriLocal, &cover)
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
				atomic.AddInt32(&migrateScanned, 1)
				defer waitGroup.Done()
				defer func() {
					<-rateLimit
				}()
				uniqueKey := music.GenerateUniqueKey(title, artist)
				songX, err := music.GetSongByUniqueKey(uniqueKey)
				if err != nil {
					songX = &music.Song{
						ID:             primitive.NewObjectID(),
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

					storeID := downloadFromSource(store.BucketSongs, songX.ID, songUrl)
					if storeID != 0 {
						songX.StoreID = storeID
						_, err := music.SaveSong(songX)
						if err != nil {
							log.Warn("Error On Save Search Result",
								zap.Error(err),
								zap.String("Title", title),
							)
							return
						}
						downloadFromSource(store.BucketCovers, songX.ID, coverUrl)
						atomic.AddInt32(&migrateDownloaded, 1)
					} else {
						atomic.AddInt32(&migrateDownloadFailed, 1)
					}

					return
				} else if songX.StoreID == 0 {
					storeID := downloadFromSource(store.BucketSongs, songX.ID, songUrl)
					if storeID != 0 {
						songX.StoreID = storeID
						_, _ = music.SaveSong(songX)
						downloadFromSource(store.BucketCovers, songX.ID, coverUrl)
						atomic.AddInt32(&migrateDownloaded, 1)
					} else {
						_ = music.DeleteSong(songX.ID)
						atomic.AddInt32(&migrateDownloadFailed, 1)
					}
					return
				}
				atomic.AddInt32(&migrateAlreadyDownloaded, 1)
			}(artist.String, title.String, uriLocal.String, cover.String)
		}
		waitGroup.Wait()
		if rows.Err() == nil {
			_ = rows.Close()
			break
		}
		log.Warn("Error On Rows", zap.Error(rows.Err()))
		_ = rows.Close()
		rows, err = db.Query(fmt.Sprintf(
			"SELECT uid, artist, title, uri_local, cover FROM archives WHERE uri_local != '' AND uid >= '%s' ORDER by uid ASC", uid.String,
		))
		if err != nil {
			log.Warn("Error On Query", zap.Error(err))
			return
		}
	}

	migrateRunning = false
	log.Info("Migration Finished",
		zap.Int32("Scanned", migrateScanned),
		zap.Int32("Downloaded", migrateDownloaded),
	)

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
		_ = dbWriter.SetWriteDeadline(time.Now().Add(time.Minute))
		_, err = io.Copy(dbWriter, res.Body)
		if err != nil {
			log.Warn("Error On Copy", zap.Error(err))
			return 0
		}

	default:
		log.Warn("Invalid HTTP Status",
			zap.String("Status", res.Status),
			zap.String("Url", url),
		)
		return 0
	}

	log.Info("Download Successfully", zap.String("Url", url))

	return storeID
}
