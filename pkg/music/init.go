package music

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate rm -f *_easyjson.go
//go:generate easyjson song.go messages.go
var (
	songIndex  bleve.Index
	songCol    *mongo.Collection
	redisCache *redis.Cache
)

func InitMongo(c *mongo.Client) {
	songCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColSong)
}

func InitRedisCache(c *redis.Cache) {
	redisCache = c
}

func Init() {
	initializeIndex()
	updateIndex()
	go watchForSongs()

}
func initializeIndex() {
	path := filepath.Join(config.GetString(config.SongsIndexDir), "songs")
	if index, err := bleve.Open(path); err != nil {
		switch err {
		case bleve.ErrorIndexPathDoesNotExist:
			// create a mapping
			indexMapping := bleve.NewIndexMapping()
			// indexMapping, err := indexMapForSongs()
			// if err != nil {
			// 	log.Fatal("Error On Init Music", zap.Error(errors.Wrap(err, "Search(Message)")))
			// 	return
			// }
			index, err = bleve.New(path, indexMapping)
			if err != nil {
				log.Fatal("Error On Init Music", zap.Error(errors.Wrap(err, "Search(Message)")))
				return
			}
		default:
			log.Fatal("Error On Init Music", zap.Error(errors.Wrap(err, "Search(Message)")))
			return
		}
	} else {
		songIndex = index
	}
}
func updateIndex() {
	log.Info("Indexing songs, this may take time ...")
	cur, err := songCol.Find(nil, bson.D{})
	if err != nil {
		log.Fatal("Error On Initializing Music", zap.Error(err))
	}
	for cur.Next(nil) {
		songX := &Song{}
		err = cur.Decode(songX)
		if err == nil {
			err = songIndex.Index(songX.ID.Hex(), songX)
			log.WarnOnError("Error On Indexing Song", err, zap.String("SongID", songX.ID.Hex()))
		}
	}
	log.Info("Indexing songs done!.")
}
func watchForSongs() {
	for {
		stream, err := songCol.Watch(nil, mongo.Pipeline{})
		if err != nil {
			log.Warn("Error On Watch Stream for Songs", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		for stream.Next(nil) {
			songX := &Song{}
			err = stream.Decode(songX)
			if err == nil {
				_ = songIndex.Index(songX.ID.Hex(), songX)
			}
		}
		_ = stream.Close(nil)
	}

}

func indexMapForSongs() (mapping.IndexMapping, error) {
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name
	keywordFieldMapping.DocValues = true
	keywordFieldMapping.IncludeTermVectors = true

	// Song
	songMapping := bleve.NewDocumentStaticMapping()
	songMapping.AddFieldMappingsAt("lyrics", keywordFieldMapping)
	songMapping.AddFieldMappingsAt("title", keywordFieldMapping)
	songMapping.AddFieldMappingsAt("artists", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("song", songMapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = keyword.Name

	return indexMapping, nil
}
