package main

import (
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"go.uber.org/zap/zapcore"

	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	_Mongo *mongo.Client
)

func init() {
	log.InitLogger(zapcore.Level(config.GetInt(config.LogLevel)), "")

	// Initialize MongoDB
	if mongoClient, err := mongo.Connect(
		nil,
		options.Client().ApplyURI(viper.GetString(config.MongoUrl)),
	); err != nil {
		log.Fatal("Error On MongoConnect", zap.Error(err))
	} else {
		_Mongo = mongoClient
		auth.InitMongo(mongoClient)
		session.InitMongo(mongoClient)
		token.InitMongo(mongoClient)
		user.InitMongo(mongoClient)
		vas.InitMongo(mongoClient)
	}

	// Initialize RedisCache
	redisConfig := ronak.DefaultRedisConfig
	redisConfig.Host = viper.GetString(config.RedisUrl)
	redisConfig.Password = viper.GetString(config.RedisPass)
	redisCache := ronak.NewRedisCache(redisConfig)
	auth.InitRedisCache(redisCache)
	user.InitRedisCache(redisCache)

	// Initialize VAS Saba Service
	saba.Init()

	// Initialize ACR Sound Identification Service
	acr.Init()
}

func initServer() *iris.Application {
	app := iris.New()

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

	musicParty := app.Party("/music")
	musicParty.Use(auth.MustHaveAccessKey)
	musicParty.Post("/search_by_proxy", session.MustHaveSession, user.MustVasEnabled, music.SearchByProxy)
	musicParty.Post("/search_by_sound", session.MustHaveSession, user.MustVasEnabled, music.SearchBySound)
	musicParty.Post("/search_by_text", session.MustHaveSession, user.MustVasEnabled, music.SearchByText)
	musicParty.Post("/search_by_cursor", session.MustHaveSession, user.MustVasEnabled, music.SearchByCursor)
	musicParty.Get("/download", session.MustHaveSession, user.MustVasEnabled, music.Download)

	// Value Added Services
	vasParty := app.Party("/vas")
	vasParty.Get("/mci/notify", vas.MCINotification)
	vasParty.Get("/mci/mo", vas.MCIMo)

	crawlerParty := app.Party("/crawler")
	crawlerParty.Post("/search_result", crawler.SearchResult)

	// shopParty := app.Party("/shop")
	// shopParty.Post("/sep/")
	return app
}

func main() {
	Root.AddCommand(InitDB)
	_ = Root.Execute()
}
