package token

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
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
	tokenCol *mongo.Collection
)

type Token struct {
	ID        string `bson:"_id"`
	Phone     string `bson:"phone"`
	Period    int64  `bson:"period"`
	CreatedOn int64  `bson:"created_on"`
	ExpiredOn int64  `bson:"expired_on"`
	DeviceID  string `bson:"device_id"`
}

func InitMongo(c *mongo.Client) {
	tokenCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColToken)
}

var mtxLock sync.RWMutex
var tokenCache map[string]Token

func init() {
	tokenCache = make(map[string]Token, 100000)
}
