package vas

import (
	"context"
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

var writeLogToDB = flusher.New(1000, 1, time.Millisecond*500, func(items []flusher.Entry) {
	bulkWrites := make([]mongo.WriteModel, 0, len(items))
	for idx := range items {
		bulkWrites = append(bulkWrites, mongo.NewInsertOneModel().SetDocument(items[idx].Value))
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	res, err := vasLogCol.BulkWrite(ctx, bulkWrites, options.BulkWrite().SetOrdered(false))
	if err != nil {
		log.Warn("Error On Writing MCI Notification to DB", zap.Error(err))
	} else {
		log.Debug("MCI Notifications was Written on DB", zap.Int64("Total", res.InsertedCount))
	}
})
