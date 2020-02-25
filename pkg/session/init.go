package session

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

/*
   Creation Time: 2020 - Jan - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	sessionCol      *mongo.Collection
	sessionCache    map[string]*Session
	sessionCacheMtx sync.RWMutex
)

func InitMongo(c *mongo.Client) {
	sessionCol = c.Database(config.DbMain).Collection(config.ColSession)
}

func Init() {
	sessionCache = make(map[string]*Session, 100000)
}
