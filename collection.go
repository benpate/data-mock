package mockdb

import (
	"context"
	"sort"
	"strings"

	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
)

// Collection is a mock database collection
type Collection struct {
	server  *Server
	context context.Context
	name    string
}

// Context returns the context for this collection
func (collection Collection) Context() context.Context {
	return collection.context
}

// Count returns the number of records in the mock collection that match the criteria.
func (collection Collection) Count(criteria exp.Expression, _ ...option.Option) (int64, error) {

	var count int64

	for _, document := range collection.server.getCollection(collection.name) {
		if criteria.Match(MatcherFunc(document)) {
			count++
		}
	}

	return count, nil
}

// Query is not implemented by the mock collection, and always returns an error.
func (collection Collection) Query(_ any, _ exp.Expression, _ ...option.Option) error {
	return derp.NotImplemented("data-mock.collection.Query", "Not implemented")
}

// Iterator retrieves a group of records as an Iterator.
func (collection Collection) Iterator(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {

	result := []data.Object{}

	// RULE: The collection must already exist. An empty Iterator accompanies the
	// error so that callers who ignore it still get something safe to range over.
	if !collection.server.hasCollection(collection.name) {
		return NewIterator(result), derp.NotFound("data-mock.collection.Iterator", "Collection does not exist", collection)
	}

	// Collect every document that matches the criteria. A nil criteria matches all.
	for _, document := range collection.server.getCollection(collection.name) {
		if (criteria == nil) || (criteria.Match(MatcherFunc(document))) {
			result = append(result, document)
		}
	}

	iterator := NewIterator(result, options...)

	// RULE: Sort MUST be stable. Records that tie on every sort option keep their
	// insertion order, so tests that sort on a non-unique field get a repeatable
	// result instead of one that shifts with the sort algorithm.
	sort.Stable(iterator)

	return iterator, nil

}

// Load retrieves a single record from the mock collection.
func (collection Collection) Load(criteria exp.Expression, target data.Object, _ ...option.Option) error {

	if !collection.server.hasCollection(collection.name) {
		return derp.NotFound("data-mock.collection.Load", "Collection does not exist", collection)
	}

	c := collection.server.getCollection(collection.name)

	// Copy the first matching document into the caller's target
	for _, document := range c {

		if (criteria == nil) || (criteria.Match(MatcherFunc(document))) {
			populateInterface(document, target)
			return nil // Station
		}
	}

	return derp.NotFound("data-mock.collection.Load", "Document not found", criteria)
}

// Save inserts a new record into the mock database, or updates an existing one.
func (collection Collection) Save(object data.Object, comment string) error {

	const location = "data-mock.collection.Save"

	// NILCHECK: Server cannot be nil
	if collection.server == nil {
		return derp.Internal(location, "Nil Server", "Attempted to save to a nil server", object)
	}

	// RULE: Handle synthetic errors (for testing purposes)
	if strings.HasPrefix(comment, "ERROR") {
		return derp.Internal(location, "Synthetic Error", comment)
	}

	// Load the current contents of the collection
	c := collection.server.getCollection(collection.name)

	// Stamp the journal before the record lands in the datastore
	object.SetUpdated(comment)

	// Insert brand new records at the end of the collection
	if object.IsNew() {
		object.SetCreated(comment)
		collection.setObjects(append(c, object))
		return nil // Success (maybe?)
	}

	// Otherwise, overwrite the existing record in place
	if index := collection.findByObjectID(object.ID()); index >= 0 {
		c[index] = object
		collection.setObjects(c)
		return nil // I am Groot
	}

	return derp.Internal(location, "Object Not Found", "attempted to update object, but it does not exist in the datastore", object)
}

// Delete removes a record from the mock database.
func (collection Collection) Delete(object data.Object, comment string) error {

	// RULE: Handle synthetic errors (for testing purposes)
	if strings.HasPrefix(comment, "ERROR") {
		return derp.Internal("data-mock.collection.Delete", "Synthetic Error", comment)
	}

	// Drop the record from the collection. Deleting a record that is not present
	// is a no-op, so that repeated deletes stay idempotent.
	if index := collection.findByObjectID(object.ID()); index >= 0 {
		c := collection.server.getCollection(collection.name)
		collection.setObjects(append(c[:index], c[index+1:]...))
	}

	return nil // Silence is golden
}

// HardDelete is not implemented by the mock collection, and always returns an error.
func (collection Collection) HardDelete(criteria exp.Expression) error {
	return derp.NotImplemented("data-mock.collection.HardDelete", "Not implemented", criteria)
}

// getObjects retrieves the slice of objects for this collection from the server
func (collection Collection) getObjects() []data.Object {
	return (*collection.server)[collection.name]
}

// setObjects sets the slice of objects for this collection on the server
func (collection Collection) setObjects(objects []data.Object) {
	(*collection.server)[collection.name] = objects
}

// findByObjectID does a linear search on the collection for the first object with a matching ID()
func (collection Collection) findByObjectID(objectID string) int {

	objects := collection.getObjects()

	for index, object := range objects {

		if object.ID() == objectID {
			return index
		}
	}

	return -1
}
