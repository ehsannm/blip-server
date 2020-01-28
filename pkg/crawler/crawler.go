package crawler

import (
	"bytes"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/mediocregopher/radix/v3"
	"net/http"
)

/*
   Creation Time: 2019 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate easyjson -pkg
var (
	redisCache *ronak.RedisCache
	crawlerUrl string
)

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func Init() {
	crawlerUrl = config.GetString(config.CrawlerUrl)
}
func getNextRequestID() string {
	return fmt.Sprintf("%s.%s", config.ServerID, ronak.RandomID(32))
}

// easyjson:json
type searchRequest struct {
	Keyword string `json:"keyword"`
}

func SendSearchRequest(keyword string) error {
	c := http.Client{
		Timeout: config.HttpRequestTimeout,
	}
	reqID := getNextRequestID()
	redisCache.Do(radix.FlatCmd(nil, "SET", fmt.Sprintf("CR-%s", reqID)))

	req := searchRequest{
		Keyword: keyword,
	}
	reqBytes, err := req.MarshalJSON()
	if err != nil {
		return err
	}
	httpRes, err := c.Post(crawlerUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusAccepted {
		return fmt.Errorf("invalid http response, got %d", httpRes.StatusCode)
	}
	return nil
}
