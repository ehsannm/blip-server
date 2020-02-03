package store

import (
	"errors"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
)

/*
   Creation Time: 2020 - Feb - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// easyjson:json
// Store
type Store struct {
	ID       int64  `bson:"_id" json:"id"`
	Dsn      string `bson:"dsn" json:"dsn"`
	Capacity int    `bson:"cap" json:"capacity"`
	Region   string `bson:"region" json:"region"`
}

// GetStore returns a store identified by storeID
func get(storeID int64) *Store {
	storesMtx.RLock()
	s := stores[storeID]
	storesMtx.RUnlock()
	return s
}

func UploadSong(storeID int64, songID primitive.ObjectID, source io.Reader) error {
	storesMtx.RLock()
	conn := storeConns[storeID]
	storesMtx.RUnlock()
	if conn == nil {
		return errors.New("no connection exists")
	}
	bucket, err := gridfs.NewBucket(conn.Database(config.DbStore), options.GridFSBucket().SetName("songs"))
	if err != nil {
		return err
	}
	return bucket.UploadFromStreamWithID(songID, songID.Hex(), source)
}

func DownloadSong(storeID int64, songID primitive.ObjectID, dst io.Writer) error {
	storesMtx.RLock()
	conn := storeConns[storeID]
	storesMtx.RUnlock()
	if conn == nil {
		return errors.New("no connection exists")
	}
	bucket, err := gridfs.NewBucket(conn.Database(config.DbStore), options.GridFSBucket().SetName("songs"))
	if err != nil {
		return err
	}
	_, err = bucket.DownloadToStream(songID, dst)
	return err
}
