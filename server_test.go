package mockdb

import (
	"context"
	"testing"

	"github.com/benpate/data"
	"github.com/benpate/exp"
	"github.com/stretchr/testify/require"
)

func TestServer_WithTransaction(t *testing.T) {

	server := getSampleDataset()

	result, err := server.WithTransaction(context.TODO(), func(session data.Session) (any, error) {

		collection := session.Collection("Person").(Collection)
		return len(collection.getObjects()), nil
	})

	require.Nil(t, err)
	require.Equal(t, 4, result)
}

func TestServer_WithTransaction_PropagatesError(t *testing.T) {

	server := getSampleDataset()

	_, err := server.WithTransaction(context.TODO(), func(_ data.Session) (any, error) {
		return nil, errForTest
	})

	require.Equal(t, errForTest, err)
}

func TestServer_HasCollection(t *testing.T) {

	server := getSampleDataset()
	require.True(t, server.hasCollection("Person"))
	require.False(t, server.hasCollection("Missing"))
}

func TestServer_GetCollection_ReturnsEmptyWithoutCreating(t *testing.T) {

	server := New()

	// Getting a non-existent collection yields a usable empty slice...
	require.False(t, server.hasCollection("New"))
	result := server.getCollection("New")
	require.NotNil(t, result)
	require.Equal(t, 0, len(result))

	// ...but MUST NOT materialize the collection in the Server
	require.False(t, server.hasCollection("New"))
}

// TestServer_ReadsDoNotCreateCollection guards the regression where Count()
// materialized a missing collection, downgrading every later Load() on that
// name from "Collection does not exist" to "Document not found".
func TestServer_ReadsDoNotCreateCollection(t *testing.T) {

	server := New()
	session, err := server.Session(context.TODO())
	require.Nil(t, err)

	collection := session.Collection("Missing")

	count, err := collection.Count(exp.All())
	require.Nil(t, err)
	require.Equal(t, int64(0), count)

	require.False(t, server.hasCollection("Missing"))

	err = collection.Load(exp.All(), &testPerson{})
	require.ErrorContains(t, err, "Collection does not exist")

	_, err = collection.Iterator(exp.All())
	require.ErrorContains(t, err, "Collection does not exist")
}

var errForTest = errTest("synthetic test error")

type errTest string

func (e errTest) Error() string { return string(e) }
