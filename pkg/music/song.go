package music

import (
	"encoding/hex"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"github.com/gobwas/pool/pbytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	UniqueKey      string             `bson:"unique_key" json:"-"`
	Title          string             `bson:"title" json:"title"`
	Genre          string             `bson:"genre" json:"genre"`
	Lyrics         string             `bson:"lyrics" json:"lyrics"`
	Artists        string             `bson:"artists" json:"artists"`
	Cdn            string             `bson:"cdn" json:"-"`
	OriginCoverUrl string             `bson:"org_cover_url" json:"-"`
	OriginSongUrl  string             `bson:"org_song_url" json:"-"`
	Source         string             `bson:"source" json:"-"`
}

func GenerateUniqueKey(songX *Song) string {
	uniqueKeyArgs := pbytes.GetCap(len(songX.Title) + len(songX.Artists))
	uniqueKeyArgs = append(uniqueKeyArgs, tools.StrToByte(songX.Title)...)
	uniqueKeyArgs = append(uniqueKeyArgs, tools.StrToByte(songX.Artists)...)
	id, _ := tools.Sha256(uniqueKeyArgs)
	return hex.EncodeToString(id[:])
}

// DropAllSongs drop all the songs from the database
func DropAllSongs() error {
	return songCol.Drop(nil)
}

// SaveSong saves/replaces the song 's' to the database
func SaveSong(songX *Song) (primitive.ObjectID, error) {
	songX.ID = primitive.NewObjectID()
	_, err := songCol.InsertOne(nil, songX)
	if err != nil {
		return primitive.NilObjectID, err
	}
	err = UpdateLocalIndex(songX)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return songX.ID, nil
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
