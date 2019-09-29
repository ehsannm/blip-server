package saba

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Subscribe(phone string) (int, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	url := fmt.Sprintf("%s/v2/sub/%s/%s/%s",
		viper.GetString(config.ConfSmsServiceBaseUrl),
		viper.GetString(config.ConfSmsServiceName),
		viper.GetString(config.ConfSmsServiceToken),
		phone,
	)
	httpResp, err := c.Get(url)
	if err != nil {
		return 0, err
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(SubscribeResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return 0, err
	}

	return sResp.OtpID, nil
}

func Confirm(phone, phoneCode string, otpID int) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	httpResp, err := c.Get(fmt.Sprintf("%s/v2/confirm/%s/%s/%s/%s?otp_id=%d",
		viper.GetString(config.ConfSmsServiceBaseUrl),
		viper.GetString(config.ConfSmsServiceName),
		viper.GetString(config.ConfSmsServiceToken),
		phone, phoneCode, otpID,
	))
	if err != nil {
		return "", err
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(ConfirmResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return "", err
	}

	return sResp.StatusCode, nil
}
