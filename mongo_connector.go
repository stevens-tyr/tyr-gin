package tyrgin

import (
	"errors"
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

// GetSession returns a mgo session. Uses MANGO_URI env variable.
func GetSession() (*mgo.Session, error) {
	session, err := mgo.Dial(os.Getenv("MANGO_URI"))

	return session, err
}

// GetDataBase takes a string and returns a mongo db of that name.
func GetDataBase(d string, s *mgo.Session) *mgo.Database {
	return s.DB(d)
}

// GetCollection takes a string and returns a mongo collection of that name.
func GetCollection(c string, db *mgo.Database) *mgo.Collection {
	return db.C(c)
}

// SafeGetCollection takes a string and returns a mongo collection of that name
// only if it exists and returns an error if it does not.
func SafeGetCollection(c string, db *mgo.Database) (*mgo.Collection, error) {
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

	err = errors.New("Collection does not exist")

	return nil, err
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
