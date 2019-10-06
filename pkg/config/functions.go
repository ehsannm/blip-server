package config

import (
	"github.com/spf13/viper"
	"time"
)

/*
   Creation Time: 2019 - Oct - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
