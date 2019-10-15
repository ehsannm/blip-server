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
	fs.String(FlagPhone, "", "")
	fs.String(FlagPhoneCode, "", "")
	fs.String(FlagPhoneCodeHash, "", "")
	fs.String(FlagOtpID, "", "")

	RegisterCmd.Flags().String(FlagUsername, "", "")

	RootCmd.AddCommand(AuthCmd)
	AuthCmd.AddCommand(SendCodeCmd, LoginCmd, RegisterCmd)

}

var AuthCmd = &cobra.Command{
	Use: "Auth",
}

var SendCodeCmd = &cobra.Command{
	Use:   "SendCode",
	Short: "send sms code request",
	Long:  "send a sms code to the phone number for verification",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.SendCodeReq{
			Phone: cmd.Flag(FlagPhone).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/send_code", bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var LoginCmd = &cobra.Command{
	Use: "Login",
	Short: "if user has been already registered, then just login to server",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.LoginReq{
			PhoneCode:     cmd.Flag(FlagPhoneCode).Value.String(),
			PhoneCodeHash: cmd.Flag(FlagPhoneCodeHash).Value.String(),
			Phone:         cmd.Flag(FlagPhone).Value.String(),
			OperationID:   cmd.Flag(FlagOtpID).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/login", bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var RegisterCmd = &cobra.Command{
	Use: "Register",
	Short: "if user is a new one, then it registers users in the server",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.RegisterReq{
			PhoneCode:     cmd.Flag(FlagPhoneCode).Value.String(),
			PhoneCodeHash: cmd.Flag(FlagPhoneCodeHash).Value.String(),
			Phone:         cmd.Flag(FlagPhone).Value.String(),
			OperationID:   cmd.Flag(FlagOtpID).Value.String(),
		}

		if cmd.Flag(FlagUsername) != nil {
			req.Username = cmd.Flag(FlagUsername).Value.String()
		}

		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/register", bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
