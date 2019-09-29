package config

import (
	"github.com/spf13/pflag"
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

const (
	ConfMongoUrl            = "MONGO_URL"
	ConfMongoDB             = "MONGO_DB"
	ConfMongoRequestTimeout = "MONGO_REQUEST_TIMEOUT"
	ConfSmsServiceBaseUrl   = "SMS_SERVICE_BASE_URL"
	ConfSmsServiceName      = "SMS_SERVICE_NAME"
	ConfSmsServiceToken     = "SMS_SERVICE_TOKEN"
)

func init() {
	viper.SetEnvPrefix("BLIP")
	viper.AutomaticEnv()
	// viper.SetDefault(ConfSmsServiceToken, "stuimxfhyy")
	// viper.SetDefault(ConfSmsServiceName, "test")

	pflag.String(ConfMongoUrl, "mongodb://localhost:27017", "")
	pflag.String(ConfMongoDB, "blip", "")
	pflag.Duration(ConfMongoRequestTimeout, time.Second*3, "")
	pflag.String(ConfSmsServiceName, "test", "")
	pflag.String(ConfSmsServiceToken, "stuimxfhyy", "")
	pflag.String(ConfSmsServiceBaseUrl, "http://api.sabaeco.com", "")
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}
