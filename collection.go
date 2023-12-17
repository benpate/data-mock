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
	Server  *Server
	Context context.Context
	Name    string
}

func (collection Collection) Count(criteria exp.Expression, options ...option.Option) (int64, error) {

	var count int64

	for _, document := range collection.Server.getCollection(collection.Name) {
		if criteria.Match(MatcherFunc(document)) {
			count++
		}
	}

	return count, nil
}

func (collection Collection) Query(target any, criteria exp.Expression, options ...option.Option) error {
	return derp.NewInternalError("data-mock.collection.Query", "Unimplemented")
}

// List retrieves a group of records as an Iterator.
func (collection Collection) Iterator(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {

	result := []data.Object{}

	if !collection.Server.hasCollection(collection.Name) {
		return NewIterator(result), derp.NewNotFoundError("mockdb.Load", "Collection does not exist", collection)
	}

	c := collection.Server.getCollection(collection.Name)

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
func (collection Collection) Load(criteria exp.Expression, target data.Object) error {

	if !collection.Server.hasCollection(collection.Name) {
		return derp.NewNotFoundError("mockdb.Load", "Collection does not exist", collection)
	}

	c := collection.Server.getCollection(collection.Name)

	for _, document := range c {

		if (criteria == nil) || (criteria.Match(MatcherFunc(document))) {
			populateInterface(document, target)
			return nil
		}
	}

	return derp.NewNotFoundError("mockdb.Load", "Document not found", criteria)
}

// Save adds/inserts a new record into the mock database
func (collection Collection) Save(object data.Object, comment string) error {

	if strings.HasPrefix(comment, "ERROR") {
		return derp.NewInternalError("mockdb.Save", "Synthetic Error", comment)
	}

	c := collection.Server.getCollection(collection.Name)

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

	return derp.NewInternalError("mockdb.Save", "Object Not Found", "attempted to update object, but it does not exist in the datastore", object)
}

// Delete PERMANENTLY removes a record from the mock database.
func (collection Collection) Delete(object data.Object, comment string) error {

	if strings.HasPrefix(comment, "ERROR") {
		return derp.NewInternalError("mockdb.Delete", "Synthetic Error", comment)
	}

	c := collection.Server.getCollection(collection.Name)

	if index := collection.findByObjectID(object.ID()); index >= 0 {
		collection.setObjects(append(c[:index], c[index+1:]...))
	}

	return nil
}

func (collection Collection) HardDelete(criteria exp.Expression) error {
	return derp.NewBadRequestError("data-mock.connection.HardDelete", "Not implemented", criteria)
}

func (collection Collection) getObjects() []data.Object {
	return (*collection.Server)[collection.Name]
}

func (collection Collection) setObjects(objects []data.Object) {
	(*collection.Server)[collection.Name] = objects
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
