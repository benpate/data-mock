package mockdb

import (
	"context"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/stretchr/testify/assert"
)

// MODEL OBJECT

func TestDelete(t *testing.T) {

	ds := getSampleDataset()

	session, _ := ds.Session(context.TODO())
	collection := session.Collection("Person").(Collection)

	assert.Equal(t, 4, len(collection.getObjects()))

	// Delete FIRST entry
	assert.Nil(t, collection.Delete(&testPerson{PersonID: "michael"}, ""))
	assert.Equal(t, 3, len(collection.getObjects()))

	// Delete MIDDLE entry
	assert.Nil(t, collection.Delete(&testPerson{PersonID: "latoya"}, ""))
	assert.Equal(t, 2, len(collection.getObjects()))

	// Delete LAST entry
	assert.Nil(t, collection.Delete(&testPerson{PersonID: "janet"}, ""))
	assert.Equal(t, 1, len(collection.getObjects()))

	// Delete ONLY entry
	assert.Nil(t, collection.Delete(&testPerson{PersonID: "jermaine"}, ""))
	assert.Equal(t, 0, len(collection.getObjects()))
}

func TestSession(t *testing.T) {

	ds := New()

	session, _ := ds.Session(context.TODO())
	collection := session.Collection("Person")

	john := testPerson{
		PersonID: "A",
		Name:     "John Connor",
		Email:    "john@connor.com",
	}

	// CREATE
	{
		err := collection.Save(&john, "created in test suite")
		assert.Nil(t, err)
	}

	// READ & UPDATE
	{
		// Load a record from the db
		person := testPerson{}
		criteria := exp.Equal("_id", "A")
		err := collection.Load(criteria, &person)
		assert.Nil(t, err)
		assert.Equal(t, "A", person.PersonID)
		assert.Equal(t, "John Connor", person.Name)
		assert.Equal(t, "john@connor.com", person.Email)

		// Change some values
		person.Name = "Sarah Connor"
		person.Email = "sarah@sky.net"

		// Save the record
		err = collection.Save(&person, "Comment Here")
		assert.Nil(t, err)

		person2 := testPerson{}
		err = collection.Load(criteria, &person2)
		assert.Nil(t, err)
		assert.Equal(t, "Sarah Connor", person2.Name)
		assert.Equal(t, "sarah@sky.net", person2.Email)
		assert.Equal(t, "Comment Here", person2.Journal.Note)
	}
}

func TestList(t *testing.T) {

	ds := New()

	session, _ := ds.Session(context.TODO())
	collection := session.Collection("Person")

	collection.Save(&testPerson{PersonID: "A", Name: "Sarah Connor", Email: "sarah@sky.net"}, "Creating Record")
	collection.Save(&testPerson{PersonID: "B", Name: "John Connor", Email: "john@connor.com"}, "Creating Record")
	collection.Save(&testPerson{PersonID: "C", Name: "Kyle Reese", Email: "kyle@resistance.mil"}, "Creating Record")

	criteria := exp.Equal("_id", "A")

	it, err := collection.Iterator(criteria)

	assert.Nil(t, err)

	person := testPerson{}

	assert.True(t, it.Next(&person))
	assert.Equal(t, "A", person.PersonID)
	assert.Equal(t, "Sarah Connor", person.Name)
	assert.Equal(t, "sarah@sky.net", person.Email)

	assert.False(t, it.Next(&person))
}

func TestErrors(t *testing.T) {

	ds := New()

	session, _ := ds.Session(context.TODO())

	person := &testPerson{}

	{
		err := session.Collection("MissingCollection").Load(nil, person).(derp.Error)
		assert.NotNil(t, err)
		assert.Equal(t, derp.CodeNotFoundError, err.Code)
		assert.Equal(t, "mockdb.Load", err.Location)
		assert.Equal(t, "Collection does not exist", err.Message)
		// assert.Equal(t, []any{"MissingCollection"}, err.Details)
	}

	{
		err := session.Collection("Person").Save(person, "ERROR: Testing error codes").(derp.Error)
		assert.NotNil(t, err)
		assert.Equal(t, derp.CodeInternalError, err.Code)
		assert.Equal(t, "mockdb.Save", err.Location)
		assert.Equal(t, "Synthetic Error", err.Message)
	}
}
