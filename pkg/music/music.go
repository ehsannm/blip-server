package music

import (
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"github.com/gobwas/pool/pbytes"
	"net/http"
	"net/http/httputil"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var reverseProxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "ws.blipapi.xyz"
		req.URL.Path = "blip-v2/music_chi/voice/upload"
		req.URL.RawQuery = ""
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "Blip Sever v2")
		}
		req.Header.Del(auth.HdrAccessKey)
		req.Header.Del(session.HdrSessionID)
	},
	Transport:      nil,
	FlushInterval:  0,
	ErrorLog:       nil,
	BufferPool:     buffPool{},
	ModifyResponse: nil,
	ErrorHandler:   nil,
}

const (
	bufferPoolSize = 1 << 10 * 128
)

type buffPool struct{}

func (b buffPool) Get() []byte {
	return pbytes.GetLen(bufferPoolSize)
}

func (b buffPool) Put(slice []byte) {
	pbytes.Put(slice)
}
