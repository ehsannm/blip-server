package token

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const CValidated = "VALIDATED"

// easyjson:json
type Validated struct {
	RemainingDays int64 `json:"remaining_days"`
}
