package store

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
// Store
type Store struct {
	ID       int64  `bson:"_id" json:"id"`
	Dsn      string `bson:"dsn" json:"dsn"`
	Capacity int    `bson:"cap" json:"capacity"`
	Region   string `bson:"region" json:"region"`
}

// GetStore returns a store identified by storeID
func get(storeID int64) *Store {
	storesMtx.RLock()
	s := stores[storeID]
	storesMtx.RUnlock()
	return s
}
