package tyrgin

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// GetMongoSession returns a mgo session. Uses MONGO_URI env variable.
func GetMongoSession() (*mongo.Client, error) {
	session, err := mongo.NewClient("mongodb://" + os.Getenv("MONGO_URI"))
	if err != nil {
		return nil, err
	}

	err = session.Connect(context.TODO())

	return session, err
}

// GetMongoDB takes a string and returns a mongo db of that name.
func GetMongoDB(d string) (*mongo.Database, error) {
	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}

	return session.DB(d), nil
}

// GetMongoCollection takes a string and returns a mongo collection of that name.
func GetMongoCollection(c string, db *mongo.Database) *mongo.Collection {
	return db.Collection(c)
}

// SafeGetMongoCollection takes a string and returns a mongo collection of that name
// only if it exists and returns an error if it does not.
func SafeGetMongoCollection(c string, db *mongo.Database) (*mongo.Collection, error) {
	cnames, err := db.ListCollections()
	if err != nil {
		return nil, err
	}
	defer cnames.Close(context.Backround())

	for cnames.Next(context.Background()) {
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
	cnames, err := db.ListCollections(context.Backround())
	if err != nil {
		return nil, err
	}
	defer cnames.Close(context.Backround())

	collection := db.Collection(c)

	for cnames.Next(context.Background()) {
		var name string
		err = cnames.Decode(name)
		if err != nil {
			return nil, err
		}

		if name == c {
			return collection, nil
		}
	}

	collection.Insert()

	return collection, nil
}

// CheckStatus here is of the struct for checking mongo replica set statuses.
func (m MongoRPLStatusChecker) CheckStatus(name string) StatusList {
	var replResult MongoReplStatus
	m.RPL.Run("replSetGetStatus", &replResult)

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
