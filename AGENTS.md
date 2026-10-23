# data-mock — Notes for AI Agents

This repo is Go package `mockdb`, the in-memory stand-in for the `benpate/data` interfaces described in [README.md](README.md). Its whole job is to impersonate `data-mongo` in unit tests, so most of these notes are about where the impersonation is deliberately imperfect — the exact places a test can pass on the mock and fail on Mongo, or the reverse.

- **`Delete` is PHYSICAL here, virtual in data-mongo.** The mock removes the record from the slice and never touches the journal; data-mongo's `Delete` stamps `DeleteDate` and re-saves, leaving the row in place. A test that passes because "the record is gone" proves nothing about whether production queries filter on `deleteDate`. Deleting an absent record is a silent no-op, so repeated deletes stay idempotent.

- **Predicates resolve fields by EXACT bson tag; a miss means no match, ever.** `MatcherFunc` walks dotted bson-tag paths case-sensitively ([helpers.go](helpers.go)), and a field the struct does not declare makes every predicate FALSE — including `!=`, which MongoDB's `$ne` would match on a missing field. Comparison errors are also swallowed as "no match", so malformed criteria narrow a query, never widen it — the same fail-closed stance as data-mongo.

- **`IN`, `NOT IN`, `IN ALL`, `EXISTS`, and the GEO operators never match.** Matching delegates to `rosetta/compare.WithOperator`, which implements only the scalar and substring operators; anything else returns an error that [matcherFunc.go](matcherFunc.go) discards as FALSE. data-mongo implements all of them, so an expression built on these operators returns rows on Mongo and nothing at all here.

- **Sorting works; row-limiting does not.** `Iterator` honors `SortOption`s with a stable sort — ties keep insertion order, so keep `sort.Stable` — but `MaxRows`, `FirstRow`, and every other option are ignored, and the iterator always yields every match. Sorting also resolves its field by the OPPOSITE conventions from matching: case-insensitive struct field name or bson tag, no dotted paths ([iterator.go](iterator.go)).

- **`Load` ignores options entirely and returns the FIRST match in insertion order.** "Load the newest by sorting descending" works on data-mongo and silently returns the earliest-inserted match here.

- **Reads never create a collection.** `Server.getCollection` must not materialize the map key: doing so would let a `Count()` on a missing collection downgrade later `Load`/`Iterator` errors from "Collection does not exist" to "Document not found". Only `setObjects` creates the key. Note the read asymmetry that remains: `Count` on a missing collection returns 0 with no error, while `Load` and `Iterator` return NotFound.

- **`Save` stores the caller's object, not a copy.** The live `data.Object` pointer goes into the map, so mutating a struct after saving it retroactively rewrites the "database" — data-mongo serialized it to BSON at save time. Reads are insulated (`populateInterface` copies field-by-field on the way out), but the stored record itself aliases the caller's struct.

- **`Save` enforces no unique index.** A record whose `IsNew()` is true is appended even when another record already has the same `ID()`, where data-mongo rejects the duplicate key with a Conflict. Updating a non-new object whose ID is absent returns an error rather than upserting.

- **A journal comment starting with `ERROR` triggers a synthetic failure.** `Save` and `Delete` return an Internal error when the comment has that prefix — the mock's fault-injection hook. Don't pass a real comment beginning with "ERROR" unless you want the failure.

- **Only `Save` guards against a zero-value Collection.** `Count`, `Load`, `Iterator`, and `Delete` all panic on `mockdb.Collection{}` because they dereference the nil `server`. Build collections through `Server.Session(ctx).Collection(name)`. Relatedly, `populateInterface` panics (inside `reflect`) when a target is not a pointer to a struct — deliberately loud, because a silent no-op would hide caller errors in a test library.

- **No transactions, no locking, no `Query`/`HardDelete`.** `WithTransaction` runs the callback in a plain session and does NOT roll back on error; the `Server` map has no mutex, so concurrent goroutines race; `Query` and `HardDelete` always return NotImplemented even though data-mongo implements both.
