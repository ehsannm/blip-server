package music_test

import (
	"fmt"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/music"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestSong(t *testing.T) {
	Convey("Song Functionality", t, func() {
		Convey("DropAll", func(c C) {
			err := music.DropAllSongs()
			So(err, ShouldBeNil)
		})
		Convey("Set and Get", func(c C) {
			for i := 0; i < 10; i++ {
				songID, err := music.SaveSong(&music.Song{
					ID:             primitive.NilObjectID,
					Title:          fmt.Sprintf("Song %d", i),
					Genre:          fmt.Sprintf("Genre %d", i),
					Lyrics:         fmt.Sprintf("Lyrics %d", i),
					Artists:        fmt.Sprintf("Artists %d", i),
					CoverUrl:       fmt.Sprintf("http://cdn.blip.fun/cover_%d", i),
					SongUrl:        fmt.Sprintf("http://cdn.blip.fun/song_%d", i),
					OriginCoverUrl: fmt.Sprintf("http://cdn.blip.fun/org_cover_%d", i),
					OriginSongUrl:  fmt.Sprintf("http://cdn.blip.fun/org_song_%d", i),
					Source:         fmt.Sprintf("Source %d", i%3),
				})
				So(err, ShouldBeNil)
				songX, err := music.GetSong(songID)
				So(err, ShouldBeNil)
				So(songX, ShouldNotBeNil)
				So(songX.Title, ShouldEqual, fmt.Sprintf("Song %d", i))
				So(songX.Source, ShouldEqual, fmt.Sprintf("Source %d", i%3))
			}
		})
		Convey("Search", func(c C) {
			songIDs, err := music.SearchLocalIndex("song")
			c.So(err, ShouldBeNil)
			c.So(len(songIDs), ShouldBeGreaterThan, 0)
		})
	})
}
