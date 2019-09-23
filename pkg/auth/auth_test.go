package auth_test

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
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

func init() {
	_Log = log.NewConsoleLogger()
	if mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		_Log.Fatal("Error On MongoConnect", zap.Error(err))
	} else {
		auth.InitMongo(mongoClient)
	}

	app := iris.New()
	app.UseGlobal(auth.GetAuthorizationHandler)
	authParty := app.Party("/auth")
	authParty.Post("/create", auth.MustAdminHandler, auth.CreateAccessKeyHandler)

	go func() {
		_ = app.Run(iris.Addr(":80"), iris.WithOptimizations)
	}()
	time.Sleep(time.Second)
}

func TestCreateAccessKey(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:        "http://localhost:80",
		RequestFactory: httpexpect.DefaultRequestFactory{},
		Client:         http.DefaultClient,
		Reporter:       httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})
	r := e.POST("/auth/create").
		WithHeader(auth.HdrAccessKey, "ROOT").
		WithQuery("Permissions", []string{"read", "write"}).
		WithQuery("Period", 90).
		Expect().JSON().Object()

	fmt.Println(r.Raw())
	r.Value("constructor").Equal(auth.CAccessTokenCreated)
	payload := r.Value("payload").Object()
	payload.Value("access_token").String().Length().Equal(64)

}
