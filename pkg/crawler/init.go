package crawler

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
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
func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func InitMongo(c *mongo.Client) {
	crawlerCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColCrawler)
}

func Init() {
	_, _ = GetAll()
}
