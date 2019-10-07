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
	HdrAccessKey = "AccessKey"
	CtxAuth      = "Auth"
	UsernameRegex = "/^[a-zA-Z0-9]+[a-zA-Z0-9_]{3,12}$/"
)

var (
	usernameREGX	*regexp.Regexp
)

func init() {
	usernameREGX, _ = regexp.Compile(UsernameRegex)
}