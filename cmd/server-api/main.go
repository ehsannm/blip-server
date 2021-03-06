package main

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/admin"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"git.ronaksoftware.com/blip/server/pkg/device"
	"git.ronaksoftware.com/blip/server/pkg/help"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/store"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	"github.com/kataras/iris/v12"
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
		options.Client().
			ApplyURI(config.GetString(config.MongoUrl)).
			SetMaxPoolSize(200).
			SetMinPoolSize(1),
	); err != nil {
		log.Fatal("Error On Mongo Connect", zap.Error(err))
	} else {
		err := mongoClient.Ping(nil, nil)
		if err != nil {
			log.Fatal("Error On Mongo Ping", zap.Error(err))
		}
		auth.InitMongo(mongoClient)
		crawler.InitMongo(mongoClient)
		device.InitMongo(mongoClient)
		help.InitMongo(mongoClient)
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
	help.Init()
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
	adminParty.Post("/health_check_db", admin.HealthCheckDbHandler)
	adminParty.Post("/health_check_store", admin.HealthCheckStoreHandler)
	adminParty.Get("/health_check_stats", admin.HealthCheckStatsHandler)
	adminParty.Post("/vas", admin.SetVas)

	deviceParty := app.Party("/device")
	deviceParty.Use(auth.MustHaveAccessKey, session.MustHaveSession)
	deviceParty.Post("/register", device.RegisterDevice)

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

	helpParty := app.Party("/help")
	helpParty.Use(auth.MustHaveAccessKey)
	helpParty.Post("/config", auth.MustAdmin, help.SetHandler)
	helpParty.Delete("/config", auth.MustAdmin, help.UnsetHandler)
	helpParty.Get("/config", help.GetHandler)
	helpParty.Post("/feedback", session.MustHaveSession, help.FeedbackHandler)

	storeParty := app.Party("/store")
	storeParty.Use(auth.MustHaveAccessKey)
	storeParty.Post("/save", auth.MustAdmin, store.SaveHandler)
	storeParty.Get("/get", auth.MustAdmin, store.GetHandler)

	crawlerParty := app.Party("/crawler")
	crawlerParty.Use(auth.MustHaveAccessKey)
	crawlerParty.Post("/save", auth.MustAdmin, crawler.SaveHandler)
	crawlerParty.Get("/list", auth.MustAdmin, crawler.ListHandler)
	crawlerParty.Delete("/{crawlerID}", auth.MustAdmin, crawler.RemoveHandler)

	musicParty := app.Party("/music")
	musicParty.Use(auth.MustHaveAccessKey)
	musicParty.Post("/search_by_proxy", session.MustHaveSession, user.MustVasEnabled, music.SearchByProxyHandler)
	musicParty.Post("/search/sound", session.MustHaveSession, user.MustVasEnabled, music.SearchBySoundHandler)
	musicParty.Post("/search/fingerprint", session.MustHaveSession, user.MustVasEnabled, music.SearchByFingerprintHandler)
	musicParty.Post("/search/text", session.MustHaveSession, user.MustVasEnabled, music.SearchByTextHandler)
	musicParty.Post("/search/bot", session.MustHaveSession, user.MustVasEnabled, music.SearchByBotHandler)
	musicParty.Get("/search", session.MustHaveSession, user.MustVasEnabled, music.SearchByCursorHandler)
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

func debugMiddleware(ctx iris.Context) {
	if ce := log.Check(log.DebugLevel, "Request Received"); ce != nil {
		ce.Write(
			zap.String("Url", ctx.Request().RequestURI),
			zap.String("Method", ctx.Request().Method),
			zap.Bool("HashSession", ctx.Request().Header.Get(session.HdrSessionID) != ""),
		)
	}
}
