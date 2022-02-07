package coerce

import (
	"fmt"
	"reflect"
	"strconv"
)

// parseInt attempts to parse the string first as an int, then as uint, and finally
// as float.  If all three fail then ErrInvalid is returned.
func parseInt(s string) (interface{}, error) {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	} else if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return u, nil
	} else if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("%w; could not parse %v", ErrInvalid, s)
}

// Int coerces v to int.
func Int(v interface{}) (int, error) {
	var err error
	for {
		switch sw := v.(type) {
		case nil:
			return 0, nil
		case bool:
			if sw {
				return 1, nil
			}
			return 0, nil
		case int:
			return sw, nil
		case int8:
			return int(sw), nil
		case int16:
			return int(sw), nil
		case int32:
			return int(sw), nil
		case int64:
			if IntOverflowsInt(sw, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case uint:
			if UintOverflowsInt(uint64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case uint8:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), strconv.IntSize) {
			// 	return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			// }
			return int(sw), nil
		case uint16:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), strconv.IntSize) {
			// 	return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			// }
			return int(sw), nil
		case uint32:
			if UintOverflowsInt(uint64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case uint64:
			if UintOverflowsInt(sw, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case float32:
			if FloatOverflowsInt(float64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case float64:
			if FloatOverflowsInt(sw, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, sw)
			}
			return int(sw), nil
		case string:
			if v, err = parseInt(sw); err != nil {
				return 0, err
			}
			continue
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
			if reflect.ValueOf(v).Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Int()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Uint()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			f := KindFloat32(v)
			if FloatOverflowsInt(float64(f), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, f)
			}
			return int(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsInt(f, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows int", ErrOverflow, f)
			}
			return int(f), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return 0, nil
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
			return 0, nil
		}
		//
		return 0, fmt.Errorf("%w; coerce %v to int", ErrUnsupported, v)
	}
}

// Int8 coerces v to int8.
func Int8(v interface{}) (int8, error) {
	var err error
	for {
		switch sw := v.(type) {
		case nil:
			return 0, nil
		case bool:
			if sw {
				return 1, nil
			}
			return 0, nil
		case int:
			if IntOverflowsInt(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case int8:
			return sw, nil
		case int16:
			if IntOverflowsInt(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case int32:
			if IntOverflowsInt(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case int64:
			if IntOverflowsInt(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case uint:
			if UintOverflowsInt(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case uint8:
			if UintOverflowsInt(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case uint16:
			if UintOverflowsInt(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case uint32:
			if UintOverflowsInt(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case uint64:
			if UintOverflowsInt(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case float32:
			if FloatOverflowsInt(float64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case float64:
			if FloatOverflowsInt(sw, 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, sw)
			}
			return int8(sw), nil
		case string:
			if v, err = parseInt(sw); err != nil {
				return 0, err
			}
			continue
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
			if reflect.ValueOf(v).Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Int()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Uint()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			f := KindFloat32(v)
			if FloatOverflowsInt(float64(f), 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, f)
			}
			return int8(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsInt(f, 8) {
				return 0, fmt.Errorf("%w; %v overflows int8", ErrOverflow, f)
			}
			return int8(f), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return 0, nil
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
			return 0, nil
		}
		//
		return 0, fmt.Errorf("%w; coerce %v to int8", ErrUnsupported, v)
	}
}

// Int16 coerces v to int16.
func Int16(v interface{}) (int16, error) {
	var err error
	for {
		switch sw := v.(type) {
		case nil:
			return 0, nil
		case bool:
			if sw {
				return 1, nil
			}
			return 0, nil
		case int:
			if IntOverflowsInt(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case int8:
			return int16(sw), nil
		case int16:
			return sw, nil
		case int32:
			if IntOverflowsInt(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case int64:
			if IntOverflowsInt(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case uint:
			if UintOverflowsInt(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case uint8:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), 16) {
			// 	return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			// }
			return int16(sw), nil
		case uint16:
			if UintOverflowsInt(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case uint32:
			if UintOverflowsInt(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case uint64:
			if UintOverflowsInt(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case float32:
			if FloatOverflowsInt(float64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case float64:
			if FloatOverflowsInt(sw, 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, sw)
			}
			return int16(sw), nil
		case string:
			if v, err = parseInt(sw); err != nil {
				return 0, err
			}
			continue
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
			if reflect.ValueOf(v).Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Int()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Uint()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			f := KindFloat32(v)
			if FloatOverflowsInt(float64(f), 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, f)
			}
			return int16(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsInt(f, 16) {
				return 0, fmt.Errorf("%w; %v overflows int16", ErrOverflow, f)
			}
			return int16(f), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return 0, nil
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
			return 0, nil
		}
		//
		return 0, fmt.Errorf("%w; coerce %v to int16", ErrUnsupported, v)
	}
}

// Int32 coerces v to int32.
func Int32(v interface{}) (int32, error) {
	var err error
	for {
		switch sw := v.(type) {
		case nil:
			return 0, nil
		case bool:
			if sw {
				return 1, nil
			}
			return 0, nil
		case int:
			if IntOverflowsInt(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case int8:
			return int32(sw), nil
		case int16:
			return int32(sw), nil
		case int32:
			return sw, nil
		case int64:
			if IntOverflowsInt(sw, 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case uint:
			if UintOverflowsInt(uint64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case uint8:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), 32) {
			// 	return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			// }
			return int32(sw), nil
		case uint16:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), 32) {
			// 	return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			// }
			return int32(sw), nil
		case uint32:
			if UintOverflowsInt(uint64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case uint64:
			if UintOverflowsInt(sw, 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case float32:
			if FloatOverflowsInt(float64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case float64:
			if FloatOverflowsInt(sw, 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, sw)
			}
			return int32(sw), nil
		case string:
			if v, err = parseInt(sw); err != nil {
				return 0, err
			}
			continue
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
			if reflect.ValueOf(v).Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Int()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Uint()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			f := KindFloat32(v)
			if FloatOverflowsInt(float64(f), 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, f)
			}
			return int32(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsInt(f, 32) {
				return 0, fmt.Errorf("%w; %v overflows int32", ErrOverflow, f)
			}
			return int32(f), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return 0, nil
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
			return 0, nil
		}
		//
		return 0, fmt.Errorf("%w; coerce %v to int32", ErrUnsupported, v)
	}
}

// Int64 coerces v to int64.
func Int64(v interface{}) (int64, error) {
	var err error
	for {
		switch sw := v.(type) {
		case nil:
			return 0, nil
		case bool:
			if sw {
				return 1, nil
			}
			return 0, nil
		case int:
			return int64(sw), nil
		case int8:
			return int64(sw), nil
		case int16:
			return int64(sw), nil
		case int32:
			return int64(sw), nil
		case int64:
			return sw, nil
		case uint:
			if UintOverflowsInt(uint64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, sw)
			}
			return int64(sw), nil
		case uint8:
			return int64(sw), nil
		case uint16:
			return int64(sw), nil
		case uint32:
			// NB Impossible overflow
			// if UintOverflowsInt(uint64(sw), 64) {
			// 	return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, sw)
			// }
			return int64(sw), nil
		case uint64:
			if UintOverflowsInt(sw, 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, sw)
			}
			return int64(sw), nil
		case float32:
			if FloatOverflowsInt(float64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, sw)
			}
			return int64(sw), nil
		case float64:
			if FloatOverflowsInt(sw, 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, sw)
			}
			return int64(sw), nil
		case string:
			if v, err = parseInt(sw); err != nil {
				return 0, err
			}
			continue
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
			if reflect.ValueOf(v).Bool() {
				return 1, nil
			}
			return 0, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = reflect.ValueOf(v).Int()
			continue
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = reflect.ValueOf(v).Uint()
			continue
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			f := KindFloat32(v)
			if FloatOverflowsInt(float64(f), 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, f)
			}
			return int64(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsInt(f, 64) {
				return 0, fmt.Errorf("%w; %v overflows int64", ErrOverflow, f)
			}
			return int64(f), nil

		case reflect.Ptr:
			rv := reflect.ValueOf(v)
			for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
				if rv.IsNil() {
					return 0, nil
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
			return 0, nil
		}
		//
		return 0, fmt.Errorf("%w; coerce %v to int64", ErrUnsupported, v)
	}
}
