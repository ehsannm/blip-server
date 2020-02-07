package admin

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
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

func MigrateLegacyDBStatsHandler(ctx iris.Context) {
	msg.WriteResponse(ctx, CMigrateStats, MigrateStats{
		Scanned:           atomic.LoadInt32(&migrateScanned),
		Downloaded:        atomic.LoadInt32(&migrateDownloaded),
		AlreadyDownloaded: atomic.LoadInt32(&migrateAlreadyDownloaded),
		DownloadFailed:    atomic.LoadInt32(&migrateDownloadFailed),
	})
}
