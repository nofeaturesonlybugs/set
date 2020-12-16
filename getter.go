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

// MapGetter accepts a map and returns a Getter.  Map keys need to be either interface{}
// or string; i.e. the map needs to be of type map[string]* or map[interface{}]*.
func MapGetter(m interface{}) Getter {
	rv := GetterFunc(func(key string) interface{} { return nil })
	//
	v := reflect.ValueOf(m)
	k, t := v.Kind(), v.Type()
	if k != reflect.Map {
		return rv
	}
	if t.Key().Kind() != reflect.String && t.Key().Kind() != reflect.Interface {
		return rv
	}
	//
	rv = GetterFunc(func(key string) interface{} {
		if reflected := v.MapIndex(reflect.ValueOf(key)); reflected.IsValid() {
			value := V(reflected.Interface())
			if value.IsMap {
				return MapGetter(reflected.Interface())
			} else if value.IsSlice && value.Elem != nil && value.Elem.IsMap {
				// A slice of maps; convert to a []Getter
				getterSlice := []Getter{}
				for k, max := 0, value.pv.Len(); k < max; k++ {
					getterSlice = append(getterSlice, MapGetter(value.pv.Index(k).Interface()))
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
