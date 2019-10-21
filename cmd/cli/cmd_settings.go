package main

import (
	"fmt"
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
	SetSessionIDCmd.Flags().String(FlagSessionID, "", "")

	SettingsCmd.AddCommand(SetAccessTokenCmd, SetSessionIDCmd, GetAccessTokenCmd, GetSessionIDCmd)
}

var SettingsCmd = &cobra.Command{
	Use: "Settings",
}

var SetAccessTokenCmd = &cobra.Command{
	Use: "SetAccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		_ = ioutil.WriteFile(".blip-accessToken", ronak.StrToByte(cmd.Flag(FlagAccessToken).Value.String()), os.ModePerm)
	},
}

var GetAccessTokenCmd = &cobra.Command{
	Use: "GetAccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		tokenBytes, _ := ioutil.ReadFile(".blip-accessToken")
		fmt.Println(ronak.ByteToStr(tokenBytes))
	},
}

var SetSessionIDCmd = &cobra.Command{
	Use: "SetSessionID",
	Run: func(cmd *cobra.Command, args []string) {
		_ = ioutil.WriteFile(".blip-session", ronak.StrToByte(cmd.Flag(FlagSessionID).Value.String()), os.ModePerm)
	},
}

var GetSessionIDCmd = &cobra.Command{
	Use: "GetSessionID",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := ioutil.ReadFile(".blip-session")
		fmt.Println(ronak.ByteToStr(id))
	},
}
