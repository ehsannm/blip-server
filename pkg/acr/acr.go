package acr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/config"

	"go.uber.org/zap"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Oct - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/


func IdentifyByFile(fileAddr string) (*Music, error) {
	f, err := os.Open(fileAddr)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Warn("Error On IdentifyByFile", zap.Error(err), zap.String("Path", fileAddr))
		return nil, err
	}

	t := time.Now().Unix() * 1000
	stringToSign := fmt.Sprintf("POST\n/v1/identify\n%s\naudio\n1\n%d", accessKey, t)

	hm := hmac.New(func() hash.Hash { return sha1.New() }, tools.StrToByte(accessSecret))
	hm.Write(tools.StrToByte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hm.Sum(nil))

	c := http.Client{
		Timeout: config.HttpRequestTimeout,
	}
	values := url.Values{}
	values.Set("access_key", accessKey)
	values.Set("data_type", "audio")
	values.Set("sample_bytes", fmt.Sprintf("%d", len(fileBytes)))
	values.Set("sample", base64.StdEncoding.EncodeToString(fileBytes))
	values.Set("signature_version", "1")
	values.Set("signature", signature)
	values.Set("timestamp", fmt.Sprintf("%d", t))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/identify", baseUrl), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	music := new(Music)
	err = json.Unmarshal(resBytes, music)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func IdentifyByByteString(fileBytes []byte) (*Music, error) {
	t := time.Now().Unix() * 1000
	stringToSign := fmt.Sprintf("POST\n/v1/identify\n%s\naudio\n1\n%d", accessKey, t)

	hm := hmac.New(func() hash.Hash { return sha1.New() }, tools.StrToByte(accessSecret))
	hm.Write(tools.StrToByte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(hm.Sum(nil))

	c := http.Client{
		Timeout: time.Second * 3,
	}
	values := url.Values{}
	values.Set("access_key", accessKey)
	values.Set("data_type", "audio")
	values.Set("sample_bytes", fmt.Sprintf("%d", len(fileBytes)))
	values.Set("sample", base64.StdEncoding.EncodeToString(fileBytes))
	values.Set("signature_version", "1")
	values.Set("signature", signature)
	values.Set("timestamp", fmt.Sprintf("%d", t))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/identify", baseUrl), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	music := new(Music)
	err = json.Unmarshal(resBytes, music)
	if err != nil {
		return nil, err
	}

	return music, nil
}
