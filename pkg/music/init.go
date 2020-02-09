package music

import (
	"crypto/tls"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
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
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	if config.GetString(config.HttpProxy) != "" {
		if proxyURL, err := url.Parse(config.GetString(config.HttpProxy)); err == nil {
			httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
		} else {
			log.Warn("Error On Set HTTP Proxy", zap.Error(err))
		}
	}

	searchContexts = make(map[string]*searchCtx)
	log.Info("Initialize Songs Index ...")
	initSongIndex()
	log.Info("Songs Index Initialized. ")
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
