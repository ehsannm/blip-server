package saba_test

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/sms/saba"
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

func TestSubscribe(t *testing.T) {
	phone := "989122139561"
	otpID := 0
	config.Set(config.ConfSmsServiceName, "")
	config.Set(config.ConfSmsServiceToken, "")
	config.Set(config.ConfSmsServiceBaseUrl, "")
	convey.Convey("Test SMS", t, func(c convey.C) {
		convey.Convey("Subscribe", func(c convey.C) {
			var err error
			otpID, err = saba.Subscribe(phone)
			c.So(err, convey.ShouldBeNil)
			_, _ = c.Println(otpID)
		})
		convey.Convey("Confirm", func(c convey.C) {
			statusCode, err := saba.Confirm(phone, "0000", otpID)
			c.So(err, convey.ShouldBeNil)
			_, _ = c.Println(statusCode, otpID)
		})
	})
}
