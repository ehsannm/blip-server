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
type SearchResult struct {
	Songs []*Song `json:"songs"`
}

const CSoundSearchResult = "SOUND_SEARCH_RESULT"

// easyjson:json
type SoundSearchResult struct {
	Info struct {
		Artists     []string `json:"artists"`
		Title       string   `json:"title"`
		ReleaseDate string   `json:"release_date"`
		ImageUrl    string   `json:"image_url"`
	} `json:"info"`
	Songs []*Song `json:"songs"`
}
