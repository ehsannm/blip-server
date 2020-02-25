package vas

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
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
	vasLogCol = c.Database(config.DbMain).Collection(config.ColLogVas)
}
