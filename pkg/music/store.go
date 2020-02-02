package music

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

// GetStore returns a store identified by storeID
func GetStore(storeID int64) *Store {
	storesMtx.RLock()
	s := stores[storeID]
	storesMtx.RUnlock()
	return s
}
