package main

import (
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/sms/saba"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas"
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
		token.InitMongo(mongoClient)
		session.InitMongo(mongoClient)
		user.InitMongo(mongoClient)
	}

	// Initialize RedisCache
	redisConfig := ronak.DefaultRedisConfig
	redisConfig.Host = viper.GetString(config.ConfRedisUrl)
	redisConfig.Password = viper.GetString(config.ConfRedisPass)
	redisCache := ronak.NewRedisCache(redisConfig)
	auth.InitRedisCache(redisCache)

	saba.Init()
}

func initServer() *iris.Application {
	app := iris.New()

	tokenParty := app.Party("/token")
	tokenParty.Use(auth.GetAuthorizationHandler)
	tokenParty.Post("/create", auth.MustWriteAccess, token.CreateHandler)
	tokenParty.Post("/validate", auth.MustReadAccess, token.ValidateHandler)

	authParty := app.Party("/auth")
	authParty.Use(auth.GetAuthorizationHandler)
	authParty.Post("/create", auth.MustAdmin, auth.CreateAccessKeyHandler)
	authParty.Post("/send_code", auth.SendCodeHandler)
	authParty.Post("/login", auth.LoginHandler)
	authParty.Post("/register", auth.RegisterHandler)

	musicParty := app.Party("/music")
	musicParty.Use(auth.GetAuthorizationHandler)
	musicParty.Get("/search", music.Search)

	// Value Added Services
	vasParty := app.Party("/vas")
	vasParty.Get("/mci/notify", vas.MCINotification)
	vasParty.Get("/mci/mo", vas.MCIMo)

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
