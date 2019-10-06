package main

import (
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"git.ronaksoftware.com/blip/server/pkg/user"
	ronak "git.ronaksoftware.com/ronak/toolbox"

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
	log.InitLogger(log.DebugLevel, "")


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
}

func initServer() *iris.Application {
	app := iris.New()
	app.UseGlobal(auth.GetAuthorizationHandler)

	tokenParty := app.Party("/token")
	tokenParty.Post("/create", auth.MustWriteAccess, token.CreateHandler)
	tokenParty.Post("/validate", auth.MustReadAccess, token.ValidateHandler)

	authParty := app.Party("/auth")
	authParty.Post("/create", auth.MustAdmin, auth.CreateAccessKeyHandler)
	authParty.Post("/send_code", auth.SendCodeHandler)
	authParty.Post("/login", auth.LoginHandler)
	authParty.Post("/register", auth.RegisterHandler)

	musicParty := app.Party("/music")
	musicParty.Get("/search", music.Search)


	vasParty := app.Party("/vas")
	vasParty.Get("/notif", auth.VasNotification)



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
