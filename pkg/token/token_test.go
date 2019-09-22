package token_test

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/token"
	log "git.ronaksoftware.com/ronak/toolbox/logger"
	"github.com/iris-contrib/httpexpect"
	"github.com/kataras/iris"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
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

var (
	_Log log.Logger
)

func init () {
	_Log = log.NewConsoleLogger()
	if mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		_Log.Fatal("Error On MongoConnect", zap.Error(err))
	} else {
		token.InitMongo(mongoClient)
	}

	app := iris.New()
	app.UseGlobal(auth.GetAuthorization)
	tokenParty := app.Party("/token")
	tokenParty.Post("/create", auth.MustWriteAccess, token.CreateHandler)
	tokenParty.Get("/validate", auth.MustReadAccess, token.ValidateHandler)


	go func() {
		_ = app.Run(iris.Addr(":80"), iris.WithOptimizations)
	}()
	time.Sleep(time.Second)
}

func TestCreate(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://localhost:80",
		RequestFactory: httpexpect.DefaultRequestFactory{},
		Client:         http.DefaultClient,
		Reporter:       httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})
	r := e.POST("/token/create").
		WithHeader(auth.HdrAccessKey, "ROOT").
		WithQuery("Phone", "989121228718").
		WithQuery("Period", 90).
		Expect().JSON().Array()

	fmt.Println(r.Raw())
}


func TestValidate(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://localhost:80",
		RequestFactory: httpexpect.DefaultRequestFactory{},
		Client:         http.DefaultClient,
		Reporter:       httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})
	r := e.POST("/token/validate").
		WithHeader(auth.HdrAccessKey, "ROOT").
		WithQuery("DeviceID", "989121228718").
		WithQuery("Token", "90").
		Expect().JSON().Array()

	fmt.Println(r.Raw())
}