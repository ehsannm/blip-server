package gridfs

import "github.com/gobwas/pool/pbytes"

/*
   Creation Time: 2020 - Feb - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/


func init() {
	pbytes.DefaultPool = pbytes.New(128, 256 << 10)
}