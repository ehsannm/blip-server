package crawler

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	RequestID string      `json:"req_id"`
	Source    string      `json:"source"`
	Result    []FoundSong `json:"result"`
}

// easyjson:json
type FoundSong struct {
	SongUrl     string `json:"song_url"`
	CoverUrl    string `json:"cover_url"`
	Lyrics      string `json:"lyrics,omitempty"`
	Artists     string `json:"artists"`
	Title       string `json:"title"`
	Genre       string `json:"genre,omitempty"`
	UrlLifetime int    `json:"url_lifetime"`
}

// easyjson:json
type SaveReq struct {
	ID             primitive.ObjectID `json:"id"`
	Url            string             `json:"url"`
	Source         string             `json:"source"`
	Name           string             `json:"name"`
	DownloaderJobs int                `json:"downloader_jobs"`
}
