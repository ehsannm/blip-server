package crawler

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/mediocregopher/radix/v3"
	"sync"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type searchRequest struct {
	RequestID string `json:"req_id"`
	Keyword   string `json:"keyword"`
}

// easyjson:json
type searchResponse struct {
	RequestID string `json:"req_id"`
	Sources   string `json:"source"`
	Result    struct {
		SongUrl  string `json:"song_url"`
		CoverUrl string `json:"cover_url"`
		Lyrics   string `json:"lyrics,omitempty"`
		Artists  string `json:"artists"`
		Title    string `json:"title"`
		Genre    string `json:"genre"`
	} `json:"result"`
}

func getNextRequestID() string {
	return fmt.Sprintf("%s.%s", config.ServerID, ronak.RandomID(32))
}
func getRegisteredCrawlers() []*Crawler {
	list, ok := registeredCrawlersPool.Get().([]*Crawler)
	if ok {
		return list
	}
	list = make([]*Crawler, 0, len(registeredCrawlers))
	registeredCrawlersMtx.RLock()
	for _, crawlers := range registeredCrawlers {
		idx := ronak.RandomInt(len(crawlers))
		list = append(list, crawlers[idx])
	}
	registeredCrawlersMtx.RUnlock()
	return list
}
func putRegisteredCrawlers(list []*Crawler) {
	registeredCrawlersPool.Put(list)
}

func Search(keyword string) (string, error) {
	reqID := getNextRequestID()
	crawlers := getRegisteredCrawlers()
	defer putRegisteredCrawlers(crawlers)

	waitGroup := sync.WaitGroup{}

	for _, c := range crawlers {
		waitGroup.Add(1)
		go func(c *Crawler) {
			_ = c.SendRequest(keyword)
			waitGroup.Done()
		}(c)
	}
	redisCache.Do(radix.FlatCmd(nil, "HMSET", reqID, map[string]interface{}{}))

	return reqID, nil
}
