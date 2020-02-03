package store_test

import (
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/store"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
   Creation Time: 2020 - Feb - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	testEnv.Init()
}

func TestStore(t *testing.T) {
	Convey("Test Store", t, func(c C) {
		err := store.DropAll()
		c.So(err, ShouldBeNil)
		err = store.Save(&store.Store{
			ID:       101,
			Dsn:      "mongodb://localhost:27001",
			Capacity: 0,
			Region:   "",
		})
		c.So(err, ShouldBeNil)
	})
}
