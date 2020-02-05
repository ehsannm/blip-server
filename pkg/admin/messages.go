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

const CMigrateStats = "MIGRATE_STATS"

// easyjson:json
type MigrateStats struct {
	Scanned           int32 `json:"scanned"`
	Downloaded        int32 `json:"downloaded"`
	AlreadyDownloaded int32 `json:"already_downloaded"`
}
