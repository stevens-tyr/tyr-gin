package tyrgin

import (
	"bytes"
	ctx "context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// GetMongoSession returns a mgo session. Uses MONGO_URI env variable.
func GetMongoSession() (*mongo.Client, error) {
	session, err := mongo.NewClient("mongodb://" + os.Getenv("MONGO_URI"))
	if err != nil {
		return nil, err
	}

	err = session.Connect(ctx.TODO())

	return session, err
}

// GetMongoDB takes a string and returns a mongo db of that name.
func GetMongoDB(d string) (*mongo.Database, error) {
	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}

	return session.Database(d), nil
}

// GetMongoCollection takes a string and returns a mongo collection of that name.
func GetMongoCollection(c string, db *mongo.Database) *mongo.Collection {
	return db.Collection(c)
}

// SafeGetMongoCollection takes a string and returns a mongo collection of that name
// only if it exists and returns an error if it does not.
func SafeGetMongoCollection(c string, db *mongo.Database) (*mongo.Collection, error) {
	cnames, err := db.ListCollections(ctx.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer cnames.Close(ctx.Background())

	for cnames.Next(ctx.Background()) {
		var name string
		err = cnames.Decode(name)
		if err != nil {
			return nil, err
		}

		if name == c {
			collection := db.Collection(c)
			return collection, nil
		}
	}

	return nil, ErrorMongoCollectionFailure
}

// GetMongoCollectionCreate takes a string and returns a mongo collection of that name.
// Will create collection if it does not exist
func GetMongoCollectionCreate(c string, db *mongo.Database) (*mongo.Collection, error) {
	cnames, err := db.ListCollections(ctx.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer cnames.Close(ctx.Background())

	collection := db.Collection(c)

	for cnames.Next(ctx.Background()) {
		var name string
		err = cnames.Decode(name)
		if err != nil {
			return nil, err
		}

		if name == c {
			return collection, nil
		}
	}

	collection.InsertOne(ctx.Background(), nil)

	return collection, nil
}

// GetGridFSBucket returns a mongo gridfs bucket given a  name and chunk size in bytes for a bucket.
func GetGridFSBucket(db *mongo.Database, name string, size int32) (*Bucket, error) {
	bucketOptions := options.GridFSBucket()

	bucketOptions.Name = &name
	bucketOptions.ChunkSizeBytes = &size

	fsbucket, err := gridfs.NewBucket(db, bucketOptions)
	if err != nil {
		return nil, err
	}

	bucket := &Bucket{
		Bucket:         fsbucket,
		ChunkSizeBytes: &size,
		Name:           &name,
	}

	return bucket, nil
}

// GridFSUploadFile  uploads a file to a Bucket given name of file and the data as a reader object.
func (b *Bucket) GridFSUploadFile(fileID primitive.ObjectID, filename string, file io.Reader) error {
	uploadStreamOptions := options.GridFSUpload()

	uploadStreamOptions.ChunkSizeBytes = b.ChunkSizeBytes

	err := b.Bucket.UploadFromStreamWithID(fileID, filename, file, uploadStreamOptions)
	if err != nil {
		return err
	}

	return nil
}

// GridFSDownloadFile given a fileID downloads the file from the bucket.
func (b *Bucket) GridFSDownloadFile(fileID primitive.ObjectID) (bytes.Buffer, error) {
	var fsStream bytes.Buffer

	_, err := b.Bucket.DownloadToStream(fileID, &fsStream)
	if err != nil {
		return fsStream, err
	}

	return fsStream, nil
}

// GridFSDeleteFile given a fileID deletes the file from the bucket.
func (b *Bucket) GridFSDeleteFile(fileID primitive.ObjectID) error {
	err := b.Bucket.Delete(fileID)

	return err
}

// CheckStatus here is of the struct for checking mongo replica set statuses.
func (m MongoRPLStatusChecker) CheckStatus(name string) StatusList {
	var replResult MongoReplStatus
	var cmd interface{}
	cmd = bson.D{{"replSetGetStatus", 1}}
	raw := m.RPL.RunCommand(ctx.Background(), cmd)
	raw.Decode(&replResult)

	var result Status
	if replResult.OK == 0 {
		result = Status{
			Description: name,
			Result:      CRITICAL,
			Details:     fmt.Sprintf("%v check failed: %v", name, replResult.ErrorMsg),
		}
	} else {
		result = Status{
			Description: name,
			Result:      OK,
			Details:     "",
		}
	}

	return StatusList{StatusList: []Status{result}}
}
