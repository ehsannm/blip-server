package main

import (
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	_Mongo *mongo.Client
)

func init() {
	log.InitLogger(zapcore.Level(config.GetInt(config.ConfLogLevel)), "")

	// Initialize MongoDB
	if mongoClient, err := mongo.Connect(
		nil,
		options.Client().ApplyURI(viper.GetString(config.ConfMongoUrl)),
	); err != nil {
		log.Fatal("Error On MongoConnect", zap.Error(err))
	} else {
		_Mongo = mongoClient
		auth.InitMongo(mongoClient)
		session.InitMongo(mongoClient)
		token.InitMongo(mongoClient)
		user.InitMongo(mongoClient)
	}

	// Initialize RedisCache
	redisConfig := ronak.DefaultRedisConfig
	redisConfig.Host = viper.GetString(config.ConfRedisUrl)
	redisConfig.Password = viper.GetString(config.ConfRedisPass)
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

	// Value Added Services
	vasParty := app.Party("/vas")
	vasParty.Get("/mci/notify", vas.MCINotification)
	vasParty.Get("/mci/mo", vas.MCIMo)

	// devParty := app.Party("/dev")
	// devParty.Post("/unsubscribe", dev.Unsubscribe)

	// shopParty := app.Party("/shop")
	// shopParty.Post("/sep/")
	return app
}

func main() {
	Root.AddCommand(InitDB)
	_ = Root.Execute()
}

var Root = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		app := initServer()
		err := app.Run(iris.Addr(":80"), iris.WithOptimizations)
		if err != nil {
			log.Warn(err.Error())
		}
	},
}
