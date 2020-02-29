package help

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
   Creation Time: 2020 - Feb - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

//go:generate rm -f *_easyjson.go
//go:generate easyjson messages.go
var (
	helpCol     *mongo.Collection
	feedbackCol *mongo.Collection
)

func InitMongo(c *mongo.Client) {
	helpCol = c.Database(config.DbMain).Collection(config.ColHelp)
	feedbackCol = c.Database(config.DbMain).Collection(config.ColFeedback)
}

func Init() {
	// load default configs from database
	loadDefaultConfig()
}
