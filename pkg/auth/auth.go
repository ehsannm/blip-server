package auth

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/sms"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	authCol     *mongo.Collection
	redisCache  *ronak.RedisCache
	smsProvider sms.Provider
)

func InitMongo(c *mongo.Client) {
	authCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColAuth)
}

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

type Permission byte

const (
	_ Permission = 1 << iota
	Admin
	Read
	Write
)

type Auth struct {
	ID          string       `bson:"_id"`
	Permissions []Permission `bson:"perm"`
	CreatedOn   int64        `bson:"created_on"`
	ExpiredOn   int64        `bson:"expired_on"`
	AppName     string       `bson:"app_name"`
}

var authCache map[string]Auth
var mtxLock sync.RWMutex

func init() {
	authCache = make(map[string]Auth, 100000)
	smsProvider = sms.NewADP(
		config.GetString(config.SmsAdpUser),
		config.GetString(config.SmsAdpPass),
		config.GetString(config.SmsADPUrl),
		config.GetString(config.SmsAdpPhone),
		10,
	)
}
