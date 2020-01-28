package crawler_test

import (
	"fmt"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/tcplisten"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	testEnv.Init()
}

type mockCrawler struct{}

func (m mockCrawler) ServeHTTP(httpRes http.ResponseWriter, httpReq *http.Request) {
	reqData, _ := ioutil.ReadAll(httpReq.Body)
	req := crawler.SearchRequest{}
	_ = req.UnmarshalJSON(reqData)

	res := crawler.SearchResponse{
		RequestID: req.RequestID,
		Sources:   "Source 01",
		Result: struct {
			SongUrl  string `json:"song_url"`
			CoverUrl string `json:"cover_url"`
			Lyrics   string `json:"lyrics,omitempty"`
			Artists  string `json:"artists"`
			Title    string `json:"title"`
			Genre    string `json:"genre"`
		}{
			SongUrl:  "http://url.com",
			CoverUrl: "http://cover-url.com",
			Lyrics:   "This is some lyrics text",
			Artists:  "Some Famous Artist",
			Title:    req.Keyword,
			Genre:    "Rock",
		},
	}

	resData, _ := res.MarshalJSON()
	httpRes.Write(resData)
}

func initMockCrawler(port int) {
	s := httptest.NewUnstartedServer(mockCrawler{})
	tcpConfig := tcplisten.Config{}
	s.Listener, _ = tcpConfig.NewListener("tcp4", ":8080")
	s.Start()

}

func TestCrawler(t *testing.T) {
	Convey("Crawler Functionality", t, func() {
		Convey("DropAll", func(c C) {
			err := crawler.DropAll()
			c.So(err, ShouldBeNil)
		})
		Convey("Set and Get", func(c C) {
			for i := 0; i < 10; i++ {
				crawlerID, err := crawler.Save(crawler.Crawler{
					Url:         fmt.Sprintf("http://crawler%d.com/some_text", i),
					Name:        fmt.Sprintf("Crawler %d", i),
					Description: fmt.Sprintf("Description for Crawler %d", i),
					Source:      fmt.Sprintf("Source %d", i%3),
				})
				c.So(err, ShouldBeNil)
				crawlerX, err := crawler.Get(crawlerID)
				c.So(err, ShouldBeNil)
				c.So(crawlerX, ShouldNotBeNil)
				c.So(crawlerX.Url, ShouldEqual, fmt.Sprintf("http://crawler%d.com/some_text", i))
				c.So(crawlerX.Source, ShouldEqual, fmt.Sprintf("Source %d", i%3))
			}
		})
		Convey("GetAll", func(c C) {
			time.Sleep(time.Second)
			crawlers := crawler.GetAll()
			c.So(crawlers, ShouldNotBeNil)
			c.So(crawlers, ShouldHaveLength, 10)
		})
		Convey("Send Request", func(c C) {
			initMockCrawler(8080)
			time.Sleep(time.Second)
			crawlerX := crawler.Crawler{
				Url:         "http://localhost:8080",
				Name:        "Mock Crawler",
				Description: "This is a mock crawler",
				Source:      "S01",
			}
			res, err := crawlerX.SendRequest("", "Text")
			c.So(err, ShouldBeNil)
			c.So(res.Result.Title, ShouldEqual, "Text")
		})
	})
}
