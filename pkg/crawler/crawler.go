package crawler

import (
	"bytes"
	"context"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/mediocregopher/radix/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetAll returns all the crawlers from the database
func GetAll() ([]*Crawler, error) {
	registeredCrawlersMtx.Lock()
	defer registeredCrawlersMtx.Unlock()
	registeredCrawlers = make(map[string][]*Crawler)
	cur, err := crawlerCol.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	crawlers := make([]*Crawler, 0, 10)
	for cur.Next(ctx) {
		crawler := &Crawler{}
		err = cur.Decode(crawler)
		if err != nil {
			return crawlers, err
		}
		registeredCrawlers[crawler.Source] = append(registeredCrawlers[crawler.Source], crawler)
		crawlers = append(crawlers, crawler)
	}
	return crawlers, nil
}

func DropAll() error {
	return crawlerCol.Drop(nil)
}


// Crawler
type Crawler struct {
	httpClient  http.Client `bson:"-"`
	Url         string      `bson:"url"`
	Name        string      `bson:"name"`
	Description string      `bson:"desc"`
	Source      string      `bson:"source"`
}

func (c *Crawler) SendRequest(keyword string) error {
	c.httpClient.Timeout = config.HttpRequestTimeout
	reqID := getNextRequestID()
	err := redisCache.Do(radix.FlatCmd(nil, "SET", fmt.Sprintf("CR-%s", reqID)))
	if err != nil {
		return err
	}
	req := searchRequest{
		Keyword: keyword,
	}
	reqBytes, err := req.MarshalJSON()
	if err != nil {
		return err
	}
	httpRes, err := c.httpClient.Post(c.Url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusAccepted {
		return fmt.Errorf("invalid http response, got %d", httpRes.StatusCode)
	}
	return nil
}
