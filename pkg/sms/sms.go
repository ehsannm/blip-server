package sms

import (
	"time"
)

const (
	requestTimeout       = 10 * time.Second
	requestLongThreshold = 5 * time.Second
	queueSize            = 100
)

type Job struct {
	phoneNumber string
	body        string
}

// SMS Provider
type Provider interface {
	SendInBackground(phoneNumber string, txt string) (int, error)
	Send(phoneNumber string, txt string) (delivered bool, err error)
}
