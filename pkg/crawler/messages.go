package crawler

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type SearchRequest struct {
	RequestID string `json:"req_id"`
	Keyword   string `json:"keyword"`
}

// easyjson:json
type SearchResponse struct {
	RequestID string `json:"req_id"`
	Sources   string `json:"source"`
	Result    struct {
		SongUrl  string `json:"song_url"`
		CoverUrl string `json:"cover_url"`
		Lyrics   string `json:"lyrics,omitempty"`
		Artists  string `json:"artists"`
		Title    string `json:"title"`
		Genre    string `json:"genre"`
	} `json:"result"`
}

