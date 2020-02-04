package crawler

import (
	"context"
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/pools"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.uber.org/zap"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
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

// DropAll deletes all the crawlers from the database
func DropAll() error {
	return crawlerCol.Drop(nil)
}

// Save insert the crawler 'c' into the database
func Save(crawlerX *Crawler) (primitive.ObjectID, error) {
	if crawlerX.ID == primitive.NilObjectID {
		crawlerX.ID = primitive.NewObjectID()
	}
	_, err := crawlerCol.UpdateOne(nil, bson.M{"_id": crawlerX.ID}, bson.M{"$set": crawlerX}, options.Update().SetUpsert(true))
	if err != nil {
		return primitive.NilObjectID, err
	}
	return crawlerX.ID, err
}

// Get returns a crawler identified by 'crawlerID'
func Get(crawlerID primitive.ObjectID) (*Crawler, error) {
	res := crawlerCol.FindOne(nil, bson.M{"_id": crawlerID}, options.FindOne())
	c := &Crawler{}
	err := res.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetAll returns an array of the registered crawlers
func GetAll() []*Crawler {
	list := make([]*Crawler, 0, 10)
	registeredCrawlersMtx.RLock()
	for _, crawlers := range registeredCrawlers {
		list = append(list, crawlers...)
	}
	registeredCrawlersMtx.RUnlock()
	return list
}

// Search sends search request to all the crawlers and pushes the result into the channel.
// Crawlers are categorized with their 'Source' tag.
func Search(ctx context.Context, keyword string) <-chan *SearchResponse {
	crawlers := getRegisteredCrawlers()
	if len(crawlers) == 0 {
		log.Warn("No Crawler has been Registered")
		return nil
	}
	resChan := make(chan *SearchResponse, len(crawlers))

	go func() {
		defer putRegisteredCrawlers(crawlers)
		waitGroup := pools.AcquireWaitGroup()
		for _, c := range crawlers {
			waitGroup.Add(1)
			go func(c *Crawler) {
				defer waitGroup.Done()
				res, err := c.SendRequest(ctx, keyword)
				if err != nil {
					log.Warn("Error On Crawler Request",
						zap.Error(err),
						zap.String("CrawlerUrl", c.Url),
						zap.String("Keyword", keyword),
					)
					return
				}
				resChan <- res
			}(c)
		}
		waitGroup.Wait()
		pools.ReleaseWaitGroup(waitGroup)
		close(resChan)
	}()
	return resChan
}
func getRegisteredCrawlers() []*Crawler {
	list, ok := registeredCrawlersPool.Get().([]*Crawler)
	if ok {
		return list
	}
	list = make([]*Crawler, 0, len(registeredCrawlers))
	registeredCrawlersMtx.RLock()
	for _, crawlers := range registeredCrawlers {
		idx := tools.RandomInt(len(crawlers))
		list = append(list, crawlers[idx])
	}
	registeredCrawlersMtx.RUnlock()
	return list
}
func putRegisteredCrawlers(list []*Crawler) {
	registeredCrawlersPool.Put(list)
}

// Crawler
type Crawler struct {
	httpClient     http.Client        `bson:"-"`
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Url            string             `bson:"url" json:"url"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"desc" json:"description"`
	Source         string             `bson:"source" json:"source"`
	DownloaderJobs int                `bson:"downloader_jobs" json:"downloader_jobs"`
}

func (c *Crawler) SendRequest(ctx context.Context, keyword string) (*SearchResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	c.httpClient.Timeout = config.HttpRequestTimeout

	crawlerUrl := fmt.Sprintf("%s/%s/%s", c.Url, tools.RandomID(24), url.QueryEscape(keyword))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, crawlerUrl, nil)
	if err != nil {
		return nil, err
	}
	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("invalid http response, got %d", httpRes.StatusCode)
	}
	resBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	res := &SearchResponse{}
	err = res.UnmarshalJSON(resBytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}
