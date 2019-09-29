package user

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
	userCol *mongo.Collection
)

func InitMongo(c *mongo.Client) {
	userCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColUser)
}


type User struct {
	ID        string `json:"id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Phone     string `json:"phone" bson:"phone"`
	Email     string `json:"email" bson:"email"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	Disabled  bool   `json:"disabled" bson:"disabled"`
}


func Save(user User) error {
	_, err := userCol.InsertOne(nil, user)
	return err
}

func Get(userID string) (*User, error) {
	user := new(User)
	res := userCol.FindOne(nil, bson.M{"_id": userID}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetByPhone(phone string) (*User, error) {
	user := new(User)
	res := userCol.FindOne(nil, bson.M{"phone": phone}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

