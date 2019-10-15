package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

/*
   Creation Time: 2019 - Oct - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func sendHttp(method, urlSuffix string, reader io.Reader, formValues url.Values, print bool) (*msg.ResponseEnvelope, error) {
	c := http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", baseUrl, urlSuffix), reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error In Request: %v", err))
	}
	switch method {
	case http.MethodPost:
		req.PostForm = formValues
	}

	req.Header.Set(auth.HdrAccessKey, "ROOT")
	res, err := c.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error In Response: %v", err))
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error In Reading Response: %v", err))
	}

	x := msg.ResponseEnvelope{}
	err = json.Unmarshal(bodyBytes, &x)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error In Unmarshal Response: %v, %s", err, ronak.ByteToStr(bodyBytes)))
	}

	if print {
		fmt.Println(x.Constructor)
		_, _ = pretty.Println(x.Payload)

	}

	return &x, nil
}
