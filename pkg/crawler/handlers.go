package crawler

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris"
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

func Add(ctx iris.Context) {
	req := SaveReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	_, err = Save(&Crawler{
		httpClient:  http.Client{},
		Url:         req.Url,
		Name:        req.Name,
		Description: "",
		Source:      req.Source,
	})
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{
		Success: true,
	})

}
