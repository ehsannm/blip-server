package user

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/pkg/config"
	ronak "git.ronaksoftware.com/ronak/toolbox"
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
	userCol    *mongo.Collection
	redisCache *ronak.RedisCache
)

func InitMongo(c *mongo.Client) {
	userCol = c.Database(viper.GetString(config.ConfMongoDB)).Collection(config.ColUser)
}

func InitRedisCache(c *ronak.RedisCache) {
	redisCache = c
}

// easyjson:json
type User struct {
	ID        string `json:"id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Phone     string `json:"phone" bson:"phone"`
	Email     string `json:"email" bson:"email"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	Disabled  bool   `json:"disabled" bson:"disabled"`
	Premium   bool   `json:"premium" bson:"premium"`
}

func Save(user User) error {
	_, err := userCol.InsertOne(nil, user)
	return err
}

func readFromCache(userID string) (*User, error) {
	keyID := fmt.Sprintf("%s.%s", RkUser, userID)
	userBytes, err := redisCache.GetBytes(keyID)
	if err != nil || userBytes == nil {
		user, err := readFromDb(userID)
		if err != nil {
			return nil, err
		}
		userBytes, err = user.MarshalJSON()
		if err != nil {
			return nil, err
		}
		_, _ = redisCache.Set(keyID, userBytes)
		return user, nil
	}

	user := new(User)
	err = user.UnmarshalJSON(userBytes)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func readFromDb(userID string) (*User, error) {
	user := new(User)
	res := userCol.FindOne(nil, bson.M{"_id": userID}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Get(userID string) (*User, error) {
	return readFromCache(userID)
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
