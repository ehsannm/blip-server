package device

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type RegisterDeviceReq struct {
	TokenType string `json:"token_type"` // apn || fb
	Token     string `json:"token"`
}

const CRegisterDevice = "REGISTER_DEVICE"
