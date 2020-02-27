package admin

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	"github.com/kataras/iris"
	"net/http"
	"sync/atomic"
)

/*
   Creation Time: 2019 - Oct - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Unsubscribe(ctx iris.Context) {
	phone := ctx.FormValue("phone")
	if len(phone) < 5 {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrPhoneNotValid)
		return
	}

	res, err := saba.Unsubscribe(phone)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}

	msg.WriteResponse(ctx, CUnsubscribed, Unsubscribed{
		Phone:      phone,
		StatusCode: res,
	})

}

func HealthCheckDbHandler(ctx iris.Context) {
	if !healthCheckRunning {
		healthCheckRunning = true
		go HealthCheckDB()
	}
}

func HealthCheckStatsHandler(ctx iris.Context) {
	msg.WriteResponse(ctx, CHealthCheckStats, HealthCheckStats{
		Scanned:    atomic.LoadInt32(&scanned),
		CoverFixed: atomic.LoadInt32(&coverFixed),
		SongFixed:  atomic.LoadInt32(&songFixed),
	})
}

func HealthCheckStoreHandler(ctx iris.Context) {
	if !healthCheckRunning {
		healthCheckRunning = true
		go HealthCheckStore()
	}
}

// SetVas is API handler
// API: /admin/vas
// Http Method: POST
// Inputs: JSON
//	user_id: string
//	enabled: bool
// Returns: Bool (BOOL)
// Possible Errors:
//	1. 403: USER_NOT_FOUND
//	2. 500: WRITE_TO_DB
func SetVas(ctx iris.Context) {
	req := &SetVasReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	actorUser, err := user.Get(req.UserID)
	if err != nil {
		msg.WriteError(ctx, http.StatusNotFound, msg.ErrUserNotFound)
		return
	}

	actorUser.VasPaid = req.Enabled
	err = user.Save(actorUser)
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}
