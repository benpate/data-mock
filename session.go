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

func (session Session) Context() context.Context {
	return session.context
}

// Close cleans up any remaining data created by the mock session.
func (session Session) Close() {

}
