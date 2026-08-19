// Package mockdb provides an in-memory implementation of the benpate/data
// interfaces, for use in unit tests.
//
// A Server holds every record in a plain map of named collections, so tests can
// build a fixture as a struct literal and hand it to code that expects a real
// datastore. Sessions and Collections created from a Server all share it, which
// means writes made through one are immediately visible to the others.
//
// Records are matched and sorted by reflection over their bson struct tags, so
// this package is SLOW by design and belongs only in tests. Query and HardDelete
// are not implemented, and WithTransaction provides no rollback.
package mockdb
