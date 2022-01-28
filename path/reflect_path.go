package path

import (
	"reflect"
)

// ReflectPath contains the bare minimum information to safely traverse
// a path from an origin to the value described by Index+Last.
type ReflectPath struct {
	HasPointer bool
	Index      []int
	Last       int
}

// Value accepts an originating struct value and traverses Index+Last to
// reach a leaf field.
//
// v should be a struct whose type is equal to the type used when creating the original
// Path from which this ReflectPath was derived.
//
// For performance critical code consider manually inlining this implementation.
func (p ReflectPath) Value(v reflect.Value) reflect.Value {
	if p.HasPointer {
		for _, n := range p.Index {
			v = v.Field(n)
			for ; v.Kind() == reflect.Ptr; v = v.Elem() {
				if v.IsNil() && v.CanSet() {
					v.Set(reflect.New(v.Type().Elem()))
				}
			}
		}
	} else {
		for _, n := range p.Index {
			v = v.Field(n)
		}
	}
	return v.Field(p.Last)
}
