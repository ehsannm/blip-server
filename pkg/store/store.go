package store

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

func GetUploadStreamForSong(songID primitive.ObjectID) (int64, *gridfs.UploadStream, error) {
	var bucket *gridfs.Bucket
	var storeID int64
	storesMtx.RLock()
	for sID, x := range storeBuckets {
		if x != nil {
			bucket = x
			storeID = sID
			break
		}
	}
	storesMtx.RUnlock()
	if bucket == nil {
		return storeID, nil, errors.New("no connection exists")
	}
	stream, err := bucket.OpenUploadStreamWithID(songID, songID.Hex())
	return storeID, stream, err
}

func DownloadSong(storeID int64, songID primitive.ObjectID, dst io.Writer) error {
	storesMtx.RLock()
	bucket := storeBuckets[storeID]
	storesMtx.RUnlock()
	if bucket == nil {
		return errors.New("no connection exists")
	}
	_, err := bucket.DownloadToStream(songID, dst)
	return err
}
