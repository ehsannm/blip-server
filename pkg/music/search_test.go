package music_test

import (
	"fmt"
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

func TestIndex(t *testing.T) {
	Convey("Song Search Functionality", t, func() {
		Convey("Index", func(c C) {
			for i := 0; i < 10; i++ {
				s := &music.Song{
					ID:             primitive.NewObjectID(),
					Title:          fmt.Sprintf("song %d", i),
					Genre:          fmt.Sprintf("genre %d", i),
					Lyrics:         fmt.Sprintf("lyrics %d", i),
					Artists:        fmt.Sprintf("artists %d", i),
					CoverUrl:       fmt.Sprintf("http://cdn.blip.fun/cover_%d", i),
					SongUrl:        fmt.Sprintf("http://cdn.blip.fun/song_%d", i),
					OriginCoverUrl: fmt.Sprintf("http://cdn.blip.fun/org_cover_%d", i),
					OriginSongUrl:  fmt.Sprintf("http://cdn.blip.fun/org_song_%d", i),
					Source:         fmt.Sprintf("Source %d", i%3),
				}
				err := music.IndexSong(s)
				So(err, ShouldBeNil)
			}

		})
		Convey("Search", func(c C) {
			songIDs, err := music.SearchIndex("song")
			c.So(err, ShouldBeNil)
			for _, songID := range songIDs {
				_, _ = c.Println(songID.Hex())
			}
		})
	})

}
