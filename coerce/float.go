package coerce

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Float32 coerces v to float32.
func Float32(v interface{}) (float32, error) {
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
			return float32(sw), nil
		case int8:
			return float32(sw), nil
		case int16:
			return float32(sw), nil
		case int32:
			return float32(sw), nil
		case int64:
			return float32(sw), nil
		case uint:
			return float32(sw), nil
		case uint8:
			return float32(sw), nil
		case uint16:
			return float32(sw), nil
		case uint32:
			return float32(sw), nil
		case uint64:
			return float32(sw), nil
		case float32:
			return float32(sw), nil
		case float64:
			if sw > math.MaxFloat32 || sw < math.SmallestNonzeroFloat32 {
				return 0, fmt.Errorf("%w; %v overflows float32", ErrOverflow, sw)
			}
			return float32(sw), nil
		case string:
			f, err := strconv.ParseFloat(sw, 32)
			if err == nil {
				return float32(f), nil

			} else if errors.Is(err, strconv.ErrRange) {
				return 0, fmt.Errorf("%w; %v overflows float32", ErrOverflow, sw)
			} else if b, berr := strconv.ParseBool(sw); berr == nil {
				if b {
					return 1, nil
				}
				return 0, nil
			}
			return 0, fmt.Errorf("%w; %v", ErrInvalid, err.Error())
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
			return float32(reflect.ValueOf(v).Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float32(reflect.ValueOf(v).Uint()), nil
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			return KindFloat32(v), nil
		case reflect.Float64:
			f := reflect.ValueOf(v).Float()
			if f > math.MaxFloat32 || f < math.SmallestNonzeroFloat32 {
				return 0, fmt.Errorf("%w; %v overflows float32", ErrOverflow, f)
			}
			return float32(f), nil

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
		return 0, fmt.Errorf("%w; coerce %v to float32", ErrUnsupported, v)
	}
}

// Float64 coerces v to float64.
func Float64(v interface{}) (float64, error) {
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
			return float64(sw), nil
		case int8:
			return float64(sw), nil
		case int16:
			return float64(sw), nil
		case int32:
			return float64(sw), nil
		case int64:
			return float64(sw), nil
		case uint:
			return float64(sw), nil
		case uint8:
			return float64(sw), nil
		case uint16:
			return float64(sw), nil
		case uint32:
			return float64(sw), nil
		case uint64:
			return float64(sw), nil
		case float32:
			return float64(sw), nil
		case float64:
			return float64(sw), nil
		case string:
			f, err := strconv.ParseFloat(sw, 64)
			if err == nil {
				return float64(f), nil
			} else if errors.Is(err, strconv.ErrRange) {
				return 0, fmt.Errorf("%w; %v overflows float64", ErrOverflow, sw)
			} else if b, berr := strconv.ParseBool(sw); berr == nil {
				if b {
					return 1, nil
				}
				return 0, nil
			}
			return 0, fmt.Errorf("%w; %v", ErrInvalid, err.Error())
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
			return float64(reflect.ValueOf(v).Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(reflect.ValueOf(v).Uint()), nil
		case reflect.String:
			v = reflect.ValueOf(v).Convert(TypeString).Interface()
			continue
		case reflect.Float32:
			return float64(KindFloat32(v)), nil
		case reflect.Float64:
			return KindFloat64(v), nil

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
		return 0, fmt.Errorf("%w; coerce %v to float64", ErrUnsupported, v)
	}
}
