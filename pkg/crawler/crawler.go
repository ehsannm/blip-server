package crawler

import (
	"bytes"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"sync"
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
	redisCache             *ronak.RedisCache
	crawlerCol             *mongo.Collection
	registeredCrawlersMtx  sync.RWMutex
	registeredCrawlers     map[string][]*Crawler
	registeredCrawlersPool sync.Pool
)

// DropAll deletes all the crawlers from the database
func DropAll() error {
	return crawlerCol.Drop(nil)
}

// Save insert the crawler 'c' into the database
func Save(c Crawler) (primitive.ObjectID, error) {
	res, err := crawlerCol.InsertOne(nil, c, options.InsertOne())
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), err
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

// Crawler
type Crawler struct {
	httpClient  http.Client `bson:"-"`
	Url         string      `bson:"url"`
	Name        string      `bson:"name"`
	Description string      `bson:"desc"`
	Source      string      `bson:"source"`
}

func (c *Crawler) SendRequest(reqID string, keyword string) (*SearchResponse, error) {
	c.httpClient.Timeout = config.HttpRequestTimeout
	req := SearchRequest{
		Keyword: keyword,
	}
	reqBytes, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}
	httpRes, err := c.httpClient.Post(c.Url, "application/json", bytes.NewBuffer(reqBytes))
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
