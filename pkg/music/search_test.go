package music_test

import (
	"git.ronaksoftware.com/blip/server/internal/tools"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/music"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Feb - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func TestSearch(t *testing.T) {
	testEnv.InitMultiCrawlers(10, time.Second*10, 8181)
	time.Sleep(time.Second * 2)
	cursorID := tools.RandomID(12)
	keyword := "Song 1"
	Convey("Search By Text", t, func(c C) {
		Convey("Start Search (Consume All)", func(c C) {
			searchCtx := music.StartSearch(cursorID, keyword)
			for s := range searchCtx.SongChan() {
				c.So(s, ShouldNotBeNil)
			}
		})
		Convey("Resume Search (After Consume All)", func(c C) {
			songChan := music.ResumeSearch(cursorID)
			c.So(songChan, ShouldBeNil)
		})
		Convey("Start Search (Partial Consume)", func(c C) {
			searchCtx := music.StartSearch(cursorID, keyword)
			c.So(searchCtx, ShouldNotBeNil)
			songX := <-searchCtx.SongChan()
			c.So(songX, ShouldNotBeNil)
		})
		Convey("Resume Search (After Partial Consume)", func(c C) {
			searchCtx := music.ResumeSearch(cursorID)
			c.So(searchCtx, ShouldNotBeNil)
			for s := range searchCtx.SongChan() {
				c.So(s, ShouldNotBeNil)
			}
		})
	})
}
