package music

import (
	"crypto/tls"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
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
	searchContexts    map[string]*searchCtx
	searchContextsMtx sync.RWMutex
	httpClient        http.Client
)

func InitMongo(c *mongo.Client) {
	songCol = c.Database(config.Db).Collection(config.ColSong)
}

func Init() {
	if proxyURL, err := url.Parse(config.GetString(config.HttpProxy)); err == nil {
		http.DefaultTransport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	} else {
		log.Warn("Error On Set HTTP Proxy", zap.Error(err))
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	searchContexts = make(map[string]*searchCtx)
	initSongIndex()
	// updateSongIndex()
	go watchForSongs()
}
func initSongIndex() {
	path := filepath.Join(config.GetString(config.SongsIndexDir), "songs")
	if index, err := bleve.Open(path); err != nil {
		switch err {
		case bleve.ErrorIndexPathDoesNotExist:
			indexMapping := bleve.NewIndexMapping()
			index, err = bleve.New(path, indexMapping)
			if err != nil {
				log.Fatal("Error On Init Music", zap.Error(errors.Wrap(err, "Search(Message)")))
				return
			}
			songIndex = index
		default:
			log.Fatal("Error On Init Music", zap.Error(errors.Wrap(err, "Search(Message)")))
			return
		}
	} else {
		songIndex = index
	}
}
func updateSongIndex() {
	startTime := time.Now()
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
			updateLocalIndex(songX)
			if cnt%1000 == 0 {
				log.Info("Still indexing ...", zap.Int("Indexed", cnt))
			}
		}
	}
	log.Info("Indexing songs done!.",
		zap.Int("Indexed", cnt),
		zap.Duration("Elapsed Time", time.Now().Sub(startTime)),
	)
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

				updateLocalIndex(songX)
				log.Debug("Song Indexed", zap.String("ID", songX.ID.Hex()))
			}
		}
		_ = stream.Close(nil)
	}
}
