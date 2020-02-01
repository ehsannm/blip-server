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
		Convey("Prepare Database", func(c C) {
			// for i := 0; i < 10; i++ {
			// 	title := fmt.Sprintf("Title %d", i)
			// 	artists := fmt.Sprintf("Artist %d", i)
			// 	_, err := music.SaveSong(&music.Song{
			// 		UniqueKey:      music.GenerateUniqueKey(title, artists),
			// 		Title:          title,
			// 		Genre:          "Genre",
			// 		Lyrics:         fmt.Sprintf("Lyrics %d", i),
			// 		Artists:        artists,
			// 		Cdn:            "",
			// 		OriginCoverUrl: fmt.Sprintf("http://url.com/file/%d", i),
			// 		OriginSongUrl:  fmt.Sprintf("http://url.com/file/%d", i),
			// 		Source:         fmt.Sprintf("Source %d", i),
			// 	})
			// 	c.So(err, ShouldBeNil)
			// }
			err := music.DropAllSongs()
			c.So(err, ShouldBeNil)
		})
		Convey("Start Search (Consume All)", func(c C) {
			songChan := music.StartSearch(cursorID, keyword)
			for s := range songChan {
				_, _ = c.Println(s.ID, s.Title, s.Artists)
			}
		})
		Convey("Resume Search (After Consume All)", func(c C) {
			songChan := music.ResumeSearch(cursorID)
			c.So(songChan, ShouldBeNil)
		})
		Convey("Start Search (Partial Consume)", func(c C) {
			songChan := music.StartSearch(cursorID, keyword)
			c.So(songChan, ShouldNotBeNil)
			songX := <-songChan
			c.So(songX, ShouldNotBeNil)
			c.Println(songX.Source, songX.Title, songX.Artists)

		})
		Convey("Resume Search (After Partial Consume)", func(c C) {
			songChan := music.ResumeSearch(cursorID)
			c.So(songChan, ShouldNotBeNil)
			for s := range songChan {
				_, _ = c.Println(s.ID, s.Title, s.Artists)
			}
		})
	})
}
