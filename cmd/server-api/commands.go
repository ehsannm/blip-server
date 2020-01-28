package main

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/auth"
	"git.ronaksoftware.com/blip/server/pkg/config"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"github.com/kataras/iris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
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
		// Create Root if there is no auth exist in the database
		c := _Mongo.Database(viper.GetString(config.MongoDB)).Collection(config.ColAuth)
		cnt, err := c.CountDocuments(nil, bson.D{})
		if err != nil {
			fmt.Println(err)
			return
		}
		if cnt == 0 {
			_, err := c.InsertOne(nil, auth.Auth{
				ID:          "ROOT",
				Permissions: []auth.Permission{auth.Admin},
				CreatedOn:   time.Now().Unix(),
				ExpiredOn:   0,
			})
			if err != nil {
				fmt.Println(err)
			}
		}


		// Create Magic User
		c = _Mongo.Database(config.GetString(config.MongoDB)).Collection(config.ColUser)
		_, err = c.InsertOne(nil, user.User{
			ID:        "MAGIC_USER",
			Username:  "MAGIC_USER",
			Phone:     "2374002374",
			Email:     "support@blip.fun",
			CreatedOn: 0,
			Disabled:  false,
			VasPaid:   true,
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}
