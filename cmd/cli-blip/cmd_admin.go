package main

import (
	"bytes"
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/admin"
	"git.ronaksoftware.com/blip/server/pkg/help"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strings"
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
	RootCmd.AddCommand(AdminCmd)
	AdminCmd.AddCommand(UnsubscribeCmd, HealthCheckDb, HealthCheckStore, HealthCheckStats, SetVasCmd, SetConfig, DeleteConfig)
	SetVasCmd.Flags().String(FlagUserID, "", "")
	SetVasCmd.Flags().Bool(FlagEnabled, false, "")
	SetConfig.Flags().String(FlagKey, "", "")
	SetConfig.Flags().String(FlagValue, "", "")
	DeleteConfig.Flags().String(FlagKey, "", "")
}

var AdminCmd = &cobra.Command{
	Use: "Admin",
}

var UnsubscribeCmd = &cobra.Command{
	Use: "Unsubscribe",
	Run: func(cmd *cobra.Command, args []string) {
		v := url.Values{}
		v.Set("phone", cmd.Flag(FlagPhone).Value.String())
		_, err := sendHttp(http.MethodPost, "dev/unsubscribe", ContentTypeForm,
			strings.NewReader(v.Encode()),
			true,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var HealthCheckDb = &cobra.Command{
	Use: "HealthCheckDb",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodPost, "admin/health_check_db", ContentTypeJSON, nil, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
var HealthCheckStore = &cobra.Command{
	Use: "HealthCheckStore",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodPost, "admin/health_check_store", ContentTypeJSON, nil, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var HealthCheckStats = &cobra.Command{
	Use: "HealthCheckStats",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := sendHttp(http.MethodGet, "admin/health_check_stats", ContentTypeJSON, nil, false)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch res.Constructor {
		case admin.CHealthCheckStats:
			v := res.Payload.(map[string]interface{})
			color.HiGreen("Scanned: %s", color.BlueString("%d", int(v["scanned"].(float64))))
			color.HiRed("CoverFixed: %s", color.BlueString("%d", int(v["cover_fixed"].(float64))))
			color.HiRed("SongFixed: %s", color.BlueString("%d", int(v["song_fixed"].(float64))))
		}

	},
}

var SetVasCmd = &cobra.Command{
	Use: "SetVas",
	Run: func(cmd *cobra.Command, args []string) {
		req := admin.SetVasReq{
			UserID:  cmd.Flag(FlagUserID).Value.String(),
			Enabled: cmd.Flag(FlagEnabled).Changed,
		}
		reqBytes, _ := req.MarshalJSON()
		_, err := sendHttp(http.MethodPost, "admin/vas", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var SetConfig = &cobra.Command{
	Use: "SetConfig",
	Run: func(cmd *cobra.Command, args []string) {
		req := help.SetDefaultConfig{
			Key:   cmd.Flag(FlagKey).Value.String(),
			Value: cmd.Flag(FlagValue).Value.String(),
		}
		reqBytes, _ := req.MarshalJSON()
		_, err := sendHttp(http.MethodPost, "help/config", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var DeleteConfig = &cobra.Command{
	Use: "DeleteConfig",
	Run: func(cmd *cobra.Command, args []string) {
		req := help.UnsetDefaultConfig{
			Key: cmd.Flag(FlagKey).Value.String(),
		}
		reqBytes, _ := req.MarshalJSON()
		_, err := sendHttp(http.MethodDelete, "help/config", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
