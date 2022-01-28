package path

import (
	"reflect"
)

// ReflectPath contains the bare minimum information to safely traverse
// a path from an origin to the value described by Index+Last.
//
// NB  If a full index is []int{1,2,3,4} then it is stored in this type as
//	Index  = []int{1,2,3}
//  Last   = 4
//
// When iterating Index we instantiate and follow any nil pointers except
// the reflect.Value for the last index:
//   lastK := len(Index)-1 // necessary if Index has **all** elements.
//   for k, n := range Index {
//       if k < lastK && HasPointer {
//           v = v.Field(n)
//           // possibly create pointer
//           continue
//       }
//       v = v.Field(n)
//   }
//
// By storing the last true index in Last we no longer need to know len(Index)
// or make a comparison during every iteration.
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
		// NB  Since the last true index is stored in Last we know every element
		//     during the range needs to be checked for instantiation without
		//     further checks.
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
