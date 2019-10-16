package main

import (
	"fmt"
	"net/http"
	"net/url"
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
	v := url.Values{}
	v.Set("phone", "989121228718")
	req, err := http.NewRequest(http.MethodPost, "https://webhook.site/4c6e5117-b7fd-42ae-a1b8-1c4118e1510e", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = v
	req.Form = v
	err = req.ParseForm()
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
