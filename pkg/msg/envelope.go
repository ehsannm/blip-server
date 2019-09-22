package msg

import (
	"github.com/kataras/iris"
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
	resBytes, _ := ResponseEnvelope{
		Constructor: constructor,
		Payload:     payload,
	}.MarshalJSON()
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.Write(resBytes)
	ctx.StopExecution()
}
