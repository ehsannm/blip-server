package saba

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var SuccessfulCode = "SC000"

var SabaCodes = map[string]string{
	"SC000": "successful",
	"SC001": "service not found",
	"SC002": "service token is invalid",
	"SC003": "operator not found",
	"SC004": "cellphone is blocked",
	"SC005": "transaction_id is invald",
	"SC006": "authentication failed",
	"SC007": "phone number is invalid",
	"SC008": "authorization failed",
	"SC009": "pin already sent to destination",
	"SC010": "TPS exceeded",
	"SC011": "whitelist error",
	"SC012": "subscription already exists",
	"SC013": "subscription not exist",
	"SC014": "method is unavailable",
	"SC015": "service exception",
	"SC016": "waiting for callback",
	"SC017": "service is inactive",
	"SC018": "service_id is invalid",
}
