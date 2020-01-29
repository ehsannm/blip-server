package redis_test

import (
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"testing"
)

/*
   Creation Time: 2020 - Jan - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/


func TestPubSub(t *testing.T) {
	redisConfig := redis.DefaultConfig
	redisConfig.Host = config.GetString(config.RedisUrl)
	redisConfig.Password = config.GetString(config.RedisPass)

	_ = redis.NewPubSub(redisConfig)

}