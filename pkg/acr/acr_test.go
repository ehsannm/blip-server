package acr_test

import (
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

func TestIdentifyByFile(t *testing.T) {
	acr.Init()
	music, err := acr.IdentifyByFile("./testdata/test2.m4a")
	if err != nil {
		t.Error(err)
	}
	pretty.Println(music)
}