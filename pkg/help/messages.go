package help

/*
   Creation Time: 2020 - Feb - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
// Config
type Config struct {
	UpdateAvailable bool   `json:"update_available"`
	UpdateForce     bool   `json:"update_force"`
	StoreLink       string `json:"store_link"`
	ShowBlipLink    bool   `json:"show_blip_link"`
	Authorized      bool   `json:"authorized"`
	VasEnabled      bool   `json:"vas_enabled"`
}
