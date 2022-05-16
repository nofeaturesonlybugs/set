package set

import (
	"reflect"
)

// Field is a struct field; it contains a Value and a reflect.StructField.
type Field struct {
	Value    Value
	Field    reflect.StructField
	TagValue string
}
