package saba

import (
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/config"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	baseUrl = strings.TrimRight(config.GetString(config.VasSabaServiceBaseUrl), "/")
	serviceName = config.GetString(config.VasSabaServiceName)
	serviceToken = config.GetString(config.VasSabaServiceToken)
}

func Subscribe(phone string) (*SubscribeResponse, error) {
	c := http.Client{
		Timeout: time.Second * 30,
	}

	httpResp, err := c.Get(fmt.Sprintf("%s/v2/sub/%s/%s/%s",
		baseUrl, serviceName, serviceToken, phone,
	))
	if err != nil {
		return nil, errors.Wrap(err, "Error In Response")
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(SubscribeResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return nil, err
	}

	if ce := log.Check(log.DebugLevel, "Saba SMS Subscribe Response"); ce != nil {
		ce.Write(
			zap.String("Code", sResp.StatusCode),
			zap.String("Status", sResp.Status),
			zap.String("Res", tools.ByteToStr(httpBytes)),
		)
	}

	return sResp, nil
}

func Unsubscribe(phone string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 30,
	}
	httpResp, err := c.Get(fmt.Sprintf("%s/v2/unsub/%s/%s/%s",
		baseUrl, serviceName, serviceToken, phone,
	))
	if err != nil {
		return "", errors.Wrap(err, "Error In Response")
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(UnsubscribeResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return "", errors.Wrap(err, "Error In Unmarshal Response")
	}
	if ce := log.Check(log.DebugLevel, "Saba SMS Unsubscribe Response"); ce != nil {
		ce.Write(
			zap.String("Code", sResp.StatusCode),
			zap.String("Status", sResp.Status),
			zap.String("Res", tools.ByteToStr(httpBytes)),
		)
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
		return "", errors.Wrap(err, "Error In Response")
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(ConfirmResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return "", errors.Wrap(err, "Error In Unmarshal Response")
	}
	if ce := log.Check(log.DebugLevel, "Saba SMS Confirm Response"); ce != nil {
		ce.Write(
			zap.String("Code", sResp.StatusCode),
			zap.String("Status", sResp.Status),
			zap.String("Res", tools.ByteToStr(httpBytes)),
		)
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
		viper.GetString(config.VasSabaServiceBaseUrl),
		viper.GetString(config.VasSabaServiceName),
		viper.GetString(config.VasSabaServiceToken),
		phone, v.Encode(),
	))
	if err != nil {
		return nil, errors.Wrap(err, "Error In Response")
	}

	httpBytes, _ := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	sResp := new(SendSmsResponse)
	err = sResp.UnmarshalJSON(httpBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Error In Unmarshal Response")
	}

	return sResp, nil

}
