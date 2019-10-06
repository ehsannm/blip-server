package main

import (
	"encoding/json"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"os"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - Sep - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func TestMain(m *testing.M) {
	app := initServer()
	go app.Run(iris.Addr(":80"), iris.WithOptimizations)

	time.Sleep(time.Second)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestAuth(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://localhost:80",
		RequestFactory: httpexpect.DefaultRequestFactory{},
		Client:         http.DefaultClient,
		Reporter:       httpexpect.NewAssertReporter(t),
		Printers: nil ,
	})

	Convey("Auth Tests", t, func() {
		Convey("CreateAuth", func(c C) {
			r := e.POST("/auth/create").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithFormField("Permissions", "read").
				WithFormField("Permissions", "write").
				WithFormField("Period", 90).
				Expect().JSON().Object()

			c.Println(r.Raw())
			// r.Value("constructor").Equal(auth.CAccessTokenCreated)

		})

		Convey("SendCode", func(c C) {
			reqBytes, _ := json.Marshal(auth.SendCodeReq{
				Phone: "989121228718",
			})
			r := e.POST("/auth/send_code").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithBytes(reqBytes).Expect().JSON()
			c.Println(r.Raw())
		})

		Convey("Register", func(c C) {
			reqBytes, _ := json.Marshal(auth.RegisterReq{
				PhoneCode:     "0",
				PhoneCodeHash: "",
				Phone:         "989121228718",
				OperationID:   0,
				Username:      "ehsan",
			})
			r := e.POST("/auth/register").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithBytes(reqBytes).Expect().JSON()
			c.Println(r.Raw())
		})
	})

}
