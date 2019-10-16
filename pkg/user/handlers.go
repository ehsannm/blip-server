package user

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/kataras/iris"
	"net/http"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func MustVasEnabled(ctx iris.Context) {
	s, ok := ctx.Values().Get(session.CtxSession).(session.Session)
	if !ok {
		msg.Error(ctx, http.StatusForbidden, msg.ErrSessionInvalid)
		return
	}
	u, err := Get(s.UserID)
	if err != nil {
		msg.Error(ctx, http.StatusForbidden, msg.ErrSessionInvalid)
		return
	}

	if !u.VasPaid {
		msg.Error(ctx, http.StatusForbidden, msg.ErrVasIsNotEnabled)
		return
	}
}
