package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/crawler"
	"github.com/spf13/cobra"
	"net/http"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	RootCmd.AddCommand(CrawlerCmd)
	CrawlerCmd.AddCommand(CrawlerSaveCmd)
	CrawlerSaveCmd.Flags().String(FlagSource, "MeloBit", "")
	CrawlerSaveCmd.Flags().String(FlagName, "", "")
	CrawlerSaveCmd.Flags().String(FlagUrl, "http://ws.blipapi.xyz/crowler/melobit/", "")
}

var CrawlerCmd = &cobra.Command{
	Use: "Crawler",
}

var CrawlerSaveCmd = &cobra.Command{
	Use: "Save",
	Run: func(cmd *cobra.Command, args []string) {
		req := crawler.SaveReq{
			Url:    cmd.Flag(FlagUrl).Value.String(),
			Name:   cmd.Flag(FlagName).Value.String(),
			Source: cmd.Flag(FlagSource).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		fmt.Println(string(reqBytes))
		_, err := sendHttp(http.MethodPost, "crawler/save", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}

	},
}
