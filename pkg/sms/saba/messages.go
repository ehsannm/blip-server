package saba

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type SubscribeResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"status_code"`
	OtpID            string `json:"otp_id"`
	OperatorResponse struct {
		StatusInfo struct {
			ReferenceCode       string `json:"referenceCode"`
			StatusCode          string `json:"statusCode"`
			ServerReferenceCode string `json:"serverReferenceCode"`
			OtpTransactionID    string `json:"otpTransactionId"`
		} `json:"statusInfo"`
	} `json:"operator_response"`
}

// easyjson:json
type UnsubscribeResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"status_code"`
	OperatorResponse struct {
		StatusInfo struct {
			ReferenceCode       string `json:"referenceCode"`
			StatusCode          string `json:"statusCode"`
			ServerReferenceCode string `json:"serverReferenceCode"`
			ErrorInfo           struct {
				Code        string `json:"errorCode"`
				Description string `json:"errorDescription"`
			} `json:"errorInfo"`
		} `json:"statusInfo"`
	} `json:"operator_response"`
}

// easyjson:json
type ConfirmResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"status_code"`
	OperatorResponse struct {
		StatusInfo struct {
			ReferenceCode       string `json:"referenceCode"`
			StatusCode          string `json:"statusCode"`
			ServerReferenceCode string `json:"serverReferenceCode"`
			TotalAmountCharged  string `json:"totalAmountCharged"`
		} `json:"statusInfo"`
	} `json:"operator_response"`
}

// easyjson:json
type SendSmsResponse struct {
	Status           string `json:"status"`
	StatusCode       string `json:"status_code"`
	OperatorResponse string `json:"operator_response"`
}
