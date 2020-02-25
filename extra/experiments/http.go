package main

import (
	"fmt"
	"strings"
)

/*
   Creation Time: 2019 - Oct - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func main() {
	x := "1|2|3"
	fmt.Println(len(strings.Split(x, "|")))
	x = "1||2"
	fmt.Println(len(strings.Split(x, "|")))
}
