package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
   Creation Time: 2019 - Sep - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	ConfMongoUrl = "MONGO_URL"
	ConfMongoDB  = "MONGO_DB"
)

func init() {
	pflag.String(ConfMongoUrl, "mongodb://localhost:27017", "")
	pflag.String(ConfMongoDB, "blip", "")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}
