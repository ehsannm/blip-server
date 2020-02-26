// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package gridfs

import (
	"context"
	"errors"
	"time"

	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// UploadBufferSize is the size in bytes of one stream batch. Chunks will be written to the db after the sum of chunk
// lengths is equal to the batch size.
const UploadBufferSize = 16 * 1024 * 1024 // 16 MiB

// ErrStreamClosed is an error returned if an operation is attempted on a closed/aborted stream.
var ErrStreamClosed = errors.New("stream is closed or aborted")

// UploadStream is used to upload files in chunks.
type UploadStream struct {
	*Upload // chunk size and metadata
	FileID  interface{}

	chunkIndex    int
	chunksColl    *mongo.Collection // collection to store file chunks
	filename      string
	filesColl     *mongo.Collection // collection to store file metadata
	closed        bool
	buffer        []byte
	bufferIndex   int
	fileLen       int64
	writeDeadline time.Time
	encryptFunc   EncryptFunc
}

// NewUploadStream creates a new upload stream.
func newUploadStream(upload *Upload, encryptFunc EncryptFunc, fileID interface{}, filename string, chunks, files *mongo.Collection) *UploadStream {
	if encryptFunc == nil {
		encryptFunc = func(bytes []byte) []byte {
			return bytes
		}
	}
	buf := largePool.GetLen(UploadBufferSize)
	return &UploadStream{
		Upload: upload,
		FileID: fileID,

		chunksColl:  chunks,
		filename:    filename,
		filesColl:   files,
		buffer:      buf,
		encryptFunc: encryptFunc,
	}
}

// Close closes this upload stream.
func (us *UploadStream) Close() error {
	if us.closed {
		return ErrStreamClosed
	}

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil {
		defer cancel()
	}

	if us.bufferIndex != 0 {
		if err := us.uploadChunks(ctx, true); err != nil {
			return err
		}
	}

	if err := us.createFilesCollDoc(ctx); err != nil {
		return err
	}

	us.closed = true
	largePool.Put(us.buffer)
	return nil
}

// SetWriteDeadline sets the write deadline for this stream.
func (us *UploadStream) SetWriteDeadline(t time.Time) error {
	if us.closed {
		return ErrStreamClosed
	}

	us.writeDeadline = t
	return nil
}

// Write transfers the contents of a byte slice into this upload stream. If the stream's underlying buffer fills up,
// the buffer will be uploaded as chunks to the server. Implements the io.Writer interface.
func (us *UploadStream) Write(p []byte) (int, error) {
	if us.closed {
		return 0, ErrStreamClosed
	}

	var ctx context.Context

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil {
		defer cancel()
	}

	origLen := len(p)
	for {
		if len(p) == 0 {
			break
		}

		n := copy(us.buffer[us.bufferIndex:], p) // copy as much as possible
		p = p[n:]
		us.bufferIndex += n

		if us.bufferIndex == UploadBufferSize {
			err := us.uploadChunks(ctx, false)
			if err != nil {
				return 0, err
			}
		}
	}
	return origLen, nil
}

// Abort closes the stream and deletes all file chunks that have already been written.
func (us *UploadStream) Abort() error {
	if us.closed {
		return ErrStreamClosed
	}

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil {
		defer cancel()
	}

	id, err := convertFileID(us.FileID)
	if err != nil {
		return err
	}
	_, err = us.chunksColl.DeleteMany(ctx, bsonx.Doc{{Key: "files_id", Value: id}})
	if err != nil {
		return err
	}

	us.closed = true
	return nil
}

// uploadChunks uploads the current buffer as a series of chunks to the bucket
// if uploadPartial is true, any data at the end of the buffer that is smaller than a chunk will be uploaded as a partial
// chunk. if it is false, the data will be moved to the front of the buffer.
// uploadChunks sets us.bufferIndex to the next available index in the buffer after uploading
func (us *UploadStream) uploadChunks(ctx context.Context, uploadPartial bool) error {
	chunks := float64(us.bufferIndex) / float64(us.chunkSize)
	numChunks := int(math.Ceil(chunks))
	if !uploadPartial {
		numChunks = int(math.Floor(chunks))
	}

	docs := make([]interface{}, int(numChunks))

	id, err := convertFileID(us.FileID)
	if err != nil {
		return err
	}
	begChunkIndex := us.chunkIndex
	for i := 0; i < us.bufferIndex; i += int(us.chunkSize) {
		endIndex := i + int(us.chunkSize)
		if us.bufferIndex-i < int(us.chunkSize) {
			// partial chunk
			if !uploadPartial {
				break
			}
			endIndex = us.bufferIndex
		}
		chunkData := us.encryptFunc(us.buffer[i:endIndex])
		docs[us.chunkIndex-begChunkIndex] = bsonx.Doc{
			{Key: "_id", Value: bsonx.ObjectID(primitive.NewObjectID())},
			{Key: "files_id", Value: id},
			{Key: "n", Value: bsonx.Int32(int32(us.chunkIndex))},
			{Key: "data", Value: bsonx.Binary(0x00, chunkData)},
		}
		us.chunkIndex++
		us.fileLen += int64(len(chunkData))
	}

	_, err = us.chunksColl.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	// copy any remaining bytes to beginning of buffer and set buffer index
	bytesUploaded := numChunks * int(us.chunkSize)
	if bytesUploaded != UploadBufferSize && !uploadPartial {
		copy(us.buffer[0:], us.buffer[bytesUploaded:us.bufferIndex])
	}
	us.bufferIndex = UploadBufferSize - bytesUploaded
	return nil
}

func (us *UploadStream) createFilesCollDoc(ctx context.Context) error {
	id, err := convertFileID(us.FileID)
	if err != nil {
		return err
	}
	doc := bsonx.Doc{
		{Key: "_id", Value: id},
		{Key: "length", Value: bsonx.Int64(us.fileLen)},
		{Key: "chunkSize", Value: bsonx.Int32(us.chunkSize)},
		{Key: "uploadDate", Value: bsonx.DateTime(time.Now().UnixNano() / int64(time.Millisecond))},
		{Key: "filename", Value: bsonx.String(us.filename)},
		{Key: "encrypted", Value: bsonx.Boolean(true)},
	}

	if us.metadata != nil {
		doc = append(doc, bsonx.Elem{
			Key:   "metadata",
			Value: bsonx.Document(us.metadata),
		})
	}

	_, err = us.filesColl.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}
