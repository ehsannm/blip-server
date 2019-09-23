package msg

import "github.com/kataras/iris"

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Item string

var (
	ErrAccessTokenInvalid Item = "ACCESS_KEY_INVALID"
	ErrAccessTokenExpired Item = "ACCESS_KEY_EXPIRED"
	ErrNoPermission       Item = "NO_PERMISSION"
	ErrWriteToDb          Item = "WRITE_TO_DB"
	ErrTokenInvalid       Item = "TOKEN_INVALID"
	ErrTokenExpired       Item = "TOKEN_EXPIRED"
	ErrPermissionIsNotSet Item = "PERMISSION_NOT_SET"
	ErrPhoneNotValid      Item = "PHONE_NOT_VALID"
	ErrPeriodNotValid     Item = "PERIOD_NOT_VALID"
)

func Error(ctx iris.Context, httpStatus int, errItem Item) {
	ctx.ContentType("application/json")
	resBytes, _ := CreateEnvelope("err", errItem).MarshalJSON()
	_, _ = ctx.Write(resBytes)
	ctx.StatusCode(httpStatus)
	ctx.StopExecution()
}
