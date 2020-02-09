package music

import (
	"context"
	"encoding/hex"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"github.com/gobwas/pool/pbytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
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
	StoreID        int64              `bson:"store_id" json:"-"`
	OriginCoverUrl string             `bson:"org_cover_url" json:"-"`
	OriginSongUrl  string             `bson:"org_song_url" json:"-"`
	Source         string             `bson:"source" json:"-"`
}

// GenerateUniqueKey returns a unique hash which help us to identify similar songs to prevent from double storing
// those songs in the database.
func GenerateUniqueKey(title, artists string) string {
	uniqueKeyArgs := pbytes.GetCap(len(title) + len(artists))
	uniqueKeyArgs = append(uniqueKeyArgs, tools.StrToByte(title)...)
	uniqueKeyArgs = append(uniqueKeyArgs, tools.StrToByte(artists)...)
	id, _ := tools.Sha256(uniqueKeyArgs)
	return hex.EncodeToString(id[:])
}

// DropAllSongs drop all the songs from the database
func DropAllSongs() error {
	return songCol.Drop(nil)
}

// DeleteSong deletes song from the database
func DeleteSong(songID primitive.ObjectID) error {
	_, err := songCol.DeleteOne(nil, bson.M{"_id": songID})
	if err != nil {
		return err
	}
	return songIndex.Delete(songID.Hex())
}

// SaveSong saves/replaces the song 'songX' to the database
func SaveSong(songX *Song) (primitive.ObjectID, error) {
	if songX.ID == primitive.NilObjectID {
		songX.ID = primitive.NewObjectID()
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
	defer cancelFunc()
	_, err := songCol.UpdateOne(
		ctx,
		bson.M{"_id": songX.ID}, bson.M{"$set": songX}, options.Update().SetUpsert(true))
	if err != nil {
		return primitive.NilObjectID, err
	}
	return songX.ID, nil
}

// GetSongByID returns a song identified by songID
func GetSongByID(songID primitive.ObjectID) (*Song, error) {
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

// GetSongByUniqueKey returns a song identified by uniqueKey
func GetSongByUniqueKey(uniqueKey string) (*Song, error) {
	res := songCol.FindOne(nil, bson.M{"unique_key": uniqueKey}, options.FindOne())
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

// GetManySongs returns a list of songs
func GetManySongs(songIDs []primitive.ObjectID) ([]*Song, error) {
	cur, err := songCol.Find(nil, bson.M{"_id": bson.M{"$in": songIDs}})
	if err != nil {
		return nil, err
	}
	songs := make([]*Song, 0, len(songIDs))
	for cur.Next(nil) {
		songX := &Song{}
		err = cur.Decode(songX)
		if err == nil {
			songs = append(songs, songX)
		}
	}
	err = cur.Close(nil)
	return songs, err
}

func ForEachSong(f func(songX *Song) bool) error {
	var lastID primitive.ObjectID
	cur, err := songCol.Find(context.Background(), bson.D{}, options.Find().SetNoCursorTimeout(true))
	if err != nil {
		return err
	}
	waitGroup := sync.WaitGroup{}
	defer waitGroup.Wait()
	rateLimit := make(chan struct{}, 20)
	for {
		for cur.Next(nil) {
			songX := &Song{}
			err = cur.Decode(songX)
			if err != nil {
				return err
			}
			lastID = songX.ID
			waitGroup.Add(1)
			rateLimit <- struct{}{}
			go func(songX *Song) {
				f(songX)
				waitGroup.Done()
				<-rateLimit
			}(songX)
		}
		if cur.Err() == nil {
			_ = cur.Close(nil)
			break
		}
		log.Warn("Error On Cursor", zap.Error(err))
		_ = cur.Close(nil)
		cur, err = songCol.Find(
			context.Background(),
			bson.M{"_id": bson.M{"$gt": lastID}},
			options.Find().SetNoCursorTimeout(true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
