package token

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	tokenCol   *mongo.Collection
	mtxLock    sync.RWMutex
	tokenCache map[string]Token
)

// Token
type Token struct {
	ID        string `bson:"_id"`
	Phone     string `bson:"phone"`
	Period    int64  `bson:"period"`
	CreatedOn int64  `bson:"created_on"`
	ExpiredOn int64  `bson:"expired_on"`
	DeviceID  string `bson:"device_id"`
}
