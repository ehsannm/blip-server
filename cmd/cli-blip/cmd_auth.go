package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"github.com/spf13/cobra"
	"net/http"
)

func init() {
	RootCmd.AddCommand(AuthCmd)
	AuthCmd.AddCommand(SendCodeCmd, LoginCmd, RegisterCmd, CreateAccessKeyCmd)
	AuthCmd.PersistentFlags().String(FlagPhone, "", "")
	AuthCmd.PersistentFlags().String(FlagPhoneCode, "", "")
	AuthCmd.PersistentFlags().String(FlagPhoneCodeHash, "", "")
	AuthCmd.PersistentFlags().String(FlagOtpID, "", "")
	RegisterCmd.Flags().String(FlagUsername, "", "")
	CreateAccessKeyCmd.Flags().Bool(FlagPermRead, true, "")
	CreateAccessKeyCmd.Flags().Bool(FlagPermWrite, false, "")
	CreateAccessKeyCmd.Flags().Bool(FlagPermAdmin, false, "")
	CreateAccessKeyCmd.Flags().Int64(FlagPeriod, 0, "")
	CreateAccessKeyCmd.Flags().String(FlagAppName, "", "")
}

var AuthCmd = &cobra.Command{
	Use: "Auth",
}

var CreateAccessKeyCmd = &cobra.Command{
	Use: "CreateAccessKey",
	Run: func(cmd *cobra.Command, args []string) {

		req := auth.CreateAccessToken{
			Period:  tools.StrToInt64(cmd.Flag(FlagPeriod).Value.String()),
			AppName: cmd.Flag(FlagAppName).Value.String(),
		}
		if b, _ := cmd.Flags().GetBool(FlagPermRead); b {
			req.Permissions = append(req.Permissions, "read")
		}
		if b, _ := cmd.Flags().GetBool(FlagPermWrite); b {
			req.Permissions = append(req.Permissions, "write")
		}
		if b, _ := cmd.Flags().GetBool(FlagPermAdmin); b {
			req.Permissions = append(req.Permissions, "admin")
		}

		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/create", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
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
		_, err := sendHttp(http.MethodPost, "auth/send_code", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var LoginCmd = &cobra.Command{
	Use:   "Login",
	Short: "if user has been already registered, then just login to server",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.LoginReq{
			PhoneCode:     cmd.Flag(FlagPhoneCode).Value.String(),
			PhoneCodeHash: cmd.Flag(FlagPhoneCodeHash).Value.String(),
			Phone:         cmd.Flag(FlagPhone).Value.String(),
		}
		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/login", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var RegisterCmd = &cobra.Command{
	Use:   "Register",
	Short: "if user is a new one, then it registers users in the server",
	Run: func(cmd *cobra.Command, args []string) {
		req := auth.RegisterReq{
			PhoneCode:     cmd.Flag(FlagPhoneCode).Value.String(),
			PhoneCodeHash: cmd.Flag(FlagPhoneCodeHash).Value.String(),
			Phone:         cmd.Flag(FlagPhone).Value.String(),
		}

		if cmd.Flag(FlagUsername) != nil {
			req.Username = cmd.Flag(FlagUsername).Value.String()
		}

		reqBytes, _ := json.Marshal(req)
		_, err := sendHttp(http.MethodPost, "auth/register", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
