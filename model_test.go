package mockdb

import (
	"github.com/benpate/data/journal"
	"github.com/benpate/derp"
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

func (person testPerson) GetPath(p string) (any, bool) {
	return nil, false
}

func (person *testPerson) SetPath(p string, value any) error {
	return derp.NewInternalError("data-mock.testPerson", "Unsupported GetPath", p)
}
