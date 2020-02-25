package help

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

/*
   Creation Time: 2020 - Feb - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// var defaultConfig = struct {
// 	ShowBlipLink         string `bson:"show_blip_link" json:"show_blip_link"`
// 	AndroidLatestVersion string `bson:"andorid_latest_version" json:"android_latest_version"`
// 	IosLatestVersion     string `bson:"ios_latest_version" json:"ios_latest_version"`
// }{}

var defaultConfig = map[string]string{}

func loadDefaultConfig() {
	res := helpCol.FindOne(nil, bson.M{"_id": "defaults"})
	err := res.Decode(&defaultConfig)
	if err != nil {
		log.Warn("Error On Reading Default Configs", zap.Error(err))
	}
}
