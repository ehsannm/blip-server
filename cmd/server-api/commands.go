package main

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"github.com/kataras/iris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
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
		c := _Mongo.Database(viper.GetString(config.MongoDB)).Collection(config.ColAuth)
		_, err := c.InsertOne(nil, auth.Auth{
			ID:          "ROOT",
			Permissions: []auth.Permission{auth.Admin},
			CreatedOn:   time.Now().Unix(),
			ExpiredOn:   0,
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}
