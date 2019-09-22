package auth

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
	authCol *mongo.Collection
)

func InitMongo(c *mongo.Client) {
	authCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColAuth)
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
}

var authCache map[string]Auth
var mtxLock sync.RWMutex

func init() {
	authCache = make(map[string]Auth, 100000)
}
