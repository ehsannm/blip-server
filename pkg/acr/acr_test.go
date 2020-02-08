package acr_test

import (
	testEnv "git.ronaksoftware.com/blip/server/pkg"
	"git.ronaksoftware.com/blip/server/pkg/acr"
	"github.com/kr/pretty"
	"testing"
)

/*
   Creation Time: 2019 - Oct - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	testEnv.Init()
}

func TestIdentifyByFile(t *testing.T) {
	music, err := acr.IdentifyByFile("./testdata/test2.m4a")
	if err != nil {
		t.Error(err)
	}
	pretty.Println(music)
}
