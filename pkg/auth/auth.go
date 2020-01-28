package auth

import (
	"git.ronaksoftware.com/blip/server/pkg/sms"
	ronak "git.ronaksoftware.com/ronak/toolbox"
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

