package music

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pkg/errors"
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
	songIndex         bleve.Index
	songCol           *mongo.Collection
	storeCol          *mongo.Collection
	stores            map[int64]*Store
	storesMtx         sync.RWMutex
	searchContexts    map[string]*searchCtx
	searchContextsMtx sync.RWMutex
)

func InitMongo(c *mongo.Client) {
	songCol = c.Database(config.GetString(config.MongoDB)).Collection(config.ColSong)
	storeCol = c.Database(config.GetString(config.MongoDB)).Collection(config.ColStore)
}

func Init() {
	searchContexts = make(map[string]*searchCtx)
	stores = make(map[int64]*Store)
	initSongIndex()
	updateSongIndex()
	go watchForSongs()
	go watchForStores()

}
func initSongIndex() {
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
func updateSongIndex() {
	log.Info("Indexing songs, this may take time ...")
	cur, err := songCol.Find(nil, bson.D{})
	if err != nil {
		log.Fatal("Error On Initializing Music", zap.Error(err))
	}
	cnt := 0
	for cur.Next(nil) {
		songX := &Song{}
		err = cur.Decode(songX)
		if err == nil {
			cnt++
			err = UpdateLocalIndex(songX)
			log.WarnOnError("Error On Indexing Song", err, zap.String("SongID", songX.ID.Hex()))
			if cnt%1000 == 0 {
				log.Info("Still indexing ...", zap.Int("Indexed", cnt))
			}
		}
	}
	log.Info("Indexing songs done!.", zap.Int("Indexed", cnt))
}
func watchForSongs() {
	var resumeToken bson.Raw
	for {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		if resumeToken != nil {
			opts.SetStartAfter(resumeToken)
		}
		stream, err := songCol.Watch(nil, mongo.Pipeline{}, opts)
		if err != nil {
			log.Warn("Error On Watch Stream for Songs", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		for stream.Next(nil) {
			songX := &Song{}
			resumeToken = stream.ResumeToken()
			operationType := strings.Trim(stream.Current.Lookup("operationType").String(), "\"")
			switch operationType {
			case "insert", "update":
				err = stream.Current.Lookup("fullDocument").UnmarshalWithRegistry(bson.DefaultRegistry, songX)
				if err != nil {
					log.Warn("Error On Decoding Song", zap.Error(err))
					continue
				}

				_ = UpdateLocalIndex(songX)
				log.Debug("Song Indexed", zap.String("ID", songX.ID.Hex()))
			}
		}
		_ = stream.Close(nil)
	}
}
func watchForStores() {
	var resumeToken bson.Raw
	for {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		if resumeToken != nil {
			opts.SetStartAfter(resumeToken)
		}
		stream, err := storeCol.Watch(nil, mongo.Pipeline{}, opts)
		if err != nil {
			log.Warn("Error On Watch Stream for Stores", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		for stream.Next(nil) {
			storeX := &Store{}
			resumeToken = stream.ResumeToken()
			operationType := strings.Trim(stream.Current.Lookup("operationType").String(), "\"")
			switch operationType {
			case "insert", "update":
				err = stream.Current.Lookup("fullDocument").UnmarshalWithRegistry(bson.DefaultRegistry, storeX)
				if err != nil {
					log.Warn("Error On Decoding Store", zap.Error(err))
					continue
				}
				storesMtx.Lock()
				stores[storeX.ID] = storeX
				storesMtx.Unlock()
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
