package tyrgin

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

// GetMongoSession returns a mgo session. Uses MONGO_URI env variable.
func GetMongoSession() (*mgo.Session, error) {
	session, err := mgo.Dial(os.Getenv("MONGO_URI"))

	return session, err
}

// GetMongoDB takes a string and returns a mongo db of that name.
func GetMongoDB(d string) (*mgo.Database, error) {
	session, err := GetMongoSession()
	if err != nil {
		return nil, err
	}

	return session.DB(d), nil
}

// GetMongoCollection takes a string and returns a mongo collection of that name.
func GetMongoCollection(c string, db *mgo.Database) *mgo.Collection {
	return db.C(c)
}

// SafeGetMongoCollection takes a string and returns a mongo collection of that name
// only if it exists and returns an error if it does not.
func SafeGetMongoCollection(c string, db *mgo.Database) (*mgo.Collection, error) {
	cnames, err := db.CollectionNames()
	if err != nil {
		return nil, err
	}

	for _, name := range cnames {
		if name == c {
			collection := db.C(c)
			return collection, nil
		}
	}

	return nil, ErrorMongoCollectionFailure
}

// GetMongoCollectionCreate takes a string and returns a mongo collection of that name.
// Will create collection if it does not exist
func GetMongoCollectionCreate(c string, db *mgo.Database) (*mgo.Collection, error) {
	cnames, err := db.CollectionNames()
	if err != nil {
		return nil, err
	}

	collection := db.C(c)

	for _, name := range cnames {
		if name == c {
			return collection, nil
		}
	}

	collection.Insert()

	return nil, ErrorMongoCollectionFailure
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
