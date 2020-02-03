package dev

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

// MigrateLegacyDB
type MigrateLegacyDBReq struct {
	MysqlHost string `json:"mysql_host"`
	MysqlUser string `json:"mysql_user"`
	MysqlPass string `json:"mysql_pass"`
}
