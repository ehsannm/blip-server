package saba_test

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
   Creation Time: 2019 - Sep - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	// phone = "989122139561"
	phone = "989121228718"
	// phone = "989916695980"
	// phone = "989124861985"
)

func init() {
	config.Set(config.VasSabaServiceName, "musicchi")
	config.Set(config.VasSabaServiceToken, "65rejoptjb")
	config.Set(config.VasSabaServiceBaseUrl, "http://api.sabaeco.com")
	saba.Init()
}

func TestSubscribe(t *testing.T) {
	convey.Convey("Test SMS", t, func(c convey.C) {
		convey.Convey("Subscribe", func(c convey.C) {
			subscribeRes, err := saba.Subscribe(phone)
			if err != nil {
				_, _ = c.Println(err)
			}
			c.So(err, convey.ShouldBeNil)
			_, _ = c.Println(subscribeRes.OtpID)
		})
	})
}

func TestUnsubscribe(t *testing.T) {
	convey.Convey("Test SMS", t, func(c convey.C) {
		convey.Convey("Subscribe", func(c convey.C) {
			statusCode, err := saba.Unsubscribe(phone)
			if err != nil {
				_, _ = c.Println("Error", err)
			}
			c.So(err, convey.ShouldBeNil)
			_, _ = c.Println("statusCode", statusCode)
		})
	})
}

func TestConfirm(t *testing.T) {
	otpID := "15710572349449"
	convey.Convey("Test SMS", t, func(c convey.C) {
		convey.Convey("Confirm", func(c convey.C) {
			statusCode, err := saba.Confirm(phone, "4995", otpID)
			c.So(err, convey.ShouldBeNil)
			_, _ = c.Println(statusCode, otpID)
		})
	})
}
