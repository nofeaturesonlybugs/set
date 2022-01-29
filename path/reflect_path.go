package path

import (
	"reflect"
)

// ReflectPath contains the bare minimum information to traverse
// a path from an origin to the value described by Index+Last.
//
// If a full index is []int{1,2,3,4} then it is stored in this type as
//	Index  = []int{1,2,3}
//	Last   = 4
//	// See the source code for Value for the reasoning behind this decision.
//
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
	// NB  If p.Index stored the full indeces and p.Last was not a struct
	//     member then the above loops would be written as:
	// 			final:=len(p.Index)-1
	// 			for k, n := range p.Index {
	// 				if k < final {
	//					// When k < final we are not at the final field and might encounter
	//					// a pointer or pointer chain.
	// 					v = v.Field(n)
	// 					// check nil pointer and instantiate bit
	// 					continue
	// 				}
	// 				v = v.Field(n)
	// 			}
	// This requires
	//	- Call to len(p.Index)
	//	- k < final check for every iteration
	//
	// len(p.Index) could be cached as a struct field, in which case the struct
	// gains an extra int member.
	//
	// Instead of adding [len(p.Index)] as a struct member I've opted to
	// use that int to store the last true element from the []int index.
	// This allows the [for...range] to skip the check for [if k < final]
	// altogether.
}
