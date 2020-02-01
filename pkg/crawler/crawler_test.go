package crawler_test

import (
	"fmt"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			testEnv.InitMockCrawler(0, 8079)
			time.Sleep(time.Second)
			crawlerX := crawler.Crawler{
				Url:         "http://localhost:8079",
				Name:        "Mock Crawler",
				Description: "This is a mock crawler",
				Source:      "S01",
			}
			_, err := crawlerX.SendRequest(nil, "Text")
			c.So(err, ShouldBeNil)
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
			testEnv.InitMultiCrawlers(5, time.Second*10, portStart)
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
