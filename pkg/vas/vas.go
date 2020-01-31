package vas

import (
	"git.ronaksoftware.com/blip/server/internal/flusher"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	vasLogCol *mongo.Collection
)

var writeToDB = flusher.NewLifo(1000, 5, time.Millisecond*500, func(items []flusher.Entry) {
	docs := make([]interface{}, 0, len(items))
	for idx := range items {
		docs = append(docs, items[idx].Value)
	}
	res, err := vasLogCol.InsertMany(nil, docs, options.InsertMany().SetOrdered(false))
	if err != nil {
		log.Warn("Error On Writing MCI Notification to DB", zap.Error(err))
	} else {
		log.Debug("MCI Notifications was Written on DB", zap.Int("Total", len(res.InsertedIDs)))
	}
})
