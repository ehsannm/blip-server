package admin

import (
	"crypto/tls"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"time"
)

/*
   Creation Time: 2020 - Feb - 11
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate rm -f *_easyjson.go
//go:generate easyjson messages.go
var (
	httpClient http.Client
)

func init() {
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	if config.GetString(config.HttpProxy) != "" {
		if proxyURL, err := url.Parse(config.GetString(config.HttpProxy)); err == nil {
			httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
		} else {
			log.Warn("Error On Set HTTP Proxy", zap.Error(err))
		}
	}
}
