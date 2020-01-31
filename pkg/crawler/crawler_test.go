package crawler_test

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/tcplisten"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type mockCrawler struct {
	MaxDelay time.Duration
}

func (m mockCrawler) ServeHTTP(httpRes http.ResponseWriter, httpReq *http.Request) {
	if m.MaxDelay > 0 {
		time.Sleep(time.Duration(tools.RandomInt(int(m.MaxDelay))))
	}
	reqData, _ := ioutil.ReadAll(httpReq.Body)
	req := crawler.SearchRequest{}
	_ = req.UnmarshalJSON(reqData)
	res := crawler.SearchResponse{
		RequestID: req.RequestID,
		Source:    "Source 01",
		Result: struct {
			SongUrl  string `json:"song_url"`
			CoverUrl string `json:"cover_url"`
			Lyrics   string `json:"lyrics,omitempty"`
			Artists  string `json:"artists"`
			Title    string `json:"title"`
			Genre    string `json:"genre,omitempty"`
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

func initMockCrawler(maxDelay time.Duration, port int) {
	s := httptest.NewUnstartedServer(mockCrawler{
		MaxDelay: maxDelay,
	})
	tcpConfig := tcplisten.Config{}
	s.Listener, _ = tcpConfig.NewListener("tcp4", fmt.Sprintf(":%d", port))
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
				crawlerID, err := crawler.Save(&crawler.Crawler{
					ID:          primitive.NewObjectID(),
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
			initMockCrawler(0, 8079)
			time.Sleep(time.Second)
			crawlerX := crawler.Crawler{
				Url:         "http://localhost:8079",
				Name:        "Mock Crawler",
				Description: "This is a mock crawler",
				Source:      "S01",
			}
			res, err := crawlerX.SendRequest(nil, "Text")
			c.So(err, ShouldBeNil)
			c.So(res.Result.Title, ShouldEqual, "Text")
		})
	})
}

func TestSearch(t *testing.T) {
	Convey("Crawler Search Functionality", t, func() {
		Convey("DropAll", func(c C) {
			err := crawler.DropAll()
			c.So(err, ShouldBeNil)
		})
		Convey("Run Mock Servers", func(c C) {
			portStart := 8081
			for i := 0; i < 5; i++ {
				initMockCrawler(time.Second*10, portStart+i)
				crawlerX := &crawler.Crawler{
					ID:          primitive.NewObjectID(),
					Url:         fmt.Sprintf("http://localhost:%d", portStart+i),
					Name:        fmt.Sprintf("Crawler %d", i),
					Description: "This is a Mock Crawler",
					Source:      fmt.Sprintf("Source %d", i%3),
				}
				_, err := crawler.Save(crawlerX)
				c.So(err, ShouldBeNil)
			}
		})
		Convey("Wait For Crawlers to Run", func(c C) {
			time.Sleep(time.Second)
		})
		Convey("Send Search Request", func(c C) {
			keyword := "Some Text"

			resChan := crawler.Search(nil, keyword)
			for x := range resChan {
				_, _ = c.Println("Response:", x.RequestID)
			}
		})
	})
}
