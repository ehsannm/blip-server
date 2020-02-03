package session

import (
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

// Session
type Session struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	CreatedOn  int64  `json:"created_on" bson:"created_on"`
	LastAccess int64  `json:"last_access" bson:"last_access"`
	App        string `json:"app" bson:"app,omitempty"`
}

// Save inserts/updates session 's' into the database
func Save(s *Session) error {
	_, err := sessionCol.UpdateOne(nil, bson.M{"_id": s.ID}, bson.M{"$set": s}, options.Update().SetUpsert(true))
	return err
}

// Remove removes all the sessions of userID associated with appName
func Remove(userID, appName string) error {
	session := &Session{}
	res := sessionCol.FindOneAndDelete(nil,
		bson.M{
			"user_id": userID,
			"$or": bson.A{
				bson.M{"app": bson.M{"$exists": false}},
				bson.M{"app": appName},
			},
		})
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

// Get returns the session identified by sessionID, or return error
func Get(sessionID string) (*Session, error) {
	session := &Session{}
	res := sessionCol.FindOne(nil, bson.M{"_id": sessionID}, options.FindOne().SetMaxTime(config.MongoRequestTimeout))
	err := res.Decode(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// GetAll returns a list of all the sessions of the userID on different apps
func GetAll(userID string) ([]*Session, error) {
	sessions := make([]*Session, 0, 10)
	cur, err := sessionCol.Find(nil, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	for cur.Next(nil) {
		session := &Session{}
		err := cur.Decode(session)
		if err == nil {
			sessions = append(sessions, session)
		}
	}
	err = cur.Close(nil)
	return sessions, err
}
