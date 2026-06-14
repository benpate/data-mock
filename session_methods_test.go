package mockdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSession_Context(t *testing.T) {

	type ctxKey string
	ctx := context.WithValue(context.Background(), ctxKey("id"), 42)

	session, err := New().Session(ctx)
	require.Nil(t, err)

	require.Equal(t, ctx, session.Context())
	require.Equal(t, 42, session.Context().Value(ctxKey("id")))
}

func TestSession_Close(t *testing.T) {

	session, err := New().Session(context.TODO())
	require.Nil(t, err)

	// Close is a no-op but must not panic
	require.NotPanics(t, func() {
		session.Close()
	})
}
