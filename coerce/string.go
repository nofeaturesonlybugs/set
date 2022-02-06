package coerce

import (
	"fmt"
	"reflect"
	"strconv"
)

// String coerces v to string.
func String(v interface{}) (string, error) {
	for {
		switch sw := v.(type) {
		case nil:
			return "", nil
		case bool:
			return strconv.FormatBool(sw), nil
		case int:
			return strconv.FormatInt(int64(sw), 10), nil
		case int8:
			return strconv.FormatInt(int64(sw), 10), nil
		case int16:
			return strconv.FormatInt(int64(sw), 10), nil
		case int32:
			return strconv.FormatInt(int64(sw), 10), nil
		case int64:
			return strconv.FormatInt(sw, 10), nil
		case uint:
			return strconv.FormatUint(uint64(sw), 10), nil
		case uint8:
			return strconv.FormatUint(uint64(sw), 10), nil
		case uint16:
			return strconv.FormatUint(uint64(sw), 10), nil
		case uint32:
			return strconv.FormatUint(uint64(sw), 10), nil
		case uint64:
			return strconv.FormatUint(sw, 10), nil
		case float32:
			return strconv.FormatFloat(float64(sw), 'g', -1, 32), nil
		case float64:
			return strconv.FormatFloat(sw, 'g', -1, 64), nil
		case string:
			return sw, nil
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
			return strconv.FormatBool(reflect.ValueOf(v).Bool()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(reflect.ValueOf(v).Int(), 10), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10), nil
		case reflect.String:
			return reflect.ValueOf(v).Convert(TypeString).Interface().(string), nil
		case reflect.Float32:
			f := KindFloat32(v)
			return strconv.FormatFloat(float64(f), 'g', -1, 32), nil
		case reflect.Float64:
			f := KindFloat64(v)
			return strconv.FormatFloat(f, 'g', -1, 64), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return "", nil
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
			return "", nil
		}
		//
		return "", fmt.Errorf("%w; coerce %v to string", ErrUnsupported, v)
	}
}
