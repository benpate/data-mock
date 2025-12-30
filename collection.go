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

// Query retrieves multiple records from the mock collection.
func (collection Collection) Query(_ any, _ exp.Expression, _ ...option.Option) error {
	return derp.Internal("data-mock.collection.Query", "Unimplemented")
}

// Iterator retrieves a group of records as an Iterator.
func (collection Collection) Iterator(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {

	result := []data.Object{}

	if !collection.server.hasCollection(collection.name) {
		return NewIterator(result), derp.NotFound("mockdb.Load", "Collection does not exist", collection)
	}

	c := collection.server.getCollection(collection.name)

	for _, document := range c {
		if (criteria == nil) || (criteria.Match(MatcherFunc(document))) {
			result = append(result, document)
		}
	}

	iterator := NewIterator(result, options...)

	sort.Sort(iterator)

	return iterator, nil

}

// Load retrieves a single record from the mock collection.
func (collection Collection) Load(criteria exp.Expression, target data.Object, _ ...option.Option) error {

	if !collection.server.hasCollection(collection.name) {
		return derp.NotFound("mockdb.Load", "Collection does not exist", collection)
	}

	c := collection.server.getCollection(collection.name)

	for _, document := range c {

		if (criteria == nil) || (criteria.Match(MatcherFunc(document))) {
			populateInterface(document, target)
			return nil
		}
	}

	return derp.NotFound("mockdb.Load", "Document not found", criteria)
}

// Save adds/inserts a new record into the mock database
func (collection Collection) Save(object data.Object, comment string) error {

	const location = "mockdb.Save"

	// NILCHECK: Server cannot be nil
	if collection.server == nil {
		return derp.Internal(location, "Nil Server", "Attempted to save to a nil server", object)
	}

	// RULE: Handle synthetic errors (for testing purposes)
	if strings.HasPrefix(comment, "ERROR") {
		return derp.Internal(location, "Synthetic Error", comment)
	}

	c := collection.server.getCollection(collection.name)

	object.SetUpdated(comment)

	if object.IsNew() {
		object.SetCreated(comment)
		collection.setObjects(append(c, object))
		return nil
	}

	if index := collection.findByObjectID(object.ID()); index >= 0 {
		c[index] = object
		collection.setObjects(c)
		return nil
	}

	return derp.Internal(location, "Object Not Found", "attempted to update object, but it does not exist in the datastore", object)
}

// Delete soft deletes removes a record from the mock database.
func (collection Collection) Delete(object data.Object, comment string) error {

	if strings.HasPrefix(comment, "ERROR") {
		return derp.Internal("mockdb.Delete", "Synthetic Error", comment)
	}

	if index := collection.findByObjectID(object.ID()); index >= 0 {
		c := collection.server.getCollection(collection.name)
		collection.setObjects(append(c[:index], c[index+1:]...))
	}

	return nil
}

// HardDelete PERMANENTLY removes records from the mock database that match the criteria.
func (collection Collection) HardDelete(criteria exp.Expression) error {
	return derp.NotImplemented("data-mock.connection.HardDelete", "Not implemented", criteria)
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
