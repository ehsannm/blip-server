package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/music"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"time"
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
	MusicCmd.AddCommand(SearchByProxyCmd, SearchByTextCmd, SearchResumeCmd, SearchBySoundCmd, DownloadCmd)

	SearchByProxyCmd.Flags().String(FlagFilePath, "", "")
	SearchBySoundCmd.Flags().String(FlagFilePath, "", "")
	SearchByTextCmd.Flags().String(FlagKeyword, "", "")
	DownloadCmd.Flags().String(FlagSongID, "", "")
	DownloadCmd.Flags().String(FlagFilePath, "./song.mp3", "")
}

var MusicCmd = &cobra.Command{
	Use: "Music",
}

var SearchByProxyCmd = &cobra.Command{
	Use: "SearchByProxy",
	Run: func(cmd *cobra.Command, args []string) {
		err := sendFile("music/search_by_proxy", "File", cmd.Flag(FlagFilePath).Value.String(), true)
		if err != nil {
			fmt.Println(err)
		}

	},
}

var SearchByTextCmd = &cobra.Command{
	Use: "SearchByText",
	Run: func(cmd *cobra.Command, args []string) {
		keyword := cmd.Flag(FlagKeyword).Value.String()
		keyword = strings.Join(strings.Split(keyword, ","), " ")
		req := music.SearchReq{
			Keyword: keyword,
		}
		reqBytes, _ := json.Marshal(req)
		res, err := sendHttp(http.MethodPost, "music/search/text", ContentTypeJSON, bytes.NewBuffer(reqBytes), false)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch res.Constructor {
		case music.CSearchResult:
			color.Green("Result: %s", res.Constructor)
			for _, s := range res.Payload.(map[string]interface{})["songs"].([]interface{}) {
				songX := s.(map[string]interface{})
				color.HiCyan("%s (%s) --> %s", songX["title"].(string), songX["artists"].(string), songX["id"])
			}
		default:
			color.Red("%s %v", res.Constructor, res.Payload)
		}

		// for {
		// 	res, err := sendHttp(http.MethodGet, "music/search", ContentTypeJSON, nil, false)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		continue
		// 	}
		// 	switch res.Constructor {
		// 	case music.CSearchResult:
		// 		color.Green("Result: %s", res.Constructor)
		// 		for _, s := range res.Payload.(map[string]interface{})["songs"].([]interface{}) {
		// 			songX := s.(map[string]interface{})
		// 			color.HiBlue("%s (%s) --> %s", songX["title"].(string), songX["artists"].(string), songX["id"])
		// 		}
		// 	default:
		// 		color.Red("%s %v", res.Constructor, res.Payload)
		// 	}
		// 	if res.Constructor == "err" {
		// 		break
		// 	}
		// }
	},
}

var SearchResumeCmd = &cobra.Command{
	Use: "SearchResume",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodGet, "music/search", ContentTypeJSON, nil, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var SearchBySoundCmd = &cobra.Command{
	Use: "SearchBySound",
	Run: func(cmd *cobra.Command, args []string) {
		err := sendFile("music/search/sound", "sound", cmd.Flag(FlagFilePath).Value.String(), true)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var DownloadCmd = &cobra.Command{
	Use: "Download",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/music/download/songs/%s", baseUrl, cmd.Flag(FlagSongID).Value.String())
		timeStart := time.Now()
		err := getFile(url, cmd.Flag(FlagFilePath).Value.String())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Download Time:", time.Now().Sub(timeStart))
	},
}
