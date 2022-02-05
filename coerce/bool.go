package coerce

import (
	"fmt"
	"reflect"
	"strconv"
)

// Bool coerces v to bool.
func Bool(v interface{}) (bool, error) {
	for {
		switch sw := v.(type) {
		case nil:
			return false, nil
		case bool:
			return sw, nil
		case int:
			return sw != 0, nil
		case int8:
			return sw != 0, nil
		case int16:
			return sw != 0, nil
		case int32:
			return sw != 0, nil
		case int64:
			return sw != 0, nil
		case uint:
			return sw != 0, nil
		case uint8:
			return sw != 0, nil
		case uint16:
			return sw != 0, nil
		case uint32:
			return sw != 0, nil
		case uint64:
			return sw != 0, nil
		case float32:
			return sw != 0, nil
		case float64:
			return sw != 0, nil
		case string:
			b, err := strconv.ParseBool(sw)
			if err != nil {
				return false, fmt.Errorf("%w; %v", ErrInvalid, err.Error())
			}
			return b, nil
		}
		//
		// Beyond this point we need reflection.
		T := reflect.TypeOf(v)
		//
		// - T.Kind() is a primitive
		//		convert to actual primitive and try again
		// - T.Kind() is a pointer
		//		dereference pointer and try again
		// - T.Kind() is a slice
		//		pick last element and try again
		switch T.Kind() {
		case reflect.Bool:
			return reflect.ValueOf(v).Convert(TypeBool).Interface().(bool), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Convert(TypeInt64).Interface()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Convert(TypeUint64).Interface()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32, reflect.Float64:
			v = reflect.ValueOf(v).Convert(TypeFloat64).Interface()
			continue

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return false, nil
				}
			}
			v = rv.Interface()
			continue

		case reflect.Slice:
			rv := reflect.ValueOf(v)
			if n := rv.Len(); n > 0 {
				v = rv.Index(n - 1).Interface()
				continue
			}
			return false, nil
		}
		//
		return false, fmt.Errorf("%w; coerce %v to bool", ErrUnsupported, v)
	}
}
