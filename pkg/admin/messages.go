package admin

/*
   Creation Time: 2019 - Oct - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const CUnsubscribed = "UNSUBSCRIBED"

// easyjson:json
type Unsubscribed struct {
	Phone      string `json:"phone"`
	StatusCode string `json:"status_code"`
}

const CHealthCheckStats = "HEALTH_CHECK_STATS"

// easyjson:json
type HealthCheckStats struct {
	Scanned    int32 `json:"scanned"`
	CoverFixed int32 `json:"cover_fixed"`
	SongFixed  int32 `json:"song_fixed"`
}

// easyjson:json
// @Function
type SetVasReq struct {
	UserID  string `json:"user_id"`
	Enabled bool   `json:"enabled"`
}
