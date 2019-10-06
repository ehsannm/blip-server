package saba

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

var (
	baseUrl      string
	serviceName  string
	serviceToken string
)

func Init() {
	baseUrl = strings.TrimRight(config.GetString(config.ConfSmsServiceBaseUrl), "/")
	serviceName = config.GetString(config.ConfSmsServiceName)
	serviceToken = config.GetString(config.ConfSmsServiceToken)
}

func Subscribe(phone string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	url := fmt.Sprintf("%s/v2/sub/%s/%s/%s",
		baseUrl, serviceName, serviceToken, phone,
	)
	httpResp, err := c.Get(url)
	if err != nil {
		return "", err
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(SubscribeResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return "", err
	}

	return sResp.OtpID, nil
}

func Unsubscribe(phone string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	url := fmt.Sprintf("%s/v2/unsub/%s/%s/%s",
		baseUrl, serviceName, serviceToken, phone,
	)
	httpResp, err := c.Get(url)
	if err != nil {
		return "", err
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(UnsubscribeResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return "", err
	}

	return sResp.StatusCode, nil
}

func Confirm(phone, phoneCode string, otpID string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	httpResp, err := c.Get(fmt.Sprintf("%s/v2/confirm/%s/%s/%s/%s?otp_id=%s",
		baseUrl, serviceName, serviceToken, phone, phoneCode, otpID,
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

func SendMessage(phone, message string) (*SendSmsResponse, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	v := url.Values{}
	v.Set("message", message)


	httpResp, err := c.Get(fmt.Sprintf("%s/v2/send/%s/%s/%s?%s",
		viper.GetString(config.ConfSmsServiceBaseUrl),
		viper.GetString(config.ConfSmsServiceName),
		viper.GetString(config.ConfSmsServiceToken),
		phone, v.Encode(),
	))
	if err != nil {
		return nil, err
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(SendSmsResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return nil, err
	}

	return sResp, nil

}
