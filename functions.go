package set

import (
	"reflect"
)

// Writable attempts to make a reflect.Value usable for writing.  It will follow and instantiate nil pointers if necessary.
func Writable(v reflect.Value) (V reflect.Value, CanWrite bool) {
	if !v.IsValid() {
		return
	}
	for V = v; V.Kind() == reflect.Ptr; V = V.Elem() {
		if V.IsNil() && V.CanSet() {
			V.Set(reflect.New(V.Type().Elem()))
		}
	}
	CanWrite = V.CanSet()
	return
}
