package coerce

import (
	"fmt"
	"reflect"
	"strconv"
)

// parseUint attempts to parse the string first as an uint, then as int, and finally
// as float.  If all three fail then ErrInvalid is returned.
func parseUint(s string) (interface{}, error) {
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return u, nil
	} else if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	} else if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("%w; could not parse %v", ErrInvalid, s)
}

// Uint coerces v to uint.
func Uint(v interface{}) (uint, error) {
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
			if IntOverflowsUint(int64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case int8:
			if IntOverflowsUint(int64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case int16:
			if IntOverflowsUint(int64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case int32:
			if IntOverflowsUint(int64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case int64:
			if IntOverflowsUint(sw, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case uint:
			return sw, nil
		case uint8:
			return uint(sw), nil
		case uint16:
			return uint(sw), nil
		case uint32:
			return uint(sw), nil
		case uint64:
			if UintOverflowsUint(uint64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case float32:
			if FloatOverflowsUint(float64(sw), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case float64:
			if FloatOverflowsUint(sw, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, sw)
			}
			return uint(sw), nil
		case string:
			if v, err = parseUint(sw); err != nil {
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
			if FloatOverflowsUint(float64(f), strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, f)
			}
			return uint(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsUint(f, strconv.IntSize) {
				return 0, fmt.Errorf("%w; %v overflows uint", ErrOverflow, f)
			}
			return uint(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to uint", ErrUnsupported, v)
	}
}

// Uint8 coerces v to uint8.
func Uint8(v interface{}) (uint8, error) {
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
			if IntOverflowsUint(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case int8:
			if IntOverflowsUint(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case int16:
			if IntOverflowsUint(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case int32:
			if IntOverflowsUint(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case int64:
			if IntOverflowsUint(int64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case uint:
			if UintOverflowsUint(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case uint8:
			return sw, nil
		case uint16:
			if UintOverflowsUint(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case uint32:
			if UintOverflowsUint(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case uint64:
			if UintOverflowsUint(uint64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case float32:
			if FloatOverflowsUint(float64(sw), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case float64:
			if FloatOverflowsUint(sw, 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, sw)
			}
			return uint8(sw), nil
		case string:
			if v, err = parseUint(sw); err != nil {
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
			if FloatOverflowsUint(float64(f), 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, f)
			}
			return uint8(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsUint(f, 8) {
				return 0, fmt.Errorf("%w; %v overflows uint8", ErrOverflow, f)
			}
			return uint8(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to uint8", ErrUnsupported, v)
	}
}

// Uint16 coerces v to uint16.
func Uint16(v interface{}) (uint16, error) {
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
			if IntOverflowsUint(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case int8:
			if IntOverflowsUint(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case int16:
			if IntOverflowsUint(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case int32:
			if IntOverflowsUint(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case int64:
			if IntOverflowsUint(int64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case uint:
			if UintOverflowsUint(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case uint8:
			return uint16(sw), nil
		case uint16:
			return sw, nil
		case uint32:
			if UintOverflowsUint(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case uint64:
			if UintOverflowsUint(uint64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case float32:
			if FloatOverflowsUint(float64(sw), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case float64:
			if FloatOverflowsUint(sw, 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, sw)
			}
			return uint16(sw), nil
		case string:
			if v, err = parseUint(sw); err != nil {
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
			if FloatOverflowsUint(float64(f), 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, f)
			}
			return uint16(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsUint(f, 16) {
				return 0, fmt.Errorf("%w; %v overflows uint16", ErrOverflow, f)
			}
			return uint16(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to uint16", ErrUnsupported, v)
	}
}

// Uint32 coerces v to uint32.
func Uint32(v interface{}) (uint32, error) {
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
			if IntOverflowsUint(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case int8:
			if IntOverflowsUint(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case int16:
			if IntOverflowsUint(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case int32:
			if IntOverflowsUint(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case int64:
			if IntOverflowsUint(int64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case uint:
			if UintOverflowsUint(uint64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case uint8:
			return uint32(sw), nil
		case uint16:
			return uint32(sw), nil
		case uint32:
			return sw, nil
		case uint64:
			if UintOverflowsUint(uint64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case float32:
			if FloatOverflowsUint(float64(sw), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case float64:
			if FloatOverflowsUint(sw, 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, sw)
			}
			return uint32(sw), nil
		case string:
			if v, err = parseUint(sw); err != nil {
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
			if FloatOverflowsUint(float64(f), 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, f)
			}
			return uint32(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsUint(f, 32) {
				return 0, fmt.Errorf("%w; %v overflows uint32", ErrOverflow, f)
			}
			return uint32(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to uint32", ErrUnsupported, v)
	}
}

// Uint64 coerces v to uint64.
func Uint64(v interface{}) (uint64, error) {
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
			if IntOverflowsUint(int64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case int8:
			if IntOverflowsUint(int64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case int16:
			if IntOverflowsUint(int64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case int32:
			if IntOverflowsUint(int64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case int64:
			if IntOverflowsUint(int64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case uint:
			return uint64(sw), nil
		case uint8:
			return uint64(sw), nil
		case uint16:
			return uint64(sw), nil
		case uint32:
			return uint64(sw), nil
		case uint64:
			return uint64(sw), nil
		case float32:
			if FloatOverflowsUint(float64(sw), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case float64:
			if FloatOverflowsUint(sw, 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, sw)
			}
			return uint64(sw), nil
		case string:
			if v, err = parseUint(sw); err != nil {
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
			if FloatOverflowsUint(float64(f), 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, f)
			}
			return uint64(f), nil
		case reflect.Float64:
			f := KindFloat64(v)
			if FloatOverflowsUint(f, 64) {
				return 0, fmt.Errorf("%w; %v overflows uint64", ErrOverflow, f)
			}
			return uint64(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to uint64", ErrUnsupported, v)
	}
}
