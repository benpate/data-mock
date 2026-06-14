package mockdb

import (
	"context"
	"testing"

	"github.com/benpate/data"
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

	_, err := server.WithTransaction(context.TODO(), func(session data.Session) (any, error) {
		return nil, errForTest
	})

	require.Equal(t, errForTest, err)
}

func TestServer_HasCollection(t *testing.T) {

	server := getSampleDataset()
	require.True(t, server.hasCollection("Person"))
	require.False(t, server.hasCollection("Missing"))
}

func TestServer_GetCollection_CreatesEmpty(t *testing.T) {

	server := New().(Server)

	// Getting a non-existent collection creates an empty one
	require.False(t, server.hasCollection("New"))
	result := server.getCollection("New")
	require.NotNil(t, result)
	require.Equal(t, 0, len(result))
	require.True(t, server.hasCollection("New"))
}

var errForTest = errTest("synthetic test error")

type errTest string

func (e errTest) Error() string { return string(e) }
