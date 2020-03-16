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
// SetDefaultConfig
type SetDefaultConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// easyjson:json
// UnsetDefaultConfig
type UnsetDefaultConfig struct {
	Key string `json:"key"`
}

const CConfig = "CONFIG"

// easyjson:json
// Config
type Config struct {
	UpdateAvailable bool   `json:"update_available"`
	UpdateForce     bool   `json:"update_force"`
	StoreLink       string `json:"store_link"`
	ShowBlipLink    bool   `json:"show_blip_link"`
	ShowShareLink   bool   `json:"show_share_link"`
	Authorized      bool   `json:"authorized"`
	VasEnabled      bool   `json:"vas_enabled"`
	MetrixToken     string `json:"metrix_token"`
}

// easyjson:json
// Feedback
type Feedback struct {
	Text string `json:"text"`
	Rate int    `json:"rate"`
}
