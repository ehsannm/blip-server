package main

import (
	"fmt"
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
	AdminCmd.AddCommand(UnsubscribeCmd, MigrateLegacyDB, MigrateLegacyDBStats)
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

var MigrateLegacyDB = &cobra.Command{
	Use: "MigrateLegacyDB",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodPost, "admin/migrate_legacy_db", ContentTypeJSON, nil, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var MigrateLegacyDBStats = &cobra.Command{
	Use: "MigrateLegacyDBStats",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := sendHttp(http.MethodGet, "admin/migrate_legacy_db_stats", ContentTypeJSON, nil, true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
