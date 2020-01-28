package crawler

import (
	"context"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	crawlerCol *mongo.Collection
)

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

func InitMongo(c *mongo.Client) {
	crawlerCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColCrawler)
}

// Crawler
type Crawler struct {
	Url         string `bson:"url"`
	Name        string `bson:"name"`
	Description string `bson:"desc"`
	Source      string `bson:"source"`
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

// GetAll returns all the crawlers from the database
func GetAll() ([]*Crawler, error) {
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
		crawlers = append(crawlers, crawler)
	}
	return crawlers, nil
}

func DropAll() error {
	return crawlerCol.Drop(nil)
}
