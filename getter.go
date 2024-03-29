package set

import (
	"reflect"
)

// Getter returns a value by name.
type Getter interface {
	// Get accepts a name and returns the value.
	Get(name string) interface{}
}

// GetterFunc casts a function into a Getter.
type GetterFunc func(name string) interface{}

// Get accepts a name and returns the value.
func (me GetterFunc) Get(name string) interface{} {
	return me(name)
}

// MapGetter accepts a map and returns a Getter.
//
// Map keys must be string or interface{}.
func MapGetter(m interface{}) Getter {
	rv := GetterFunc(func(key string) interface{} { return nil })
	//
	RV := reflect.ValueOf(m)
	K, T := RV.Kind(), RV.Type()
	if K != reflect.Map {
		return rv
	}
	if T.Key().Kind() != reflect.String && T.Key().Kind() != reflect.Interface {
		return rv
	}
	//
	rv = GetterFunc(func(key string) interface{} {
		if reflected := RV.MapIndex(reflect.ValueOf(key)); reflected.IsValid() {
			value := V(reflected.Interface())
			if value.IsMap {
				return MapGetter(reflected.Interface())
			} else if value.IsSlice && value.ElemTypeInfo.IsMap {
				getterSlice := []Getter{}
				for k, max := 0, value.WriteValue.Len(); k < max; k++ {
					getterSlice = append(getterSlice, MapGetter(value.WriteValue.Index(k).Interface()))
				}
				return getterSlice
			} else {
				return reflected.Interface()
			}
		} else {
			return nil
		}
	})
	//
	return rv
}
