package crawler

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
//go:generate easyjson crawler.go request.go
func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func InitMongo(c *mongo.Client) {
	crawlerCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColCrawler)
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
	go func() {
		stream, err := crawlerCol.Watch(nil, mongo.Pipeline{})
		if err != nil {
			log.Warn("Crawler Watcher has been Closed.", zap.Error(err))
			return
		}
		defer stream.Close(nil)
		for stream.Next(nil) {
			c := &Crawler{}
			err := stream.Decode(c)
			if err == nil {
				registeredCrawlersMtx.Lock()
				registeredCrawlers[c.Source] = append(registeredCrawlers[c.Source], c)
				registeredCrawlersMtx.Unlock()
			}
		}
		log.Warn("Crawler Watcher has been Closed.")
	}()
}
