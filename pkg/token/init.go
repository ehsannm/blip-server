package token

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   Creation Time: 2020 - Jan - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func InitMongo(c *mongo.Client) {
	tokenCol = c.Database(viper.GetString(config.MongoDB)).Collection(config.ColToken)
}

func Init() {
	tokenCache = make(map[string]Token, 100000)
}
