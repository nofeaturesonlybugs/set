package set

import (
	"fmt"
	"reflect"
)

var (
	// ErrInvalidSlicePtr is returned by NewSlicePtr when the in coming value is not pointer-to-slice.
	ErrInvalidSlicePtr = fmt.Errorf("set: expected pointer-to-slice *[]T")
)

// SlicePtr wraps around a *[]T and can be used by deserializers expecting
// their argument to be passed as a pointer to a slice.
//	func Deserializer(dst interface{}) error {
// 		// dst is expected to be pointer to slice
// 		ptr, err := NewSlicePtr(dst)
// 		if err != nil { // was not *[]T
// 			return err
// 		}
// 		// ... continue deserialization of data into dst
// 	}
type SlicePtr struct {
	Slice
	ptr reflect.Value
}

// NewSlicePtr validates the incoming value is a *[]T and returns an error
// if it is not.
func NewSlicePtr(v interface{}) (SlicePtr, error) {
	var slicePtr SlicePtr
	var pv reflect.Value
	//
	// Allow reflect.Value as an argument type.
	switch sw := v.(type) {
	case reflect.Value:
		pv = sw
	default:
		pv = reflect.ValueOf(v)
	}
	//
	// Validate incoming type is *[]T or any amount of ****[]T
	pt := pv.Type()
	if pt.Kind() != reflect.Ptr {
		return slicePtr, fmt.Errorf("%w; got %T", ErrInvalidSlicePtr, v)
	}
	for pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() != reflect.Slice {
		return slicePtr, fmt.Errorf("%w; got %T", ErrInvalidSlicePtr, v)
	}
	//
	// Type is valid so instantiate pointer chain to final []T.
	for pv.Kind() == reflect.Ptr {
		// Update our pointer reference each iteration.  We want our ptr reference
		// to be the **last** pointer in the chain such that slicePtr.ptr.Elem() is the []T.
		slicePtr.ptr = pv
		//
		if pv.IsNil() {
			if !pv.CanSet() {
				return slicePtr, fmt.Errorf("%w; %T %v", ErrReadOnly, v, v)
			}
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		pv = pv.Elem()
	}
	slicePtr.V = pv
	slicePtr.Type = pv.Type().Elem()
	for slicePtr.EndType = slicePtr.Type; slicePtr.EndType.Kind() == reflect.Ptr; slicePtr.EndType = slicePtr.EndType.Elem() {
	}
	return slicePtr, nil
}

// Commit updates the original pointer to the contain the slice that was built.
func (ptr *SlicePtr) Commit() {
	ptr.ptr.Elem().Set(ptr.V)
}

// Slice wraps around []T and facilitates element creation and appending.
type Slice struct {
	// V represents []T.
	V reflect.Value

	// Type and EndType describe the type of elements in the slice.
	//
	// Type!=reflect.Ptr means Type equals EndType; they describe the same type.
	// Type==reflect.Ptr means EndType is the type at the end of the pointer chain.
	Type    reflect.Type
	EndType reflect.Type
}

// Append appends an element created by the Elem method.
//
// If the slice is []T then Elem returns a *T.  Append automatically dereferences
// the incoming value so that a T is appended as expected.
func (s *Slice) Append(elem reflect.Value) {
	s.V = reflect.Append(s.V, reflect.Indirect(elem))
}

// Elem returns a newly allocated slice element.
//
// If the slice is []T then Elem returns a *T.  This is so the created element
// can be passed directly into function wanting to populate T.
func (s Slice) Elem() reflect.Value {
	return reflect.New(s.Type)
}
