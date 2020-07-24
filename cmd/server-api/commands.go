package main

import (
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"github.com/kataras/iris/v12"
	"github.com/spf13/cobra"
)

/*
   Creation Time: 2019 - Sep - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var Root = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		app := initServer()
		err := app.Run(iris.Addr(":80"), iris.WithOptimizations)
		if err != nil {
			log.Warn(err.Error())
		}
	},
}

var InitDB = &cobra.Command{
	Use: "initDB",
	Run: func(cmd *cobra.Command, args []string) {
		// Create Root if there is no auth exist in the database

	},
}
