package main

import (
	"encoding/json"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
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
	testEnv.Init()
	app := initServer()
	go app.Run(iris.Addr(":8080"), iris.WithOptimizations)

	time.Sleep(time.Second)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestAuth(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://localhost:8080",
		RequestFactory: httpexpect.DefaultRequestFactory{},
		Client:         http.DefaultClient,
		Reporter:       httpexpect.NewAssertReporter(t),
		Printers:       nil,
	})

	Convey("Auth Tests", t, func() {
		Convey("CreateAuth", func(c C) {
			reqBytes, _ := auth.CreateAccessToken{
				Permissions: []string{"read", "write"},
				Period:      90,
				AppName:     auth.AppNameMusicChi,
			}.MarshalJSON()
			r := e.POST("/auth/create").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithBytes(reqBytes).Expect().JSON().Object()
			c.So(r.Value("constructor").String().Raw(), ShouldEqual, auth.CAccessTokenCreated)
			_, _ = c.Println(r.Raw())

		})

		Convey("SendCode", func(c C) {
			reqBytes, _ := json.Marshal(auth.SendCodeReq{
				Phone: "989121228718",
			})
			r := e.POST("/auth/send_code").
				WithHeader(auth.HdrAccessKey, "ROOT").
				WithBytes(reqBytes).Expect().JSON().Object()
			c.So(r.Value("constructor").String().Raw(), ShouldEqual, auth.CPhoneCodeSent)

			r = r.Value("payload").Object()
			phoneCodeHash := r.Value("phone_code_hash").String().Raw()
			registered := r.Value("registered").Boolean().Raw()
			_, _ = c.Println(r.Raw(), registered, phoneCodeHash)
			if registered {
				Convey("Login", func(c C) {
					reqBytes, _ := json.Marshal(auth.LoginReq{
						PhoneCode:     "2374",
						PhoneCodeHash: phoneCodeHash,
						Phone:         "989121228718",
					})
					r := e.POST("/auth/login").
						WithHeader(auth.HdrAccessKey, "ROOT").
						WithBytes(reqBytes).Expect().JSON().Object()
					_, _ = c.Println(r.Raw())
					c.So(r.Value("constructor").String().Raw(), ShouldEqual, auth.CAuthorization)
				})
			} else {
				Convey("Register", func(c C) {
					reqBytes, _ := json.Marshal(auth.RegisterReq{
						PhoneCode:     "2374",
						PhoneCodeHash: phoneCodeHash,
						Phone:         "989121228718",
						Username:      "ehsan",
					})
					r := e.POST("/auth/register").
						WithHeader(auth.HdrAccessKey, "ROOT").
						WithBytes(reqBytes).Expect().JSON().Object()
					c.So(r.Value("constructor").String().Raw(), ShouldEqual, auth.CAuthorization)
				})
			}
		})
	})

}
