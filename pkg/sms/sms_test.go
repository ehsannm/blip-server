package sms

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"testing"
)

/*
   Creation Time: 2019 - Oct - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func TestPayamak(t *testing.T) {
	p := NewPayamak(
		config.GetString(config.SmsPayamakUser),
		config.GetString(config.SmsPayamakPass),
		config.GetString(config.SmsPayamakUrl),
		config.GetString(config.SmsPayamakPhone),
		10,
	)
	p.Send("989121228718", "Hi this is payamak test")
}
