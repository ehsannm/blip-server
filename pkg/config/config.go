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
	TestMode        = "TEST_MODE"
	HttpProxy       = "HTTP_PROXY"
	LogLevel        = "LOG_LEVEL"
	MongoUrl        = "MONGO_URL"
	RedisUrl        = "REDIS_URL"
	RedisPass       = "REDIS_PASS"
	MagicPhone      = "MAGIC_PHONE"
	MagicPhoneCode  = "MAGIC_PHONE_CODE"
	SongsIndexDir   = "SONGS_INDEX_DIR"
	ProfilerEnabled = "PROFILER_ENABLED"
	ProfilerPort    = "PROFILER_PORT"

	// VAS Saba Configs
	VasSabaServiceBaseUrl = "VAS_SABA_SERVICE_BASE_URL"
	VasSabaServiceName    = "VAS_SABA_SERVICE_NAME"
	VasSabaServiceToken   = "VAS_SABA_SERVICE_TOKEN"

	// ACR Configs
	ACRAccessKey    = "ACR_ACCESS_KEY"
	ACRAccessSecret = "ACR_ACCESS_SECRET"
	ACRBaseUrl      = "ACR_BASE_URL"

	// Sms Configs
	SmsADPUrl       = "SMS_ADP_URL"
	SmsAdpUser      = "SMS_ADP_USER"
	SmsAdpPass      = "SMS_ADP_PASS"
	SmsAdpPhone     = "SMS_ADP_PHONE"
	SmsPayamakUrl   = "SMS_PAYAMAK_URL"
	SmsPayamakUser  = "SMS_PAYAMAK_USER"
	SmsPayamakPass  = "SMS_PAYAMAK_PASS"
	SmsPayamakPhone = "SMS_PAYAMAK_PHONE"
)

func Init() {
	viper.SetEnvPrefix("BLIP")
	viper.AutomaticEnv()

	pflag.Bool(ProfilerEnabled, false, "")
	pflag.Int(ProfilerPort, 6060, "")
	pflag.String(HttpProxy, "***REMOVED***", "")
	pflag.Bool(TestMode, false, "")
	pflag.Int(LogLevel, 0, "")
	pflag.String(MongoUrl, "mongodb://localhost:27017", "")
	pflag.String(RedisUrl, "localhost:6379", "")
	pflag.String(RedisPass, "ehsan2374", "")
	pflag.String(MagicPhone, "237400", "")
	pflag.String(MagicPhoneCode, "2374", "")
	pflag.String(SongsIndexDir, ".", "")

	pflag.String(VasSabaServiceName, "test", "")
	pflag.String(VasSabaServiceToken, "stuimxfhyy", "")
	pflag.String(VasSabaServiceBaseUrl, "http://api.sabaeco.com", "")

	pflag.String(ACRBaseUrl, "http://identify-eu-west-1.acrcloud.com", "")
	pflag.String(ACRAccessKey, "7f808c9dcbb700bf7018ffc92c49ff93", "")
	pflag.String(ACRAccessSecret, "EpbFGBwZcFtDUH4OSxPMj6247nb5WIy6yaTbIOiq", "")

	pflag.String(SmsADPUrl, "https://ws.adpdigital.com/url/send", "")
	pflag.String(SmsAdpUser, "ronak", "")
	pflag.String(SmsAdpPass, "E2e2374k19743", "")
	pflag.String(SmsAdpPhone, "98200049112", "")

	pflag.String(SmsPayamakUrl, "http://37.228.138.118/post/sendsms.ashx", "")
	pflag.String(SmsPayamakUser, "9122139561", "")
	pflag.String(SmsPayamakPass, "539ma6", "")
	pflag.String(SmsPayamakPhone, "50001060010920", "")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}
