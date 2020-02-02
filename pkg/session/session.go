package session

import (
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

/*
   Creation Time: 2019 - Sep - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/


// Session
type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	CreatedOn  int64  `json:"created_on" bson:"created_on"`
	LastAccess int64  `json:"last_access" bson:"last_access"`
}

func Save(s Session) error {
	_, err := sessionCol.UpdateOne(nil, bson.M{"_id": s.ID}, bson.M{"$set": bson.M{
		"user_id":     s.UserID,
		"created_on":  s.CreatedOn,
		"last_access": s.LastAccess,
	}}, options.Update().SetUpsert(true))
	return err
}

func RemoveAll(userID string) error {
	session := &Session{}
	res := sessionCol.FindOneAndDelete(nil, bson.M{"user_id": userID})
	if res.Err() == mongo.ErrNoDocuments {
		return nil
	}
	err := res.Decode(session)
	if err != nil {
		return err
	}
	sessionCacheMtx.Lock()
	delete(sessionCache, session.ID)
	sessionCacheMtx.Unlock()
	return err
}

func Get(sessionID string) (*Session, error) {
	session := &Session{}
	res := sessionCol.FindOne(nil, bson.M{"_id": sessionID}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
