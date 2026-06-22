package mockdb

import (
	"context"

	"github.com/benpate/data"
)

// Server is a mock database
type Server map[string][]data.Object

// New returns a fully initialized Server. The pointer is the store's identity:
// every Session and Collection shares it, so writes are visible across them.
func New() *Server {
	return &Server{}
}

// Session returns a session that can be used as a mock database.
func (server *Server) Session(ctx context.Context) (data.Session, error) {

	return Session{
		server:  server,
		context: ctx,
	}, nil
}

// WithTransaction executes a callback function within the context of a transaction.
func (server *Server) WithTransaction(ctx context.Context, fn data.TransactionCallbackFunc) (any, error) {
	session, err := server.Session(ctx)

	if err != nil {
		return nil, err
	}

	return fn(session)
}

// hasCollection returns TRUE if the designated collection already exists in the Server
func (server *Server) hasCollection(collection string) bool {

	_, ok := (*server)[collection]

	return ok
}

// getCollection loads (and creates, if necessary) the named collection in this datastore
func (server *Server) getCollection(collection string) []data.Object {

	if result, exists := (*server)[collection]; exists {
		return result
	}

	// Create and store an empty (non-nil) collection so callers always get a usable slice.
	result := []data.Object{}
	(*server)[collection] = result
	return result
}
