package store

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate rm -f *_easyjson.go
//go:generate easyjson store.go messages.go
var (
	storeCol     *mongo.Collection
	stores       map[int64]*Store
	storeBuckets map[int64]*mongo.Client
	storesMtx    sync.RWMutex
)

func InitMongo(c *mongo.Client) {
	storeCol = c.Database(config.Db).Collection(config.ColStore)
}

func Init() {
	storesMtx.Lock()
	defer storesMtx.Unlock()
	stores = make(map[int64]*Store)
	storeBuckets = make(map[int64]*mongo.Client)
	cur, err := storeCol.Find(nil, bson.D{})
	if err != nil {
		log.Warn("Error On Initializing Stores", zap.Error(err))
	}
	for cur.Next(nil) {
		storeX := &Store{}
		err = cur.Decode(storeX)
		if err != nil {
			continue
		}
		if err := createStoreConnection(storeX); err != nil {
			log.Warn("Err On Create Store Connection",
				zap.Int64("StoreID", storeX.ID),
				zap.String("Dsn", storeX.Dsn),
				zap.Error(err),
			)
			continue
		}
	}
	go watchForStores()
}
func createStoreConnection(storeX *Store) error {
	return tools.Try(5, time.Second, func() error {
		mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(storeX.Dsn))
		if err != nil {
			return err
		}
		err = mongoClient.Ping(nil, nil)
		if err != nil {
			return err
		}
		storeBuckets[storeX.ID] = mongoClient
		stores[storeX.ID] = storeX
		return nil
	})
}
func watchForStores() {
	var resumeToken bson.Raw
	for {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		if resumeToken != nil {
			opts.SetStartAfter(resumeToken)
		}
		stream, err := storeCol.Watch(nil, mongo.Pipeline{}, opts)
		if err != nil {
			log.Warn("Error On Watch Stream for Stores", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		for stream.Next(nil) {
			storeX := &Store{}
			resumeToken = stream.ResumeToken()
			operationType := strings.Trim(stream.Current.Lookup("operationType").String(), "\"")
			switch operationType {
			case "insert", "update":
				err = stream.Current.Lookup("fullDocument").UnmarshalWithRegistry(bson.DefaultRegistry, storeX)
				if err != nil {
					log.Warn("Error On Decoding Store", zap.Error(err))
					continue
				}
				storesMtx.Lock()
				if err := createStoreConnection(storeX); err != nil {
					log.Warn("Err On Create Store Connection",
						zap.Int64("StoreID", storeX.ID),
						zap.String("Dsn", storeX.Dsn),
						zap.Error(err),
					)
					continue
				}
				storesMtx.Unlock()
			}
		}
		_ = stream.Close(nil)
	}

}
