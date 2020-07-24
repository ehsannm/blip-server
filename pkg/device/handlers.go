package device

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

/*
   Creation Time: 2020 - Mar - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// RegisterHandler is API handler
// API: /device/register
// Http Method: POST
// Inputs: JSON
//	token_type: string	(apn | fb)
//	token: string
// Returns: Bool (BOOL)
// Possible Errors:
//	1. 500: with error text
func RegisterDevice(ctx iris.Context) {
	req := &RegisterDeviceReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}
	s, _ := ctx.Values().Get(session.CtxSession).(session.Session)

	_, err = deviceCol.UpdateOne(nil,
		bson.M{"_id": s.ID},
		bson.M{"$set": bson.M{
			"token":      req.Token,
			"token_type": req.TokenType,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}
