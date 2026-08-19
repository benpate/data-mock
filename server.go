package mockdb

import (
	"context"

	"github.com/benpate/data"
)

// Server is a mock database
type Server map[string][]data.Object

// New returns a fully initialized Server.
func New() *Server {

	// The pointer is the store's identity: every Session and Collection built from
	// this Server shares it, so writes made through one are visible to all the others.
	return &Server{}
}

// Session returns a session that can be used as a mock database.
func (server *Server) Session(ctx context.Context) (data.Session, error) {

	return Session{
		server:  server,
		context: ctx,
	}, nil
}

// WithTransaction executes a callback function within a new Session.
func (server *Server) WithTransaction(ctx context.Context, fn data.TransactionCallbackFunc) (any, error) {

	// The mock has no transaction semantics -- writes made by the callback are
	// applied directly to the Server and are NOT rolled back when it returns an error.
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

// getCollection returns the named collection, or an empty slice if it does not exist.
func (server *Server) getCollection(collection string) []data.Object {

	if result, exists := (*server)[collection]; exists {
		return result
	}

	// RULE: A missing collection MUST NOT be created here. Reads (Count, Load,
	// Iterator) all pass through this method, and materializing the key would let a
	// read change what later reads report -- "Collection does not exist" would
	// silently become "Document not found". Writers create the key via setObjects.
	return []data.Object{}
}
