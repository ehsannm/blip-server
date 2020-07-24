package token_test

import (
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/token"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - Sep - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	BaseUrl = "http://localhost:8989"
)

func init() {
	testEnv.Init()

	app := iris.New()
	app.UseGlobal(auth.MustHaveAccessKey)
	tokenParty := app.Party("/token")
	tokenParty.Post("/create", auth.MustWriteAccess, token.CreateHandler)
	tokenParty.Post("/validate", auth.MustReadAccess, token.ValidateHandler)

	go func() {
		_ = app.Run(iris.Addr(":8989"), iris.WithOptimizations)
	}()
	time.Sleep(time.Second)
}

func TestToken(t *testing.T) {
	genToken := ""
	Convey("TestToken", t, func(c C) {
		e := httpexpect.New(t, BaseUrl)
		Convey("Create Token", func(c C) {
			r := e.POST("/token/create").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithFormField("Phone", "989121228718").
				WithFormField("Period", 90).
				Expect().JSON().Array()
			genToken = r.Element(0).Array().Element(0).String().Raw()
		})

		Convey("Validate Token", func(c C) {
			r := e.POST("/token/validate").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithFormField("DeviceID", "989121228718").
				WithFormField("Token", genToken).
				Expect().JSON().Object()
			r.Value("constructor").String().Equal(token.CValidated)
			_, _ = c.Println(r.Value("payload").Object().Value("remaining_days").Number().Raw())
		})

	})

}
