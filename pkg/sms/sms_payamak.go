package sms

import (
	"errors"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
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
   Creation Time: 2019 - Oct - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/


type Payamak struct {
	url string
	username string
	password string
	srcNumber string
	jobsChannel chan Job
}

func NewPayamak(un, pass, url, srcNumber string, jobConcurrency int) *Payamak {
	payamak := new(Payamak)
	payamak.username = un
	payamak.password = pass
	payamak.srcNumber = srcNumber
	payamak.url = url
	payamak.jobsChannel = make(chan Job, queueSize)
	for i := 0; i < jobConcurrency; i++ {
		go payamak.sendingJob()
	}
	return payamak
}

func (p *Payamak) SendInBackground(phoneNumber string, txt string) (int, error) {
	select {
	case p.jobsChannel <- Job{
		phoneNumber: phoneNumber,
		body:        txt,
	}:
		return len(p.jobsChannel), nil
	default:
		return len(p.jobsChannel), errors.New("sms queue is full")
	}
}

func (p *Payamak) Send(phoneNumber string, txt string) (delivered bool, err error) {
	v := url.Values{}
	v.Set("username", p.username)
	v.Set("password", p.password)
	v.Set("from", p.srcNumber)
	v.Set("to", phoneNumber)
	v.Set("text", txt)
	rb := strings.NewReader(v.Encode())
	c := http.DefaultClient
	c.Timeout = requestTimeout
	req, err := http.NewRequest("POST", p.url, rb)
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

func (p *Payamak) sendingJob() {
	for smsJob := range p.jobsChannel {
		startTime := time.Now()
		p.Send(smsJob.phoneNumber, smsJob.body)
		duration := time.Now().Sub(startTime)
		if duration > requestLongThreshold {
			log.Warn("SMS Request too Long",
				zap.Duration("Duration", duration),
			)
		}
	}
}
