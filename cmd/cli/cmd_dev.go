package main

import (
	"fmt"
	"github.com/spf13/cobra"
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

func init() {
	RootCmd.AddCommand(DevCmd)
	DevCmd.AddCommand(UnsubscribeCmd)
}

var DevCmd = &cobra.Command{
	Use: "Dev",
}

var UnsubscribeCmd = &cobra.Command{
	Use: "Unsubscribe",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodPost, "dev/unsubscribe", nil,
			url.Values{
				"phone": []string{cmd.Flag(FlagPhone).Value.String()},
			},
			true,
		)
		if err != nil {
			fmt.Println("HERE", err)
			return
		}
	},
}
