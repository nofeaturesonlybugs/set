package set

import (
	"reflect"
)

// Writable attempts to make a reflect.Value usable for writing.  It will follow and instantiate nil pointers if necessary.
func Writable(v reflect.Value) (V reflect.Value, Info TypeInfo, CanWrite bool) {
	if !v.IsValid() {
		return
	} else if Info = TypeCache.StatType(v.Type()); Info.Type == nil || Info.Kind == reflect.Invalid {
		return
	}
	var K reflect.Kind
	var T reflect.Type
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
