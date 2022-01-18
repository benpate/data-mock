package mockdb

import (
	"github.com/benpate/data/journal"
	"github.com/benpate/path"
)

type testPerson struct {
	PersonID        string `bson:"_id"`
	Name            string `bson:"name"`
	Email           string `bson:"email"`
	Age             int    `bson:"age"`
	journal.Journal `bson:"journal"`
}

func (person testPerson) ID() string {
	return person.PersonID
}

func (person testPerson) GetPath(p path.Path) (interface{}, error) {
	return nil, nil
}

func (person *testPerson) SetPath(p path.Path, value interface{}) error {
	return nil
}
