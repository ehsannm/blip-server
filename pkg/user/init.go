package user

import (
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   Creation Time: 2020 - Jan - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate easyjson user.go

func InitMongo(c *mongo.Client) {
	userCol = c.Database(config.DbMain).Collection(config.ColUser)
}

func InitRedisCache(c *redis.Cache) {
	redisCache = c
}

func Init() {
	_, _ = userCol.InsertOne(nil, User{
		ID:        "MAGIC_USER",
		Username:  "MAGIC_USER",
		Phone:     "2374002374",
		Email:     "support@blip.fun",
		CreatedOn: 0,
		Disabled:  false,
		VasPaid:   true,
	})
}
