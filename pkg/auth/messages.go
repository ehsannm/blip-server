package auth

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const CAccessTokenCreated = "ACCESS_TOKEN_CREATED"

// easyjson:json
type AccessTokenCreated struct {
	AccessToken string `json:"access_token"`
	ExpireOn    int64  `json:"expire_on"`
}
