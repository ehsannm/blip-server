package main

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/admin"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/store"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initModules() {
	config.Init()
	log.InitLogger(zapcore.Level(config.GetInt(config.LogLevel)), "")

	// Initialize MongoDB
	if mongoClient, err := mongo.Connect(nil,
		options.Client().ApplyURI(config.GetString(config.MongoUrl)),
	); err != nil {
		log.Fatal("Error On Mongo Connect", zap.Error(err))
	} else {
		err := mongoClient.Ping(nil, nil)
		if err != nil {
			log.Fatal("Error On Mongo Ping", zap.Error(err))
		}
		auth.InitMongo(mongoClient)
		crawler.InitMongo(mongoClient)
		music.InitMongo(mongoClient)
		session.InitMongo(mongoClient)
		store.InitMongo(mongoClient)
		token.InitMongo(mongoClient)
		user.InitMongo(mongoClient)
		vas.InitMongo(mongoClient)
	}

	// Initialize RedisCache
	redisConfig := redis.DefaultConfig
	redisConfig.Host = viper.GetString(config.RedisUrl)
	redisConfig.Password = viper.GetString(config.RedisPass)
	redisCache := redis.New(redisConfig)
	auth.InitRedisCache(redisCache)
	user.InitRedisCache(redisCache)

	acr.Init()
	auth.Init()
	crawler.Init()
	music.Init()
	saba.Init()
	session.Init()
	store.Init()
	token.Init()
	user.Init()
}

func initServer() *iris.Application {
	app := iris.New()

	adminParty := app.Party("/admin")
	adminParty.Use(auth.MustHaveAccessKey, auth.MustAdmin)
	adminParty.Post("/migrate_legacy_db", admin.MigrateLegacyDBHandler)
	adminParty.Get("/migrate_legacy_db_stats", admin.MigrateLegacyDBStatsHandler)

	tokenParty := app.Party("/token")
	tokenParty.Use(auth.MustHaveAccessKey)
	tokenParty.Post("/create", auth.MustWriteAccess, token.CreateHandler)
	tokenParty.Post("/validate", auth.MustReadAccess, token.ValidateHandler)

	authParty := app.Party("/auth")
	authParty.Use(auth.MustHaveAccessKey)
	authParty.Post("/create", auth.MustAdmin, auth.CreateAccessKeyHandler)
	authParty.Post("/send_code", auth.SendCodeHandler)
	authParty.Post("/login", auth.LoginHandler)
	authParty.Post("/register", auth.RegisterHandler)
	authParty.Post("/logout", session.MustHaveSession, auth.LogoutHandler)

	storeParty := app.Party("/store")
	storeParty.Use(auth.MustHaveAccessKey)
	storeParty.Post("/save", auth.MustAdmin, store.SaveHandler)
	storeParty.Get("/get", auth.MustAdmin, store.GetHandler)

	crawlerParty := app.Party("/crawler")
	crawlerParty.Use(auth.MustHaveAccessKey)
	crawlerParty.Post("/save", auth.MustAdmin, crawler.SaveHandler)

	musicParty := app.Party("/music")
	musicParty.Use(auth.MustHaveAccessKey)
	musicParty.Post("/search_by_proxy", session.MustHaveSession, user.MustVasEnabled, music.SearchByProxyHandler)
	musicParty.Post("/search_by_sound", session.MustHaveSession, user.MustVasEnabled, music.SearchBySoundHandler)
	musicParty.Post("/search_by_text", session.MustHaveSession, user.MustVasEnabled, music.SearchByTextHandler)
	musicParty.Post("/search_resume", session.MustHaveSession, user.MustVasEnabled, music.SearchByCursorHandler)
	musicParty.Get("/download/{bucket}/{downloadID}", session.MustHaveSession, user.MustVasEnabled, music.DownloadHandler)

	// Value Added Services
	vasParty := app.Party("/vas")
	vasParty.Get("/mci/notify", vas.MCINotification)
	vasParty.Get("/mci/mo", vas.MCIMo)

	// shopParty := app.Party("/shop")
	// shopParty.Post("/sep/")
	return app
}

func main() {
	// Initialize all the required modules and packages
	initModules()
	Root.AddCommand(InitDB)
	_ = Root.Execute()
}
