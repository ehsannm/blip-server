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
	ConfLogLevel            = "LOG_LEVEL"
	ConfMongoUrl            = "MONGO_URL"
	ConfMongoDB             = "MONGO_DB"
	ConfMongoRequestTimeout = "MONGO_REQUEST_TIMEOUT"
	ConfRedisUrl            = "REDIS_URL"
	ConfRedisPass           = "REDIS_PASS"

	// VAS Saba Configs
	ConfVasSabaServiceBaseUrl = "VAS_SABA_SERVICE_BASE_URL"
	ConfVasSabaServiceName    = "VAS_SABA_SERVICE_NAME"
	ConfVasSabaServiceToken   = "VAS_SABA_SERVICE_TOKEN"

	// ACR Configs
	ConfACRAccessKey    = "ACR_ACCESS_KEY"
	ConfACRAccessSecret = "ACR_ACCESS_SECRET"
	ConfACRBaseUrl      = "ACR_BASE_URL"
)

func init() {
	viper.SetEnvPrefix("BLIP")
	viper.AutomaticEnv()

	pflag.Int(ConfLogLevel, 0, "")
	pflag.String(ConfMongoUrl, "mongodb://localhost:27017", "")
	pflag.String(ConfMongoDB, "blip", "")
	pflag.String(ConfRedisUrl, "localhost:6379", "")
	pflag.String(ConfRedisPass, "ehsan2374", "")
	pflag.String(ConfVasSabaServiceName, "test", "")
	pflag.String(ConfVasSabaServiceToken, "stuimxfhyy", "")
	pflag.String(ConfVasSabaServiceBaseUrl, "http://api.sabaeco.com", "")

	pflag.String(ConfACRBaseUrl, "http://identify-eu-west-1.acrcloud.com", "")
	pflag.String(ConfACRAccessKey, "7f808c9dcbb700bf7018ffc92c49ff93", "")
	pflag.String(ConfACRAccessSecret, "EpbFGBwZcFtDUH4OSxPMj6247nb5WIy6yaTbIOiq", "")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}
