package mockdb

import (
	"context"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

func TestCollection_Context(t *testing.T) {

	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("trace"), "abc")

	session, err := getSampleDataset().Session(ctx)
	require.Nil(t, err)

	collection := session.Collection("Person").(Collection)
	require.Equal(t, ctx, collection.Context())
	require.Equal(t, "abc", collection.Context().Value(ctxKey("trace")))
}

func TestCollection_Count(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person")

	// Count all records (matching everything)
	count, err := collection.Count(exp.All())
	require.Nil(t, err)
	require.Equal(t, int64(4), count)

	// Count records matching a specific criteria
	count, err = collection.Count(exp.Equal("_id", "michael"))
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	// Count records matching nothing
	count, err = collection.Count(exp.Equal("_id", "missing"))
	require.Nil(t, err)
	require.Equal(t, int64(0), count)
}

func TestCollection_Query_Unimplemented(t *testing.T) {

	session, err := New().Session(context.TODO())
	require.Nil(t, err)

	err = session.Collection("Person").Query(nil, exp.All())
	require.NotNil(t, err)
}

func TestCollection_HardDelete_NotImplemented(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person").(Collection)
	err = collection.HardDelete(exp.All())
	require.NotNil(t, err)
}

func TestCollection_Iterator_MissingCollection(t *testing.T) {

	session, err := New().Session(context.TODO())
	require.Nil(t, err)

	it, err := session.Collection("Missing").Iterator(exp.All())
	require.NotNil(t, err)
	require.NotNil(t, it)            // returns an empty iterator, not nil
	require.Equal(t, 0, it.Count())
}

func TestCollection_Iterator_NilCriteria(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	// A nil criteria matches every document
	it, err := session.Collection("Person").Iterator(nil)
	require.Nil(t, err)
	require.Equal(t, 4, it.Count())
}

func TestCollection_Load_MissingDocument(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	person := testPerson{}
	err = session.Collection("Person").Load(exp.Equal("_id", "missing"), &person)
	require.NotNil(t, err)
	require.Equal(t, 404, err.(derp.Error).Code)
}

func TestCollection_Save_NilServer(t *testing.T) {

	// A Collection with a nil server cannot save
	collection := Collection{name: "Person"}
	err := collection.Save(&testPerson{PersonID: "x"}, "comment")
	require.NotNil(t, err)
}

func TestCollection_Save_UpdateMissingObject(t *testing.T) {

	session, err := New().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person")

	// An object that is not new (has a journal) but does not exist -> error
	person := &testPerson{PersonID: "ghost"}
	person.SetCreated("seed") // mark as not-new

	err = collection.Save(person, "update")
	require.NotNil(t, err)
}

func TestCollection_Delete_SyntheticError(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person")
	err = collection.Delete(&testPerson{PersonID: "michael"}, "ERROR: synthetic")
	require.NotNil(t, err)
}

func TestCollection_Delete_MissingObject(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person").(Collection)

	// Deleting an object that does not exist is a no-op (no error, no change)
	require.Equal(t, 4, len(collection.getObjects()))
	require.Nil(t, collection.Delete(&testPerson{PersonID: "missing"}, ""))
	require.Equal(t, 4, len(collection.getObjects()))
}

func TestCollection_FindByObjectID(t *testing.T) {

	session, err := getSampleDataset().Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Person").(Collection)

	require.Equal(t, 0, collection.findByObjectID("michael"))
	require.Equal(t, 3, collection.findByObjectID("janet"))
	require.Equal(t, -1, collection.findByObjectID("missing"))
}
