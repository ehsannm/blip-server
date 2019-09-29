package config

import "time"

/*
   Creation Time: 2019 - Sep - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	ColAuth    = "auth"
	ColToken   = "token"
	ColUser    = "user"
	ColSession = "session"
)

const (
	RkPhoneCode 	= "PHONE_CODE"
)

var (
	RegionCode = "R01"
	MongoRequestTimeout = time.Second * 3
)
