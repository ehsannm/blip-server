package sms

import (
	"errors"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Apr - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// ADP is a SMS Provider
type ADP struct {
	username    string
	password    string
	url         string
	phone       string
	jobsChannel chan Job
}

func NewADP(un, pass, url, phone string, jobConcurrency int) *ADP {
	adp := new(ADP)
	adp.username = un
	adp.password = pass
	adp.url = url
	adp.phone = phone
	adp.jobsChannel = make(chan Job, queueSize)
	for i := 0; i < jobConcurrency; i++ {
		go adp.sendingJob()
	}
	return adp
}

func (adp *ADP) SendInBackground(phoneNumber string, txt string) (int, error) {
	select {
	case adp.jobsChannel <- Job{
		phoneNumber: phoneNumber,
		body:        txt,
	}:
		return len(adp.jobsChannel), nil
	default:
		return len(adp.jobsChannel), errors.New("sms queue is full")
	}
}

func (adp *ADP) Send(phoneNumber string, txt string) (delivered bool, err error) {
	v := url.Values{}
	v.Set("username", adp.username)
	v.Set("password", adp.password)
	v.Set("dstaddress", phoneNumber)
	v.Set("srcaddress", adp.phone)
	v.Set("body", txt)
	v.Set("unicode", "1")
	rb := strings.NewReader(v.Encode())
	c := http.DefaultClient
	c.Timeout = requestTimeout
	req, err := http.NewRequest("POST", adp.url, rb)
	if err != nil {
		log.Warn("Error ADP Send",
			zap.Error(err),
			zap.String("Phone", phoneNumber),
		)
		return false, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.Do(req)
	if err != nil {
		log.Warn("Error ADP Response",
			zap.Error(err),
			zap.String("Phone", phoneNumber),
		)
		return false, err
	}
	if ce := log.Check(zapcore.DebugLevel, "SMS Response"); ce != nil {
		rb, _ := ioutil.ReadAll(res.Body)
		ce.Write(
			zap.String("Res", string(rb)),
		)
	} else {
		io.Copy(ioutil.Discard, res.Body)
	}
	res.Body.Close()

	return true, nil
}

func (adp *ADP) sendingJob() {
	for smsJob := range adp.jobsChannel {
		startTime := time.Now()
		adp.Send(smsJob.phoneNumber, smsJob.body)
		duration := time.Now().Sub(startTime)
		if duration > requestLongThreshold {
			log.Warn("SMS Request too Long",
				zap.Duration("Duration", duration),
			)
		}
	}
}
