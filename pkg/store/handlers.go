package store

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Save(ctx iris.Context) {
	req := &SaveStoreReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	storeX := &Store{
		ID:       req.StoreID,
		Dsn:      req.Dsn,
		Capacity: req.Capacity,
		Region:   req.Region,
	}
	_, err = storeCol.UpdateOne(nil, bson.M{"_id": storeX.ID}, bson.M{"$set": storeX}, options.Update().SetUpsert(true))
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, &msg.Bool{
		Success: true,
	})
}

func Get(ctx iris.Context) {
	req := &GetStoreReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	stores := &Stores{}
	for _, storeID := range req.StoreIDs {
		if storeX := get(storeID); storeX != nil {
			stores.Stores = append(stores.Stores, storeX)
		}
	}

	msg.WriteResponse(ctx, CStores, stores)
}
func get(storeID int64) *Store {
	storesMtx.RLock()
	s := stores[storeID]
	storesMtx.RUnlock()
	return s
}
