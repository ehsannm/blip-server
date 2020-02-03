package session_test

import (
	"git.ronaksoftware.com/blip/server/internal/tools"
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/session"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Feb - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	testEnv.Init()
}

func TestSession(t *testing.T) {
	Convey("Testing Session Functions", t, func(c C) {
		sessionID := tools.RandomID(12)
		userID := tools.RandomID(12)
		err := session.Save(&session.Session{
			ID:         sessionID,
			UserID:     userID,
			CreatedOn:  time.Now().Unix(),
			LastAccess: time.Now().Unix(),
		})
		c.So(err, ShouldBeNil)
		err = session.Remove(userID, "MusicChi")
		c.So(err, ShouldBeNil)
		sessions, err := session.GetAll(userID)
		c.So(err, ShouldBeNil)
		c.So(sessions, ShouldHaveLength, 0)
		err = session.Save(&session.Session{
			ID:         sessionID,
			UserID:     userID,
			CreatedOn:  0,
			LastAccess: 0,
			App:        "MusicChi",
		})
		c.So(err, ShouldBeNil)
		err = session.Remove(userID, "MusicChi")
		c.So(err, ShouldBeNil)
		sessions, err = session.GetAll(userID)
		c.So(err, ShouldBeNil)
		c.So(sessions, ShouldHaveLength, 0)
	})
}
