package music

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Store
type Store struct {
	ID       int64  `bson:"_id"`
	Dsn      string `bson:"dsn"`
	Capacity int    `bson:"cap"`
	Region   string `bson:"region"`
}

// SaveStore
func SaveStore(storeX *Store) error {
	_, err := storeCol.UpdateOne(nil, bson.M{"_id": storeX.ID}, bson.M{"$set": storeX}, options.Update().SetUpsert(true))
	return err
}

// GetStore returns a store identified by storeID
func GetStore(storeID int64) *Store {
	storesMtx.RLock()
	s := stores[storeID]
	storesMtx.RUnlock()
	return s
}
