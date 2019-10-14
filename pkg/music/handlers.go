package music

import (
	"encoding/base64"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"github.com/kataras/iris"
	"net/http"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func SearchBySound(ctx iris.Context) {
	sound := ctx.PostValue("sound")
	soundBytes, err := base64.StdEncoding.DecodeString(sound)
	if err != nil {
		msg.Error(ctx, http.StatusBadRequest, msg.ErrBadSoundFile)
		return
	}

	music, err := acr.IdentifyByByteString(soundBytes)
	if err != nil {
		msg.Error(ctx, http.StatusNotAcceptable, msg.Item(err.Error()))
		return
	}

	// music.Status.Code

}
