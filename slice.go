package set

import (
	"reflect"
)

// SliceValue wraps around a slice []T and facilitates element creation and appending.
type SliceValue struct {
	// Top is the original values passed to Slice.
	Top reflect.Value

	// V represents []T.
	V reflect.Value

	// ElemType and ElemEndType describe the type of elements in the slice.
	//
	// ElemType!=reflect.Ptr means ElemType equals ElemEndType; they describe the same type.
	// ElemType==reflect.Ptr means ElemEndType is the type at the end of the pointer chain.
	ElemType    reflect.Type
	ElemEndType reflect.Type
}

// Slice expects its argument to be a *[]T or a pointer chain ending in []T.
//
// As a convenience Slice will also accept a reflect.Value as long as it represents
// a writable []T value.
func Slice(v interface{}) (SliceValue, error) {
	var s SliceValue
	var pv reflect.Value
	//
	// Allow reflect.Value as an argument.
	switch sw := v.(type) {
	case reflect.Value:
		pv = sw
	default:
		pv = reflect.ValueOf(v)
	}
	//
	// Incoming type must be *[]T or any length ****[]T
	pt := pv.Type()
	if pt.Kind() != reflect.Ptr {
		return s, pkgerr{Err: ErrInvalidSlice, CallSite: "Slice", Context: "expected pointer to slice; got " + pt.String()}
	}
	for pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() != reflect.Slice {
		return s, pkgerr{Err: ErrInvalidSlice, CallSite: "Slice", Context: "expected pointer to slice; got " + pv.Type().String()}
	}
	//
	// Type is valid so instantiate pointer chain to final []T.
	top := pv
	for pv.Kind() == reflect.Ptr {
		// Update our pointer reference each iteration.  We want our ptr reference
		// to be the **last** pointer in the chain such that slicePtr.ptr.Elem() is the []T.
		// slicePtr.ptr = pv // TODO RM
		//
		if pv.IsNil() {
			if !pv.CanSet() {
				return s, pkgerr{Err: ErrReadOnly, CallSite: "Slice", Context: "can not set " + pv.Type().String()}
			}
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		pv = pv.Elem()
	}
	//
	elemType := pv.Type().Elem()
	endElemType := elemType
	for ; endElemType.Kind() == reflect.Ptr; endElemType = endElemType.Elem() {
	}
	s = SliceValue{
		Top:         top,
		V:           pv,
		ElemType:    elemType,
		ElemEndType: endElemType,
	}
	return s, nil
}

// Append appends an element created by the Elem method.
//
// If the slice is []T then Elem returns a *T.  Append automatically dereferences
// the incoming value so that a T is appended as expected.
func (s *SliceValue) Append(elem reflect.Value) {
	s.V.Set(reflect.Append(s.V, elem.Elem()))
}

// Elem returns a newly allocated slice element.
//
// If the slice is []T then Elem returns a *T.  This is so the created element
// can be passed directly into function wanting to populate T.
func (s SliceValue) Elem() reflect.Value {
	return reflect.New(s.ElemType)
}
