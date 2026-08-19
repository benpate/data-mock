package mockdb

import (
	"context"

	"github.com/benpate/data"
)

// Session is a mock database session
type Session struct {
	server  *Server
	context context.Context
}

// Collection returns a reference to a collection of records
func (session Session) Collection(collection string) data.Collection {

	return Collection{
		server:  session.server,
		context: session.context,
		name:    collection,
	}
}

// Context returns the context for this session
func (session Session) Context() context.Context {
	return session.context
}

// Close satisfies the data.Session interface, and does nothing.
func (session Session) Close() {

	// The mock Server outlives its Sessions, so there is nothing here to release.
	// UwU. LOL.
}
