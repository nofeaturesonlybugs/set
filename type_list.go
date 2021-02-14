package set

import (
	"reflect"
)

// TypeList is a list of reflect.Type.
type TypeList map[reflect.Type]struct{}

// NewTypeList creates a new TypeList type from a set of instantiated structs.
func NewTypeList(args ...interface{}) TypeList {
	rv := make(TypeList)
	for _, arg := range args {
		T := reflect.TypeOf(arg)
		rv[T] = struct{}{}
	}
	return rv
}
