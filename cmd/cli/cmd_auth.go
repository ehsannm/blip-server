package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
	fs := RootCmd.PersistentFlags()
	fs.String(FlagPhone, "", FlagPhone)
	fs.String(FlagPhoneCode, "", FlagPhoneCode)
	fs.String(FlagPhoneCodeHash, "", FlagPhoneCodeHash)
	fs.String(FlagOtpID, "", FlagOtpID)


	// markFlagRequired(LoginCmd, FlagPhone, FlagPhoneCode, FlagPhoneCodeHash, FlagOtpID)

	RootCmd.AddCommand(SendCodeCmd, LoginCmd)
	SendCodeCmd.AddCommand(Sub1Cmd)
}


var SendCodeCmd = &cobra.Command{
	Use:   "SendCodeCmd",
	Short: "send sms code request",
	Long:  "send a sms code to the phone number for verification",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.SendCodeReq{
			Phone: cmd.Flag(FlagPhone).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		res, err := sendHttp(http.MethodPost, "auth/send_code", bytes.NewBuffer(reqBytes))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

var LoginCmd = &cobra.Command{
	Use: "LoginCmd",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flag("phone")

	},
}


var Sub1Cmd = &cobra.Command{
	Use: "Sub1",
	Run: func(cmd *cobra.Command, args []string) {

	},
	
}