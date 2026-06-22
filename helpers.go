package mockdb

import (
	"reflect"
	"strings"
)

// findField walks a struct along a dotted bson-tag path and returns the matching field's type, value, and TRUE if found.
func findField(structure reflect.Type, value reflect.Value, bson string) (reflect.Type, reflect.Value, bool) {

	first, rest := split(bson)

	// Search every field in the structure
	for index := 0; index < structure.NumField(); index = index + 1 {

		field := structure.Field(index)

		// If the field has a bson tag...
		if tag, ok := field.Tag.Lookup("bson"); ok {

			// If the bson tag matches the predicate field
			if first == tag {

				if rest == "" {
					return structure.Field(index).Type, value.Field(index), true
				}

				return findField(structure.Field(index).Type, value.Field(index), rest)
			}
		}
	}

	return structure, value, false
}

// split divides a dotted path into its first segment and the remainder (empty when there is no dot).
func split(input string) (string, string) {
	index := strings.Index(input, ".")

	if index == -1 {
		return input, ""
	}

	first := input[:index]
	rest := input[index+1:]

	return first, rest
}
