package auth

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type CreateAccessToken struct {
	Permissions []string `json:"permissions"`
	Period      int64    `json:"period"`
	AppName     string   `json:"app_name"`
}

const CAccessTokenCreated = "ACCESS_TOKEN_CREATED"

// easyjson:json
type AccessTokenCreated struct {
	AccessToken string `json:"access_token"`
	ExpireOn    int64  `json:"expire_on"`
}

// easyjson:json
// @RPC
// @Returns: PhoneCodeSent
type SendCodeReq struct {
	Phone string `json:"phone"`
}

const CPhoneCodeSent = "PHONE_CODE_SENT"

// easyjson:json
type PhoneCodeSent struct {
	PhoneCodeHash string `json:"phone_code_hash"`
	OperationID   string `json:"operation_id"`
	Registered    bool   `json:"registered"`
}

// easyjson:json
// @RPC
// @Returns: Authorization
type LoginReq struct {
	PhoneCode     string `json:"phone_code"`
	PhoneCodeHash string `json:"phone_code_hash"`
	Phone         string `json:"phone"`
	OperationID   string `json:"operation_id"`
}

// easyjson:json
// @RPC
// @Returns: Authorization
type RegisterReq struct {
	PhoneCode     string `json:"phone_code"`
	PhoneCodeHash string `json:"phone_code_hash"`
	Phone         string `json:"phone"`
	OperationID   string `json:"operation_id"`
	Username      string `json:"username"`
}

const CAuthorization = "AUTHORIZATION"

// easyjson:json
type Authorization struct {
	UserID    string `json:"user_id"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
}
