package admin

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/store"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/url"
	"sync/atomic"
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
	healthCheckRunning bool
	scanned            int32
	coverFixed         int32
	songFixed          int32
)

func HealthCheckDB() {
	scanned = 0
	coverFixed = 0
	songFixed = 0

	err := music.ForEachSong(func(songX *music.Song) bool {
		atomic.AddInt32(&scanned, 1)
		downloadSong := false
		downloadCover := false
		if songX.SongStoreID != 0 {
			err := store.FileExists(store.BucketSongs, songX.SongStoreID, songX.ID)
			if err != nil {
				downloadSong = true
			}
		} else {
			downloadSong = true
		}
		if downloadSong {
			if _, err := url.Parse(songX.OriginSongUrl); err != nil {
				_ = music.DeleteSong(songX.ID)
				return false
			}
			songStoreID := music.DownloadFromSource(store.BucketSongs, songX.ID, songX.OriginSongUrl)
			if songStoreID != 0 {
				atomic.AddInt32(&songFixed, 1)
				songX.SongStoreID = songStoreID
			} else {
				_ = music.DeleteSong(songX.ID)
				return true
			}
		}

		if songX.CoverStoreID != 0 {
			err := store.FileExists(store.BucketCovers, songX.CoverStoreID, songX.ID)
			if err != nil {
				downloadCover = true
			}
		} else {
			downloadCover = true
		}
		if downloadCover {
			if _, err := url.Parse(songX.OriginCoverUrl); err != nil {
				_ = music.DeleteSong(songX.ID)
				return false
			}
			coverStoreID := music.DownloadFromSource(store.BucketCovers, songX.ID, songX.OriginCoverUrl)
			if coverStoreID != 0 {
				atomic.AddInt32(&coverFixed, 1)
				songX.CoverStoreID = coverStoreID
			}
		}

		if downloadCover || downloadSong {
			_, err := music.SaveSong(songX)
			if err != nil {
				log.Warn("Error On Save Song",
					zap.Error(err),
					zap.String("Title", songX.ID.Hex()),
				)
				return false
			}
		}

		return true
	})
	if err != nil {
		log.Warn("Error On ForEachSong", zap.Error(err))
	}

	healthCheckRunning = false
	log.Info("HealthCheckDB Finished",
		zap.Int32("Scanned", scanned),
		zap.Int32("CoverFixed", coverFixed),
		zap.Int32("SongFixed", songFixed),
	)

}

func HealthCheckStore() {
	scanned = 0
	coverFixed = 0
	songFixed = 0
	err := store.ForEachSong(store.BucketSongs, 101, func(songID primitive.ObjectID) bool {
		atomic.AddInt32(&scanned, 1)
		songX, err := music.GetSongByID(songID)
		if err == nil && songX != nil {
			return true
		}
		err = store.FileDelete(store.BucketSongs, 101, songID)
		if err != nil {
			log.Warn("Error On FileDelete", zap.Error(err))
			return false
		}
		atomic.AddInt32(&songFixed, 1)
		return true
	})
	if err != nil {
		log.Warn("Error On ForEachSong", zap.Error(err))
	}

	err = store.ForEachSong(store.BucketCovers, 101, func(songID primitive.ObjectID) bool {
		atomic.AddInt32(&scanned, 1)
		songX, err := music.GetSongByID(songID)
		if err == nil && songX != nil {
			return true
		}
		err = store.FileDelete(store.BucketCovers, 101, songID)
		if err != nil {
			log.Warn("Error On FileDelete", zap.Error(err))
			return false
		}
		atomic.AddInt32(&coverFixed, 1)
		return true
	})
	if err != nil {
		log.Warn("Error On ForEachSong", zap.Error(err))
	}

	healthCheckRunning = false
	log.Info("HealthCheckStore Finished",
		zap.Int32("Scanned", scanned),
		zap.Int32("CoverFixed", coverFixed),
		zap.Int32("SongFixed", songFixed),
	)

}
