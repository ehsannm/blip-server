package testEnv

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"

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
	log.InitLogger(log.InfoLevel, "")
	config.Set(config.MongoUrl, "mongodb://localhost:27001")
	config.Set(config.MongoDB, "blip")
	config.Set(config.RedisUrl, "localhost:6379")
	config.Set(config.LogLevel, log.DebugLevel)
	config.Set(config.SongsIndexDir, "./_hdd")

	_ = os.MkdirAll(config.GetString(config.SongsIndexDir), os.ModePerm)

	// Initialize MongoDB
	mongoClient, err := mongo.Connect(
		nil,
		options.Client().ApplyURI(config.GetString(config.MongoUrl)).SetDirect(true),
	)
	if err != nil {
		log.Fatal("Error On MongoConnect", zap.Error(err))
	}
	err = mongoClient.Ping(nil, nil)
	if err != nil {
		log.Fatal("Error On MongoConnect (Ping)", zap.Error(err))
	}
	auth.InitMongo(mongoClient)
	crawler.InitMongo(mongoClient)
	music.InitMongo(mongoClient)
	session.InitMongo(mongoClient)
	token.InitMongo(mongoClient)
	user.InitMongo(mongoClient)

	// Initialize RedisCache
	redisConfig := redis.DefaultConfig
	redisConfig.Host = config.GetString(config.RedisUrl)
	redisConfig.Password = config.GetString(config.RedisPass)
	redisCache := redis.New(redisConfig)
	crawler.InitRedisCache(redisCache)
	music.InitRedisCache(redisCache)

	// Initialize Modules
	auth.Init()
	crawler.Init()
	music.Init()
	session.Init()
	token.Init()
	user.Init()

}
