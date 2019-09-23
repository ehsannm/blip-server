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

// easyjson:json
type SendCodeReq struct {
	Phone string `json:"phone"`
}

const CPhoneCodeSent = "PHONE_CODE_SENT"

// easyjson:json
type PhoneCodeSent struct {
	PhoneCodeHash string `json:"phone_code_hash"`
}

// easyjson:json
type LoginReq struct {
	PhoneCode     string `json:"phone_code"`
	PhoneCodeHash string `json:"phone_code_hash"`
	Phone         string `json:"phone"`
}

// easyjson:json
type RegisterReq struct {
	PhoneCode     string `json:"phone_code"`
	PhoneCodeHash string `json:"phone_code_hash"`
	Phone         string `json:"phone"`
	Username      string `json:"username"`
}

const CAuthorization = "AUTHORIZATION"

// easyjson:json
type Authorization struct {
	UserID   int32  `json:"user_id"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}
