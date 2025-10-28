package mockdb

import (
	"context"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

// MODEL OBJECT

func TestDataType(t *testing.T) {
	var server data.Server
	var session data.Session
	var collection data.Collection
	var err error

	server = New()
	session, err = server.Session(context.TODO())
	require.Nil(t, err)

	collection = session.Collection("Person")

	require.NotNil(t, server)
	require.NotNil(t, session)
	require.NotNil(t, collection)
}

func TestDelete(t *testing.T) {

	ds := getSampleDataset()

	session, err := ds.Session(context.TODO())
	require.NotNil(t, session)
	require.Nil(t, err)

	collection := session.Collection("Person").(Collection)

	require.Equal(t, 4, len(collection.getObjects()))

	// Delete FIRST entry
	require.Nil(t, collection.Delete(&testPerson{PersonID: "michael"}, ""))
	require.Equal(t, 3, len(collection.getObjects()))

	// Delete MIDDLE entry
	require.Nil(t, collection.Delete(&testPerson{PersonID: "latoya"}, ""))
	require.Equal(t, 2, len(collection.getObjects()))

	// Delete LAST entry
	require.Nil(t, collection.Delete(&testPerson{PersonID: "janet"}, ""))
	require.Equal(t, 1, len(collection.getObjects()))

	// Delete ONLY entry
	require.Nil(t, collection.Delete(&testPerson{PersonID: "jermaine"}, ""))
	require.Equal(t, 0, len(collection.getObjects()))
}

func TestSession(t *testing.T) {

	ds := New()

	session, err := ds.Session(context.TODO())
	require.NotNil(t, session)
	require.Nil(t, err)

	collection := session.Collection("Person")

	john := testPerson{
		PersonID: "A",
		Name:     "John Connor",
		Email:    "john@connor.com",
	}

	// CREATE
	{
		err := collection.Save(&john, "created in test suite")
		require.Nil(t, err)
	}

	// READ & UPDATE
	{
		// Load a record from the db
		person := testPerson{}
		criteria := exp.Equal("_id", "A")
		err := collection.Load(criteria, &person)
		require.Nil(t, err)
		require.Equal(t, "A", person.PersonID)
		require.Equal(t, "John Connor", person.Name)
		require.Equal(t, "john@connor.com", person.Email)

		// Change some values
		person.Name = "Sarah Connor"
		person.Email = "sarah@sky.net"

		// Save the record
		err = collection.Save(&person, "Comment Here")
		require.Nil(t, err)

		person2 := testPerson{}
		err = collection.Load(criteria, &person2)
		require.Nil(t, err)
		require.Equal(t, "Sarah Connor", person2.Name)
		require.Equal(t, "sarah@sky.net", person2.Email)
		require.Equal(t, "Comment Here", person2.Note)
	}
}

func TestList(t *testing.T) {

	ds := New()

	session, err := ds.Session(context.TODO())
	require.NotNil(t, session)
	require.Nil(t, err)

	collection := session.Collection("Person")

	require.Nil(t, collection.Save(&testPerson{PersonID: "A", Name: "Sarah Connor", Email: "sarah@sky.net"}, "Creating Record"))
	require.Nil(t, collection.Save(&testPerson{PersonID: "B", Name: "John Connor", Email: "john@connor.com"}, "Creating Record"))
	require.Nil(t, collection.Save(&testPerson{PersonID: "C", Name: "Kyle Reese", Email: "kyle@resistance.mil"}, "Creating Record"))

	criteria := exp.Equal("_id", "A")

	it, err := collection.Iterator(criteria)

	require.Nil(t, err)

	person := testPerson{}

	require.True(t, it.Next(&person))
	require.Equal(t, "A", person.PersonID)
	require.Equal(t, "Sarah Connor", person.Name)
	require.Equal(t, "sarah@sky.net", person.Email)

	require.False(t, it.Next(&person))
}

func TestErrors(t *testing.T) {

	ds := New()

	session, err := ds.Session(context.TODO())
	require.Nil(t, err)

	person := &testPerson{}

	{
		err := session.Collection("MissingCollection").Load(nil, person).(derp.Error)
		require.NotNil(t, err)
		require.Equal(t, 404, err.Code)
		require.Equal(t, "mockdb.Load", err.Location)
		require.Equal(t, "Collection does not exist", err.Message)
		// require.Equal(t, []any{"MissingCollection"}, err.Details)
	}

	{
		err := session.Collection("Person").Save(person, "ERROR: Testing error codes").(derp.Error)
		require.NotNil(t, err)
		require.Equal(t, 500, err.Code)
		require.Equal(t, "mockdb.Save", err.Location)
		require.Equal(t, "Synthetic Error", err.Message)
	}
}
