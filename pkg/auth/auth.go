package auth

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
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Permission byte

const (
	_ Permission = 1 << iota
	Admin
	Read
	Write
)

// Auth
type Auth struct {
	ID          string       `bson:"_id"`
	Permissions []Permission `bson:"perm"`
	CreatedOn   int64        `bson:"created_on"`
	ExpiredOn   int64        `bson:"expired_on"`
	AppName     string       `bson:"app_name"`
}

var writeLogToDB = flusher.New(1000, 1, time.Millisecond*500, func(items []flusher.Entry) {
	bulkWrites := make([]mongo.WriteModel, 0, len(items))
	for idx := range items {
		bulkWrites = append(bulkWrites, mongo.NewInsertOneModel().SetDocument(items[idx].Value))
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	res, err := authLogCol.BulkWrite(ctx, bulkWrites, options.BulkWrite().SetOrdered(false))
	if err != nil {
		log.Warn("Error On Writing AuthLog to DB", zap.Error(err))
	} else {
		log.Debug("AuthLog was Written on DB", zap.Int64("Total", res.InsertedCount))
	}
})
