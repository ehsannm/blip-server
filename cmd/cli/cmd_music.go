package main

import (
	"fmt"
	"github.com/spf13/cobra"
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

	SearchByProxyCmd.Flags().String(FlagFilePath, "", "")

	MusicCmd.AddCommand(SearchByProxyCmd)
}

var MusicCmd = &cobra.Command{
	Use: "Music",
}

var SearchByProxyCmd = &cobra.Command{
	Use: "SearchByProxyCmd",
	Run: func(cmd *cobra.Command, args []string) {
		err := sendFile("music/search_by_proxy", cmd.Flag(FlagFilePath).Value.String(), true)
		if err != nil {
			fmt.Println(err)
		}

	},
}
