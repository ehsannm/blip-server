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

var (
	mediumPool *pbytes.Pool
	largePool  *pbytes.Pool
)

func init() {
	mediumPool = pbytes.New(64<<10, 512<<10)
	largePool = pbytes.New(16<<20, 16<<20)
}
