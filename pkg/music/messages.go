package music

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
// @Returns SearchResult
type SearchReq struct {
	Keyword string `json:"keyword"`
	Artist  string `json:"artist,omitempty"`
	Lyrics  string `json:"lyrics,omitempty"`
}

const CSearchResult = "SEARCH_RESULT"

// easyjson:json
// SearchResult
type SearchResult struct {
	Songs []*Song `json:"songs"`
}

const CSaveStore = "SAVE_STORE"

// easyjson:json
// SaveStoreReq
type SaveStoreReq struct {
	StoreID  int64  `json:"store_id"`
	Dsn      string `json:"dsn"`
	Region   string `json:"region"`
	Capacity int    `json:"capacity"`
}
