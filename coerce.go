package set

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/nofeaturesonlybugs/errors"
)

// coercions is a function map of type conversions.  Each entry is a function:
//	func( target, value ) error {
//		// The data in value is coerced into the type for target and assigned to target.
//	}
//
// The following provisions hold true for all functions in this map and should help to keep
// them lean.
//	+ It is taken as a given that these functions are called within a defer...recover construct
//		and can not crash the program.
//	+ It is taken as a given that these functions do not need to zero out target to a zero
//		value as that is done before they are called.
var coercions = map[string]func(reflect.Value, reflect.Value) error{
	"float-to-bool": func(target reflect.Value, value reflect.Value) error {
		if value.Float() != 0 {
			target.SetBool(true)
		} else {
			target.SetBool(false)
		}
		return nil
	},
	"int-to-bool": func(target reflect.Value, value reflect.Value) error {
		if value.Int() != 0 {
			target.SetBool(true)
		} else {
			target.SetBool(false)
		}
		return nil
	},
	"string-to-bool": func(target reflect.Value, value reflect.Value) error {
		var err error
		var parsed bool
		if parsed, err = strconv.ParseBool(value.String()); err != nil {
			return errors.Go(err)
		}
		target.SetBool(parsed)
		return nil
	},
	"uint-to-bool": func(target reflect.Value, value reflect.Value) error {
		if value.Uint() != 0 {
			target.SetBool(true)
		} else {
			target.SetBool(false)
		}
		return nil
	},

	"bool-to-float": func(target reflect.Value, value reflect.Value) error {
		if value.Bool() {
			target.SetFloat(float64(1))
		} else {
			target.SetFloat(float64(0))
		}
		return nil
	},
	"int-to-float": func(target reflect.Value, value reflect.Value) error {
		target.SetFloat(float64(value.Int()))
		return nil
	},
	"string-to-float": func(target reflect.Value, value reflect.Value) error {
		var err error
		var parsed float64
		if parsed, err = strconv.ParseFloat(value.String(), target.Type().Bits()); err != nil {
			return errors.Go(err)
		}
		target.SetFloat(parsed)
		return nil
	},
	"uint-to-float": func(target reflect.Value, value reflect.Value) error {
		target.SetFloat(float64(value.Uint()))
		return nil
	},

	"bool-to-int": func(target reflect.Value, value reflect.Value) error {
		if value.Bool() {
			target.SetInt(1)
		} else {
			target.SetInt(0)
		}
		return nil
	},
	"float-to-int": func(target reflect.Value, value reflect.Value) error {
		target.SetInt(int64(value.Float()))
		return nil
	},
	"string-to-int": func(target reflect.Value, value reflect.Value) error {
		if parsed, err := strconv.ParseInt(value.String(), 0, target.Type().Bits()); err == nil {
			target.SetInt(parsed)
		} else if parsedFloat, err := strconv.ParseFloat(value.String(), target.Type().Bits()); err == nil {
			target.SetInt(int64(parsedFloat))
		} else {
			return errors.Go(err)
		}
		return nil
	},
	"uint-to-int": func(target reflect.Value, value reflect.Value) error {
		target.SetInt(int64(value.Uint()))
		return nil
	},

	"bool-to-uint": func(target reflect.Value, value reflect.Value) error {
		if value.Bool() {
			target.SetUint(1)
		} else {
			target.SetUint(0)
		}
		return nil
	},
	"float-to-uint": func(target reflect.Value, value reflect.Value) error {
		if value.Float() < 0 {
			return errors.Errorf("Can not coerce negative float to uint.")
		}
		target.SetUint(uint64(value.Float()))
		return nil
	},
	"int-to-uint": func(target reflect.Value, value reflect.Value) error {
		if value.Int() < 0 {
			return errors.Errorf("Can not coerce negative int to uint.")
		}
		target.SetUint(uint64(value.Int()))
		return nil
	},
	"string-to-uint": func(target reflect.Value, value reflect.Value) error {
		var parsed uint64
		var parsedFloat float64
		var err error
		if len(value.String()) > 0 && rune(value.String()[0]) == '-' {
			return errors.Errorf("Can not coerce negative number to uint.")
		} else if parsed, err = strconv.ParseUint(value.String(), 0, target.Type().Bits()); err == nil {
			target.SetUint(parsed)
		} else if parsedFloat, err = strconv.ParseFloat(value.String(), target.Type().Bits()); err == nil {
			target.SetUint(uint64(parsedFloat))
		} else {
			return errors.Go(err)
		}
		return err
	},

	"bool-to-string": func(target reflect.Value, value reflect.Value) error {
		target.SetString(fmt.Sprintf("%v", value.Interface()))
		return nil
	},
	"float-to-string": func(target reflect.Value, value reflect.Value) error {
		target.SetString(fmt.Sprintf("%v", value.Interface()))
		return nil
	},
	"int-to-string": func(target reflect.Value, value reflect.Value) error {
		target.SetString(fmt.Sprintf("%v", value.Interface()))
		return nil
	},
	"uint-to-string": func(target reflect.Value, value reflect.Value) error {
		target.SetString(fmt.Sprintf("%v", value.Interface()))
		return nil
	},
}

// coerceType accepts a reflect.Value and returns a simplified logical type; for example float32 and float64
// are condensed into float; all ints (int, int8, int16, ...) are condensed into int.  Likewise for uint types.
// The second return value indicates if this type can be type-coerced.
func coerceType(v reflect.Value) (string, bool) {
	switch v.Kind() {
	case reflect.Bool:
		return "bool", true

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return "float", true

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return "int", true

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return "uint", true

	case reflect.String:
		return "string", true

	default:
		return v.Type().String(), false
	}
}

// coerce coerces the data in value to the correct type and assigns it to target.
func coerce(target reflect.Value, value reflect.Value) error {
	to, _ := coerceType(target)
	from, _ := coerceType(value)
	if fn, ok := coercions[from+"-to-"+to]; ok {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("Recovered %v", r)
				}
			}()
			target.Set(reflect.Zero(target.Type()))
			err = fn(target, value)
		}()
		return err
	}
	return errors.Errorf("Type coercion from %v to %v unsupported.", from, to)
}
