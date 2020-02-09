package crawler

import (
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func SaveHandler(ctx iris.Context) {
	req := &SaveReq{}
	err := ctx.ReadJSON(req)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.ErrCannotUnmarshalRequest)
		return
	}

	crawlerID, err := Save(&Crawler{
		httpClient:     http.Client{},
		ID:             req.ID,
		Url:            req.Url,
		Name:           req.Name,
		Description:    "",
		Source:         req.Source,
		DownloaderJobs: req.DownloaderJobs,
	})
	if err != nil {
		msg.WriteError(ctx, http.StatusInternalServerError, msg.ErrWriteToDb)
		return
	}

	msg.WriteResponse(ctx, CCrawlersCreated, CrawlersCreated{
		CrawlerID: crawlerID,
	})
}

func ListHandler(ctx iris.Context) {
	crawlers := GetAll()
	msg.WriteResponse(ctx, CCrawlersMany, CrawlersMany{
		Crawlers: crawlers,
	})
}

func RemoveHandler(ctx iris.Context) {
	crawlerID, err := primitive.ObjectIDFromHex(ctx.Params().GetString("crawlerID"))
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.Item(err.Error()))
		return
	}

	err = Remove(crawlerID)
	if err != nil {
		msg.WriteError(ctx, http.StatusBadRequest, msg.Item(err.Error()))
		return
	}

	msg.WriteResponse(ctx, msg.CBool, msg.Bool{Success: true})
}
