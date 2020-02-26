package auth

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/sms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
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
var (
	authCol      *mongo.Collection
	authCache    map[string]*Auth
	authCacheMtx sync.RWMutex
	redisCache   *redis.Cache
	smsProvider  sms.Provider
)

func InitMongo(c *mongo.Client) {
	authCol = c.Database(config.DbMain).Collection(config.ColAuth)
}

func InitRedisCache(c *redis.Cache) {
	redisCache = c
}

func Init() {
	authCache = make(map[string]*Auth, 10000)
	smsProvider = sms.NewPayamak(
		config.GetString(config.SmsPayamakUser),
		config.GetString(config.SmsPayamakPass),
		config.GetString(config.SmsPayamakUrl),
		config.GetString(config.SmsPayamakPhone),
		100,
	)

	cnt, err := authCol.CountDocuments(nil, bson.D{})
	if err != nil {
		fmt.Println(err)
		return
	}
	if cnt == 0 {
		_, err := authCol.InsertOne(nil, Auth{
			ID:          "ROOT",
			Permissions: []Permission{Admin},
			CreatedOn:   time.Now().Unix(),
			ExpiredOn:   0,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
