package main

import (
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
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
	RootCmd.AddCommand(SettingsCmd)

	SetAccessTokenCmd.Flags().String(FlagAccessToken, "", "")

	SettingsCmd.AddCommand(SetAccessTokenCmd)
}

var SettingsCmd = &cobra.Command{
	Use: "Settings",
}


var SetAccessTokenCmd = &cobra.Command{
	Use:                        "SetAccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		_ = ioutil.WriteFile(".blip-accessToken", ronak.StrToByte(cmd.Flag(FlagAccessToken).Value.String()), os.ModePerm)
	},
}