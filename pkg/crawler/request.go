package crawler

import (
	"bytes"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/mediocregopher/radix/v3"
	"net/http"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
type searchRequest struct {
	RequestID string `json:"req_id"`
	Keyword   string `json:"keyword"`
}

// easyjson:json
type searchResponse struct {
	RequestID string `json:"req_id"`
	Sources   string `json:"source"`
	Result    struct {
		SongUrl  string `json:"song_url"`
		CoverUrl string `json:"cover_url"`
		Lyrics   string `json:"lyrics,omitempty"`
		Artists  string `json:"artists"`
		Title    string `json:"title"`
		Genre    string `json:"genre"`
	} `json:"result"`
}

func getNextRequestID() string {
	return fmt.Sprintf("%s.%s", config.ServerID, ronak.RandomID(32))
}

func SendSearchRequest(crawlerUrl, keyword string) error {
	c := http.Client{
		Timeout: config.HttpRequestTimeout,
	}
	reqID := getNextRequestID()
	redisCache.Do(radix.FlatCmd(nil, "SET", fmt.Sprintf("CR-%s", reqID)))

	req := searchRequest{
		Keyword: keyword,
	}
	reqBytes, err := req.MarshalJSON()
	if err != nil {
		return err
	}
	httpRes, err := c.Post(crawlerUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusAccepted {
		return fmt.Errorf("invalid http response, got %d", httpRes.StatusCode)
	}
	return nil
}
