package crawler

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
//go:generate easyjson crawler.go messages.go
var (
	crawlerCol             *mongo.Collection
	registeredCrawlersMtx  sync.RWMutex
	registeredCrawlers     map[string][]*Crawler
	registeredCrawlersPool sync.Pool
)

func InitMongo(c *mongo.Client) {
	crawlerCol = c.Database(config.Db).Collection(config.ColCrawler)
}

func Init() {
	registeredCrawlersMtx.Lock()
	defer registeredCrawlersMtx.Unlock()
	registeredCrawlers = make(map[string][]*Crawler)
	cur, err := crawlerCol.Find(nil, bson.D{})

	if err != nil {
		log.Warn("Error On Initializing Crawlers", zap.Error(err))
	}
	for cur.Next(nil) {
		crawler := &Crawler{}
		err = cur.Decode(crawler)
		if err != nil {
			continue
		}
		registeredCrawlers[crawler.Source] = append(registeredCrawlers[crawler.Source], crawler)
	}
	go watchForCrawlers()
}
func watchForCrawlers() {
	var resumeToken bson.Raw
	for {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		if resumeToken != nil {
			opts.SetStartAfter(resumeToken)
		}
		stream, err := crawlerCol.Watch(nil, mongo.Pipeline{}, opts)
		if err != nil {
			log.Warn("Error On Watch Stream for Crawlers", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

	EventsLoop:
		for stream.Next(nil) {
			crawlerX := &Crawler{}
			resumeToken = stream.ResumeToken()
			operationType := strings.Trim(stream.Current.Lookup("operationType").String(), "\"")
			switch operationType {
			case "insert":
				err := stream.Current.Lookup("fullDocument").UnmarshalWithRegistry(bson.DefaultRegistry, crawlerX)
				if err != nil {
					log.Warn("Error On Decoding Crawler", zap.Error(err))
					continue
				}
				registeredCrawlersMtx.Lock()
				registeredCrawlers[crawlerX.Source] = append(registeredCrawlers[crawlerX.Source], crawlerX)
				registeredCrawlersMtx.Unlock()
			case "update":
				err := stream.Current.Lookup("fullDocument").UnmarshalWithRegistry(bson.DefaultRegistry, crawlerX)
				if err != nil {
					log.Warn("Error On Decoding Crawler", zap.Error(err))
					continue
				}
				registeredCrawlersMtx.Lock()
				for idx, c := range registeredCrawlers[crawlerX.Source] {
					if c.ID == crawlerX.ID {
						registeredCrawlers[crawlerX.Source][idx] = crawlerX
					}
				}
				registeredCrawlersMtx.Unlock()
			case "delete":
				crawlerID := stream.Current.Lookup("documentKey").ObjectID()
				registeredCrawlersMtx.Lock()
				for idx, c := range registeredCrawlers[crawlerX.Source] {
					if c.ID == crawlerID {
						registeredCrawlers[crawlerX.Source][idx] = registeredCrawlers[crawlerX.Source][len(registeredCrawlers[crawlerX.Source])-1]
						registeredCrawlers[crawlerX.Source] = registeredCrawlers[crawlerX.Source][:len(registeredCrawlers[crawlerX.Source])-1]
					}
				}
				registeredCrawlersMtx.Unlock()
				continue
			case "invalidate", "drop":
				registeredCrawlersMtx.Lock()
				registeredCrawlers = make(map[string][]*Crawler)
				registeredCrawlersMtx.Unlock()
				break EventsLoop
			default:
				log.Warn("Unknown Operation Type", zap.String("OT", operationType))
			}

			log.Debug("Crawler Found",
				zap.String("Url", crawlerX.Url),
				zap.String("Source", crawlerX.Source),
			)
		}
		_ = stream.Close(nil)
	}

}
