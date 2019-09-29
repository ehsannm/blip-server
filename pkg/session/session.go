package session

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	sessionCol *mongo.Collection
)

func InitMongo(c *mongo.Client) {
	sessionCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColSession)
}

type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	CreatedOn  int64  `json:"created_on" bson:"created_on"`
	LastAccess int64  `json:"last_access" bson:"last_access"`
}


func Save(session Session) error {
	_, err := sessionCol.InsertOne(nil, session)
	return err
}

func Get(sessionID string) (*Session, error) {
	session := new(Session)
	res := sessionCol.FindOne(nil, bson.M{"_id": sessionID}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
