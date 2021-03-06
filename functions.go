package set

import (
	"reflect"
)

// Writable attempts to make a reflect.Value usable for writing.  It will follow and instantiate nil pointers if necessary.
func Writable(v reflect.Value) (V reflect.Value, CanWrite bool) {
	if !v.IsValid() {
		return
	}
	T := v.Type()
	K := T.Kind()
	// if T == nil || K == reflect.Invalid { 	// N.B: Couldn't write a test case that hit this for 100% coverage
	// 	return									// but I'm not entirely convinced I can take it out.
	// }
	V, T, K = v, v.Type(), v.Kind()
	for K == reflect.Ptr {
		if V.IsNil() && V.CanSet() {
			ptr := reflect.New(T.Elem())
			V.Set(ptr)
		}
		K = T.Elem().Kind()
		T = T.Elem()
		V = V.Elem()
	}
	CanWrite = V.CanSet()
	return
}
