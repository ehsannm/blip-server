package msg

import (
	"github.com/kataras/iris/v12"
	"net/http"
)

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate rm -f *_easyjson.go
//go:generate easyjson envelope.go

// easyjson:json
type ResponseEnvelope struct {
	Constructor string      `json:"constructor"`
	Payload     interface{} `json:"payload"`
}

func CreateEnvelope(constructor string, payload interface{}) *ResponseEnvelope {
	return &ResponseEnvelope{
		Constructor: constructor,
		Payload:     payload,
	}
}

func WriteResponse(ctx iris.Context, constructor string, payload interface{}) {
	ctx.ContentType("application/json")
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ResponseEnvelope{
		Constructor: constructor,
		Payload:     payload,
	})
	ctx.StopExecution()
}

func WriteError(ctx iris.Context, httpStatus int, errItem Item) {
	ctx.StatusCode(httpStatus)
	ctx.ContentType("application/json")
	_, _ = ctx.JSON(CreateEnvelope("err", errItem))
	ctx.StopExecution()
}
