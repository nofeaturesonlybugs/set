package set

import (
	"reflect"
)

// TypeList is a list of reflect.Type.
type TypeList map[reflect.Type]struct{}

// NewTypeList creates a new TypeList type from a set of instantiated types.
func NewTypeList(args ...interface{}) TypeList {
	rv := make(TypeList)
	for _, arg := range args {
		T := reflect.TypeOf(arg)
		rv[T] = struct{}{}
	}
	return rv
}

// Has returns true if the specified type is in the list.
func (list TypeList) Has(T reflect.Type) bool {
	if list == nil {
		return false
	}
	_, has := list[T]
	return has
}

// Merge adds entries in `from` to this list.
func (list TypeList) Merge(from TypeList) {
	if list != nil {
		for k, v := range from {
			list[k] = v
		}
	}
}
