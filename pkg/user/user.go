package user

import (
	"git.ronaksoftware.com/blip/server/internal/redis"
	"git.ronaksoftware.com/blip/server/pkg/config"

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
	redisCache *redis.Cache
)

// easyjson:json
type User struct {
	ID        string `json:"id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Phone     string `json:"phone" bson:"phone"`
	Email     string `json:"email" bson:"email"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	Disabled  bool   `json:"disabled" bson:"disabled"`
	VasPaid   bool   `json:"vas_paid" bson:"vas_paid"`
}

// Save stores the object 'u' in the database, replace the old one if already exists
func Save(u *User) error {
	_, err := userCol.UpdateOne(nil, bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"username":   u.Username,
		"phone":      u.Phone,
		"email":      u.Email,
		"created_on": u.CreatedOn,
		"disabled":   u.Disabled,
		"vas_paid":   u.VasPaid,
	}}, options.Update().SetUpsert(true))
	deleteFromCache(u.ID)
	return err
}

func readFromCache(userID string) (*User, error) {
	userBytes, err := redisCache.GetBytes(userID)
	if err != nil || userBytes == nil {
		user, err := readFromDb(userID)
		if err != nil {
			return nil, err
		}
		userBytes, err = user.MarshalJSON()
		if err != nil {
			return nil, err
		}
		_ = redisCache.Set(userID, userBytes)
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

func deleteFromCache(userID string) {
	_ = redisCache.Del(userID)
}

// Get returns user identified by 'userID'
func Get(userID string) (*User, error) {
	return readFromCache(userID)
}

// GetByPhone returns user identified by 'phone'
func GetByPhone(phone string) (*User, error) {
	user := &User{}
	res := userCol.FindOne(nil, bson.M{"phone": phone}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

