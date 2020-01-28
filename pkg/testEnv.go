package testEnv

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Init() {
	config.Set(config.MongoUrl, "mongodb://localhost:27017")
	config.Set(config.MongoDB, "blip")
	config.Set(config.RedisUrl, "localhost:6379")
	config.Set(config.LogLevel, log.DebugLevel)
	config.Set(config.SongsIndexDir, "./_hdd")

	os.MkdirAll(config.GetString(config.SongsIndexDir), os.ModePerm)

	// Initialize MongoDB
	mongoClient, err := mongo.Connect(
		nil,
		options.Client().ApplyURI(viper.GetString(config.MongoUrl)),
	)
	if err != nil {
		log.Fatal("Error On MongoConnect", zap.Error(err))
	}
	crawler.InitMongo(mongoClient)
	music.InitMongo(mongoClient)

	// Initialize RedisCache
	redisConfig := ronak.DefaultRedisConfig
	redisConfig.Host = viper.GetString(config.RedisUrl)
	redisConfig.Password = viper.GetString(config.RedisPass)
	redisCache := ronak.NewRedisCache(redisConfig)
	crawler.InitRedisCache(redisCache)
	music.InitRedisCache(redisCache)

	// Initialize Modules
	log.InitLogger(log.DebugLevel, "")
	music.Init()

}
