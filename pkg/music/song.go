package music

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

/*
   Creation Time: 2020 - Jan - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Song
// easyjson:json
type Song struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Genre          string             `bson:"genre" json:"genre"`
	Lyrics         string             `bson:"lyrics" json:"lyrics"`
	Artists        string             `bson:"artists" json:"artists"`
	CoverUrl       string             `bson:"cover_url" json:"cover_url"`
	SongUrl        string             `bson:"song_url" json:"song_url"`
	OriginCoverUrl string             `bson:"org_cover_url" json:"-"`
	OriginSongUrl  string             `bson:"org_song_url" json:"-"`
	Source         string             `bson:"source" json:"-"`
}

func DropAllSongs() error {
	return songCol.Drop(nil)
}

func SaveSong(s *Song) (primitive.ObjectID, error) {
	s.ID = primitive.NewObjectID()
	_, err := songCol.InsertOne(nil, s, options.InsertOne())
	if err != nil {
		return primitive.NilObjectID, err
	}
	return s.ID, nil
}

func GetSong(songID primitive.ObjectID) (*Song, error) {
	res := songCol.FindOne(nil, bson.M{"_id": songID}, options.FindOne())
	err := res.Err()
	if err != nil {
		return nil, res.Err()
	}
	s := &Song{}
	err = res.Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
