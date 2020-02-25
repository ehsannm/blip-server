package auth

import "regexp"

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	// Header
	HdrAccessKey = "AccessKey"

	// Context
	CtxAuth       = "Auth"
	CtxClientName = "ClientName"

	// Apps
	AppNameAdmin    = "Admin"
	AppNameMusicChi = "MusicChi"
	AppNameBlip     = "Blip"

	UsernameRegex = "^[a-zA-Z0-9]+[a-zA-Z0-9_]{3,12}$"
)

var (
	usernameREGX *regexp.Regexp
)

func init() {
	usernameREGX, _ = regexp.Compile(UsernameRegex)
}

var supportedCarriers = map[string]bool{
	"98910": true,
	"98911": true,
	"98912": true,
	"98913": true,
	"98914": true,
	"98915": true,
	"98916": true,
	"98917": true,
	"98918": true,
	"98919": true,
	"98990": true,
	"98991": true,
	"98992": true,
}
