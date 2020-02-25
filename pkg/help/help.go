package help

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"sync"
)

/*
   Creation Time: 2020 - Feb - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var configMtx sync.RWMutex
var defaultConfig = map[string]string{}

func loadDefaultConfig() {
	configMtx.Lock()
	res := helpCol.FindOne(nil, bson.M{"_id": "defaults"})
	err := res.Decode(&defaultConfig)
	if err != nil {
		log.Warn("Error On Reading Default Configs", zap.Error(err))
	}
	configMtx.Unlock()
}

func getConfig(key string) string {
	configMtx.RLock()
	defer configMtx.RUnlock()
	return defaultConfig[key]
}
