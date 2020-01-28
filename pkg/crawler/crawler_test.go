package crawler_test

import (
	"fmt"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
			So(err, ShouldBeNil)
		})
		Convey("Set and Get", func(c C) {
			for i := 0; i < 10; i++ {
				crawlerID, err := crawler.Save(crawler.Crawler{
					Url:         fmt.Sprintf("http://crawler%d.com/some_text", i),
					Name:        fmt.Sprintf("Crawler %d", i),
					Description: fmt.Sprintf("Description for Crawler %d", i),
					Source:      fmt.Sprintf("Source %d", i%3),
				})
				So(err, ShouldBeNil)
				crawlerX, err := crawler.Get(crawlerID)
				So(err, ShouldBeNil)
				So(crawlerX, ShouldNotBeNil)
				So(crawlerX.Url, ShouldEqual, fmt.Sprintf("http://crawler%d.com/some_text", i))
				So(crawlerX.Source, ShouldEqual, fmt.Sprintf("Source %d", i%3))
			}
		})
		Convey("GetAll", func(c C) {
			crawlers, err := crawler.GetAll()
			So(err, ShouldBeNil)
			So(crawlers, ShouldNotBeNil)
			So(crawlers, ShouldHaveLength, 10)
		})
	})
}
