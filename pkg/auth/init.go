package auth

import (
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/sms"
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

//go:generate rm -f *_easyjson.go
//go:generate easyjson messages.go
func init() {
	authCache = make(map[string]Auth, 100000)
	smsProvider = sms.NewPayamak(
		config.GetString(config.SmsPayamakUser),
		config.GetString(config.SmsPayamakPass),
		config.GetString(config.SmsPayamakUrl),
		config.GetString(config.SmsPayamakPhone),
		10,
	)
}

func InitMongo(c *mongo.Client) {
	authCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColAuth)
}

func InitRedisCache(c *redis.Cache) {
	redisCache = c
}
