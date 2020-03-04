package acr

/*
   Creation Time: 2019 - Oct - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type Music struct {
	Status struct {
		Message string `json:"msg"`
		Code    int    `json:"code"`
		Version string `json:"version"`
	} `json:"status"`
	Metadata struct {
		PlayedDuration int `json:"played_duration"`
		Music          []struct {
			ExternalIDs struct {
				Isrc string `json:"isrc"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			SampleBeginTimeOffsetMs string `json:"sample_begin_time_offset_ms"`
			Label                   string `json:"label"`
			ExternalMetadata        struct {
				Spotify struct {
					Album struct {
						ID string `json:"id"`
					} `json:"album"`
					Artists []struct {
						ID string `json:"id"`
					}
					Track struct {
						ID string `json:"id"`
					} `json:"track"`
				} `json:"spotify"`
			} `json:"external_metadata"`
			PlayOffsetMS int `json:"play_offset_ms"`
			Artists      []struct {
				Name string `json:"name"`
			} `json:"artists"`
			SampleEndTimeOffsetMS string `json:"sample_end_time_offset_ms"`
			ReleaseDate           string `json:"release_date"`
			Title                 string `json:"title"`
			DbEndTimeOffsetMS     string `json:"db_end_time_offset_ms"`
			DbBeginTimeOffsetMS   string `json:"db_begin_time_offset_ms"`
			DurationMS            int    `json:"duration_ms"`
			Album                 struct {
				Name string `json:"name"`
			} `json:"album"`
			Score int `json:"score"`
		} `json:"music"`
		TimestampUTC string `json:"timestamp_utc"`
	} `json:"metadata"`
	ResultType int `json:"result_type"`
}
