package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/msg"
	"git.ronaksoftware.com/blip/server/pkg/session"

	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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

var (
	accessToken string
	sessionID   string
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

func sendHttp(method, urlSuffix string, contentType string, reader io.Reader, print bool) (*msg.ResponseEnvelope, error) {
	c := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", baseUrl, urlSuffix), reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error In Request: %v", err))
	}
	req.Header.Set(auth.HdrAccessKey, accessToken)
	req.Header.Set(session.HdrSessionID, sessionID)
	req.Header.Set("Content-Type", contentType)
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
		return nil, errors.New(fmt.Sprintf("Error In Unmarshal Response: %v, %s", err, tools.ByteToStr(bodyBytes)))
	}

	if print {
		fmt.Println(res.Status, x.Constructor)
		_, _ = pretty.Println(x.Payload)

	}

	return &x, nil
}

func sendFile(urlSuffix string, formField, filename string, print bool) error {
	c := http.Client{
		Timeout: 30 * time.Second,
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile(formField, filename)
	if err != nil {
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, urlSuffix), bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Set(auth.HdrAccessKey, accessToken)
	req.Header.Set(session.HdrSessionID, sessionID)
	req.Header.Set("Content-Type", contentType)
	res, err := c.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Error In Response: %v", err))
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Error In Reading Response: %v", err))
	}
	if print {
		fmt.Println(tools.ByteToStr(bodyBytes))
	}
	return nil
}

func getFile(url string, filepath string) error {
	c := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set(auth.HdrAccessKey, accessToken)
	req.Header.Set(session.HdrSessionID, sessionID)
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		f, err := os.Create(filepath)
		if err != nil {
			return err
		}
		_, _ = io.Copy(f, res.Body)
		_ = f.Close()
		return nil
	default:

	}
	fmt.Println(res.Status, res.StatusCode)
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(tools.ByteToStr(bodyBytes))
	return nil

}
