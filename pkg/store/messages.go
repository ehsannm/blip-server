package store

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const CSaveStore = "SAVE_STORE"

// easyjson:json
// SaveStoreReq
type SaveStoreReq struct {
	StoreID  int64  `json:"store_id"`
	Dsn      string `json:"dsn"`
	Region   string `json:"region"`
	Capacity int    `json:"capacity"`
}

// easyjson:json
// GetStoreReq
type GetStoreReq struct {
	StoreIDs []int64 `json:"store_ids"`
}

const CStores = "STORES"

// easyjson:json
// Stores
type Stores struct {
	Stores []*Store `json:"stores"`
}
