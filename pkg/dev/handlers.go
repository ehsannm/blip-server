package dev

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/sms/saba"
	"github.com/kataras/iris"
	"net/http"
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
	phone := ctx.PostValue("phone")
	res, err := saba.Unsubscribe(phone)
	if err != nil {
		msg.Error(ctx, http.StatusInternalServerError, msg.Item(err.Error()))
		return
	}
	ctx.WriteString(res)
}