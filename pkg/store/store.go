package store

import (
	"errors"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetUploadStream(bucketName string, songID primitive.ObjectID) (int64, *gridfs.UploadStream, error) {
	var mongoClient *mongo.Client
	var storeID int64
	storesMtx.RLock()
	for sID, x := range storeBuckets {
		if x != nil {
			mongoClient = x
			storeID = sID
			break
		}
	}
	storesMtx.RUnlock()
	if mongoClient == nil {
		return storeID, nil, errors.New("no connection exists")
	}
	bucket, err := gridfs.NewBucket(mongoClient.Database(config.DbStore), options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return storeID, nil, err
	}
	stream, err := bucket.OpenUploadStreamWithID(songID, songID.Hex())
	return storeID, stream, err
}

func Download(storeID int64, bucketName string, songID primitive.ObjectID, dst io.Writer) error {
	storesMtx.RLock()
	mongoClient := storeBuckets[storeID]
	storesMtx.RUnlock()
	if mongoClient == nil {
		return errors.New("no connection exists")
	}
	bucket, err := gridfs.NewBucket(mongoClient.Database(config.DbStore), options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}
	_, err = bucket.DownloadToStream(songID, dst)
	return err
}
