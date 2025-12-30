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

func (person testPerson) GetPath(_ string) (any, bool) {
	return nil, false
}

func (person *testPerson) SetPath(p string, _ any) error {
	return derp.Internal("data-mock.testPerson", "Unsupported GetPath", p)
}
