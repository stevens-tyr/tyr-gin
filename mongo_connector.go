package tyrgin

import (
	"errors"
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession() (*mgo.Session, error) {
	session, err := mgo.Dial(os.Getenv("MANGO_URI"))

	return session, err
}

func GetDataBase(d string, s *mgo.Session) *mgo.Database {
	return s.DB(d)
}

func GetCollection(c string, db *mgo.Database) *mgo.Collection {
	return db.C(c)
}

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

func (m MongoDBStatusChecker) CheckStatus(name string) StatusList {
	var result MongoReplStatus
	m.RPL.Run(bson.D{{"replSetGetStatus", 1}}, &result)

	if result.OK == 0 {
		result = Status{
			Description: name,
			Result:      CRITICAL,
			Details:     fmt.Sprintf("%v check failed: %v", name, result.ErrorMsg),
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
