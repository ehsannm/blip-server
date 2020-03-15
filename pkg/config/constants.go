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

// MongoDB Collections
const (
	DbMain      = "blip"
	DbStore     = "blip-store"
	ColAuth     = "auth"
	ColToken    = "token"
	ColUser     = "user"
	ColSession  = "session"
	ColLogVas   = "log.vas"
	ColLogAuth  = "log.auth"
	ColSong     = "song"
	ColCrawler  = "crawler"
	ColStore    = "store"
	ColHelp     = "help"
	ColFeedback = "feedback"
	ColDevice   = "device"
)

// Redis Keys
const (
	RkPhoneCode = "PHONE_CODE"
)

var (
	RegionCode          = "R01"
	MongoRequestTimeout = time.Second * 3
	HttpRequestTimeout  = 30 * time.Second
	ServerID            = "BLIP-01"
)
