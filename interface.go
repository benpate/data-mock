package mockdb

import (
	"reflect"
)

// populateInterface copies each same-named, settable field from the source struct into the target struct.
//
// Both arguments MUST be structs (or pointers to structs). Anything else panics
// inside the reflect package, which is loud on purpose: it means a test handed
// the mock a target it can never fill in.
func populateInterface(source any, target any) {

	// Reach through a pointer to get at the source struct itself
	sourceValue := reflect.ValueOf(source)

	if sourceValue.Kind() == reflect.Pointer {
		sourceValue = reflect.Indirect(sourceValue)
	}

	sourceType := sourceValue.Type()

	targetValue := reflect.Indirect(reflect.ValueOf(target))

	// Copy field-by-field, skipping any the target does not have (or cannot set,
	// such as unexported fields) so that partial targets are still usable.
	for index := 0; index < sourceType.NumField(); index = index + 1 {

		sourceField := sourceType.FieldByIndex([]int{index})

		if targetField := targetValue.FieldByName(sourceField.Name); targetField.CanSet() {
			targetField.Set(sourceValue.FieldByName(sourceField.Name))
		}
	}
}
