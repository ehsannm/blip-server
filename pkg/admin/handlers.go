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

func MigrateLegacyDBHandler(ctx iris.Context) {
	if !migrateRunning {
		migrateRunning = true
		go MigrateLegacyDB()
	}
}

func MigrateFilesHandler(ctx iris.Context) {
	if !migrateRunning {
		migrateRunning = true
		go MigrateFiles()
	}
}

func MigrateStatsHandler(ctx iris.Context) {
	msg.WriteResponse(ctx, CMigrateStats, MigrateStats{
		Scanned:           atomic.LoadInt32(&migrateScanned),
		Downloaded:        atomic.LoadInt32(&migrateDownloaded),
		AlreadyDownloaded: atomic.LoadInt32(&migrateAlreadyDownloaded),
		DownloadFailed:    atomic.LoadInt32(&migrateDownloadFailed),
	})
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
