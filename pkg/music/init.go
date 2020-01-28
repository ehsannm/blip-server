package music

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"path/filepath"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate easyjson
var (
	songIndex  bleve.Index
	songCol    *mongo.Collection
	redisCache *ronak.RedisCache
)

func InitMongo(c *mongo.Client) {
	songCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColSong)
}

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func Init() {
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
