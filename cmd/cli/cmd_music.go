package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/music"
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
	RootCmd.AddCommand(MusicCmd)
	MusicCmd.AddCommand(SearchByProxyCmd, SearchByTextCmd, SearchResumeCmd)

	SearchByProxyCmd.Flags().String(FlagFilePath, "", "")
	SearchByTextCmd.Flags().String(FlagKeyword, "", "")
}

var MusicCmd = &cobra.Command{
	Use: "Music",
}

var SearchByProxyCmd = &cobra.Command{
	Use: "SearchByProxy",
	Run: func(cmd *cobra.Command, args []string) {
		err := sendFile("music/search_by_proxy", cmd.Flag(FlagFilePath).Value.String(), true)
		if err != nil {
			fmt.Println(err)
		}

	},
}

var SearchByTextCmd = &cobra.Command{
	Use: "SearchByText",
	Run: func(cmd *cobra.Command, args []string) {
		req := music.SearchReq{
			Keyword: cmd.Flag(FlagKeyword).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "music/search_by_text", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var SearchResumeCmd = &cobra.Command{
	Use: "SearchResume",
	Run: func(cmd *cobra.Command, args []string) {
		req := music.SearchReq{
			Keyword: cmd.Flag(FlagKeyword).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "music/search_resume", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

