package store

import (
	"errors"
	log "git.ronaksoftware.com/blip/server/internal/logger"
	"git.ronaksoftware.com/blip/server/pkg/config"
	"git.ronaksoftware.com/blip/server/pkg/store/gridfs"
	"github.com/gobwas/pool/pbytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io"
	"sync"
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

// DropAll drop all the stores from the database
func DropAll() error {
	err := storeCol.Drop(nil)
	if err != nil {
		return err
	}
	storesMtx.Lock()
	for k := range stores {
		delete(stores, k)
		delete(storeConns, k)
	}
	storesMtx.Unlock()
	return nil
}

func Save(storeX *Store) error {
	_, err := storeCol.UpdateOne(nil, bson.M{"_id": storeX.ID}, bson.M{"$set": storeX}, options.Update().SetUpsert(true))
	return err
}

func Delete(storeID int64) error {
	_, err := storeCol.DeleteOne(nil, bson.M{"_id": storeID})
	return err
}

func FileExists(bucketName string, storeID int64, songID primitive.ObjectID) error {
	storesMtx.RLock()
	mongoClient := storeConns[storeID]
	storesMtx.RUnlock()
	if mongoClient == nil {
		return errors.New("no store exists")
	}
	bucket, err := gridfs.NewBucket(
		mongoClient.Database(config.DbStore),
		func(bytes []byte) []byte {
			return bytes
		},
		func(bytes []byte) []byte {
			return bytes
		},
		options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}

	err = bucket.Exists(songID)
	bucket.Release()
	return err
}

func FileDelete(bucketName string, storeID int64, songID primitive.ObjectID) error {
	storesMtx.RLock()
	mongoClient := storeConns[storeID]
	storesMtx.RUnlock()
	if mongoClient == nil {
		return errors.New("no store exists")
	}
	bucket, err := gridfs.NewBucket(
		mongoClient.Database(config.DbStore),
		func(bytes []byte) []byte {
			return bytes
		},
		func(bytes []byte) []byte {
			return bytes
		},
		options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}

	err = bucket.Delete(songID)
	bucket.Release()
	return err
}

func GetUploadStream(bucketName string, songID primitive.ObjectID) (int64, *gridfs.UploadStream, error) {
	var mongoClient *mongo.Client
	var storeID int64
	storesMtx.RLock()
	for sID, x := range storeConns {
		if x != nil {
			mongoClient = x
			storeID = sID
			break
		}
	}
	storesMtx.RUnlock()
	if mongoClient == nil {
		return storeID, nil, errors.New("no store exists")
	}
	bucket, err := gridfs.NewBucket(
		mongoClient.Database(config.DbStore),
		func(bytes []byte) []byte {
			return bytes
		},
		func(bytes []byte) []byte {
			return bytes
		},
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return storeID, nil, err
	}

	_ = bucket.Delete(songID)
	stream, err := bucket.OpenUploadStreamWithID(songID, songID.Hex())
	bucket.Release()
	return storeID, stream, err
}

func GetDownloadStream(bucketName string, songID primitive.ObjectID, storeID int64) (*gridfs.DownloadStream, error) {
	storesMtx.RLock()
	mongoClient := storeConns[storeID]
	storesMtx.RUnlock()
	if mongoClient == nil {
		return nil, errors.New("no connection exists")
	}
	bucket, err := gridfs.NewBucket(
		mongoClient.Database(config.DbStore),
		func(bytes []byte) []byte {
			return bytes
		},
		func(bytes []byte) []byte {
			return bytes
		},
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return nil, err
	}
	stream, err := bucket.OpenDownloadStream(songID)
	bucket.Release()
	return stream, err
}

func Copy(dst io.Writer, src io.Reader, flushFunc func()) (written int64, err error) {
	size := 32 * 1024
	if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf := pbytes.GetLen(size)
	defer pbytes.Put(buf)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				if flushFunc != nil {
					flushFunc()
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// ForEachSong iterates through all the songs, this is for INTERNAL use only.
func ForEachSong(bucketName string, storeID int64, f func(songID primitive.ObjectID) bool) error {
	var lastID primitive.ObjectID
	storesMtx.RLock()
	mongoClient := storeConns[storeID]
	storesMtx.RUnlock()
	if mongoClient == nil {
		return errors.New("no store exists")
	}

	bucket, err := gridfs.NewBucket(
		mongoClient.Database(config.DbStore),
		func(bytes []byte) []byte {
			return bytes
		},
		func(bytes []byte) []byte {
			return bytes
		},
		options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return err
	}
	defer bucket.Release()

	cur, err := bucket.Find(bson.D{}, options.GridFSFind().SetNoCursorTimeout(true))
	if err != nil {
		return err
	}
	waitGroup := sync.WaitGroup{}
	defer waitGroup.Wait()
	rateLimit := make(chan struct{}, 20)
	for {
		for cur.Next(nil) {
			lastID = cur.Current.Lookup("_id").ObjectID()
			waitGroup.Add(1)
			rateLimit <- struct{}{}
			go func(songID primitive.ObjectID) {
				f(songID)
				waitGroup.Done()
				<-rateLimit
			}(lastID)
		}
		if cur.Err() == nil {
			_ = cur.Close(nil)
			break
		}
		log.Warn("Error On Cursor", zap.Error(err))
		_ = cur.Close(nil)
		cur, err = bucket.Find(
			bson.M{"_id": bson.M{"$gt": lastID}},
			options.GridFSFind().SetNoCursorTimeout(true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
