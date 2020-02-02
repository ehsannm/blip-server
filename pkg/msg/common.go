package msg

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const CBool = "BOOL"

// easyjson:json
// Bool
type Bool struct {
	Success bool `json:"success"`
}
