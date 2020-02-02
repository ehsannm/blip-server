package main

import (
	"bytes"
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"git.ronaksoftware.com/blip/server/pkg/store"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
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
	RootCmd.AddCommand(StoreCmd)
	StoreCmd.AddCommand(SaveStoreCmd, GetStoreCmd)

	SaveStoreCmd.Flags().Int64(FlagStoreID, tools.RandomInt64(0), "")
	SaveStoreCmd.Flags().String(FlagStoreDsn, "", "")
	SaveStoreCmd.Flags().Int(FlagCapacity, 1000, "in megabytes")
	SaveStoreCmd.Flags().String(FlagRegion, "IR", "")

	GetStoreCmd.Flags().String(FlagStoreIDs, "", "")
}

var StoreCmd = &cobra.Command{
	Use: "Music",
}

var SaveStoreCmd = &cobra.Command{
	Use: "Save",
	Run: func(cmd *cobra.Command, args []string) {
		req := store.SaveStoreReq{
			StoreID:  tools.StrToInt64(cmd.Flag(FlagStoreID).Value.String()),
			Dsn:      cmd.Flag(FlagStoreDsn).Value.String(),
			Region:   cmd.Flag(FlagRegion).Value.String(),
			Capacity: int(tools.StrToInt64(cmd.Flag(FlagCapacity).Value.String())),
		}
		reqBytes, _ := req.MarshalJSON()
		_, err := sendHttp(http.MethodPost, "store/save", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var GetStoreCmd = &cobra.Command{
	Use: "Get",
	Run: func(cmd *cobra.Command, args []string) {
		storeIDsStr := cmd.Flag(FlagStoreIDs).Value.String()
		storeIDs := make([]int64, 0, 10)
		for _, storeIDstr := range strings.Split(storeIDsStr, ",") {
			storeIDs = append(storeIDs, tools.StrToInt64(strings.TrimSpace(storeIDstr)))
		}
		req := store.GetStoreReq{
			StoreIDs: storeIDs,
		}
		reqBytes, _ := req.MarshalJSON()
		_, err := sendHttp(http.MethodPost, "store/get", ContentTypeJSON, bytes.NewBuffer(reqBytes), true)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
