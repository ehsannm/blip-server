package acr

import "git.ronaksoftware.com/blip/server/pkg/config"

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	baseUrl      string
	accessKey    string
	accessSecret string
)

func Init() {
	baseUrl = config.GetString(config.ACRBaseUrl)
	accessKey = config.GetString(config.ACRAccessKey)
	accessSecret = config.GetString(config.ACRAccessSecret)
}
