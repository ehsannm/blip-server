package crawler

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
)

/*
   Creation Time: 2019 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/
var (
	redisCache *ronak.RedisCache
)

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func getNextRequestID() string {
	return fmt.Sprintf("%s.%s", config.ServerID, ronak.RandomID(32))
}

type searchRequest struct {
	Keyword string
}

func SendSearchRequest(keyword string) {
	// c := http.Client{
	// 	Timeout: shared.HttpRequestTimeout,
	// }
	// reqID := getNextRequestID()
	// redisCache.Do(radix.FlatCmd(nil, "SET", fmt.Sprintf("CR-%s", getNextRequestID())))

}
