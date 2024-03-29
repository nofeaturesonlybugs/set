package set

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/nofeaturesonlybugs/set/path"
)

// BoundMapping is returned from Mapper's Bind method.
//
// A BoundMapping must not be copied except via its Copy method.
//
// A BoundMapping should be used in iterative code that needs to read or mutate
// many instances of the same struct.  Bound mappings allow for adhoc or indeterminate
// field access within the bound data.
//	// adhoc access means different fields can be accessed between calls to Rebind()
//	var a, b T
//
//	bound := myMapper.Map(&a)
//	bound.Set("Field", 10)
//	bound.Set("Other", "Hello")
//
//	bound.Rebind(&b)
//	bound.Set("Bar", 27)
//
// In the preceding example the BoundMapping is first bound to a and later bound to b
// and each instance had different field(s) accessed.
type BoundMapping struct {
	// top is the original type used to create the value; it is needed to ensure type compatibility
	// when calling Rebind.
	// value is the bound value after passing through Writable to get to the end of any pointer chain.
	top   reflect.Type
	value reflect.Value
	err   error

	// NB  This field should be treated as read-only.
	paths map[string]path.ReflectPath
}

// Assignables returns a slice of pointers to the fields in the currently bound struct
// in the order specified by the fields argument.
//
// To alleviate pressure on the garbage collector the return slice can be pre-allocated and
// passed as the second argument to Assignables.  If non-nil it is assumed len(fields) == len(rv)
// and failure to provide an appropriately sized non-nil slice will cause a panic.
//
// During traversal this method will allocate struct fields that are nil pointers.
//
// An example use-case would be obtaining a slice of pointers for Rows.Scan() during database
// query results.
func (b BoundMapping) Assignables(fields []string, rv []interface{}) ([]interface{}, error) {
	if b.err != nil && errors.Is(b.err, ErrReadOnly) {
		return rv, b.err.(pkgerr).WithCallSite("BoundMapping.Assignables")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldN, name := range fields {
		step, ok := b.paths[name]
		if !ok {
			err := pkgerr{
				Err:      ErrUnknownField,
				CallSite: "BoundMapping.Assignables",
				Context:  "field [" + name + "] not found in type " + b.top.String(),
			}
			return rv, err
		}
		v := b.value
		if step.HasPointer { // NB  Begin manual inline of path.ReflectPath.Value
			for _, n := range step.Index {
				v = v.Field(n)
				for ; v.Kind() == reflect.Ptr; v = v.Elem() {
					if v.IsNil() && v.CanSet() {
						v.Set(reflect.New(v.Type().Elem()))
					}
				}
			}
		} else {
			for _, n := range step.Index {
				v = v.Field(n)
			}
		}
		v = v.Field(step.Last) // NB  End manual inline of path.ReflectPath.Value
		rv[fieldN] = v.Addr().Interface()
	}
	return rv, nil
}

// Copy creates an exact copy of the BoundMapping.
//
// One use case for Copy is to create a set of BoundMappings early in a program's
// init phase.  During later execution when a BoundMapping is needed for type T
// it can be obtained by calling Copy on the cached BoundMapping for that type.
func (b BoundMapping) Copy() BoundMapping {
	return BoundMapping{
		top:   b.top,
		value: b.value,
		err:   b.err,
		paths: b.paths,
	}
}

// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
// calls to Rebind()
func (b BoundMapping) Err() error {
	return b.err
}

// Field returns the Value for field.
func (b BoundMapping) Field(field string) (Value, error) {
	if b.err != nil && errors.Is(b.err, ErrReadOnly) {
		return zeroV, b.err.(pkgerr).WithCallSite("BoundMapping.Field")
	}
	step, ok := b.paths[field]
	if !ok {
		err := pkgerr{
			Err:      ErrUnknownField,
			CallSite: "BoundMapping.Field",
			Context:  "field [" + field + "] not found in type " + b.top.String(),
		}
		return zeroV, err
	}
	v := b.value
	if step.HasPointer { // NB  Begin manual inline of path.ReflectPath.Value
		for _, n := range step.Index {
			v = v.Field(n)
			for ; v.Kind() == reflect.Ptr; v = v.Elem() {
				if v.IsNil() && v.CanSet() {
					v.Set(reflect.New(v.Type().Elem()))
				}
			}
		}
	} else {
		for _, n := range step.Index {
			v = v.Field(n)
		}
	}
	v = v.Field(step.Last) // NB  End manual inline of path.ReflectPath.Value
	//
	return V(v), nil
}

// Fields returns a slice of values to the fields in the currently bound struct in the order
// specified by the fields argument.
//
// To alleviate pressure on the garbage collector the return slice can be pre-allocated and
// passed as the second argument to Fields.  If non-nil it is assumed len(fields) == len(rv)
// and failure to provide an appropriately sized non-nil slice will cause a panic.
//
// During traversal this method will allocate struct fields that are nil pointers.
//
// An example use-case would be obtaining a slice of query arguments by column name during
// database queries.
func (b BoundMapping) Fields(fields []string, rv []interface{}) ([]interface{}, error) {
	if b.err != nil && errors.Is(b.err, ErrReadOnly) {
		return rv, b.err.(pkgerr).WithCallSite("BoundMapping.Fields")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldN, name := range fields {
		step, ok := b.paths[name]
		if !ok {
			err := pkgerr{
				Err:      ErrUnknownField,
				CallSite: "BoundMapping.Fields",
				Context:  "field [" + name + "] not found in type " + b.top.String(),
			}
			return rv, err
		}
		v := b.value
		if step.HasPointer { // NB  Begin manual inline of path.ReflectPath.Value
			for _, n := range step.Index {
				v = v.Field(n)
				for ; v.Kind() == reflect.Ptr; v = v.Elem() {
					if v.IsNil() && v.CanSet() {
						v.Set(reflect.New(v.Type().Elem()))
					}
				}
			}
		} else {
			for _, n := range step.Index {
				v = v.Field(n)
			}
		}
		v = v.Field(step.Last) // NB  End manual inline of path.ReflectPath.Value
		// NB  The value we want is v.Interface() which performs a number of allocations for built-in primitives.
		//     If we switch off v's type as a pointer and it is a primitive we can skip the allocations.
		switch ptr := v.Addr().Interface().(type) {
		case *bool:
			rv[fieldN] = *ptr
		case *int:
			rv[fieldN] = *ptr
		case *int8:
			rv[fieldN] = *ptr
		case *int16:
			rv[fieldN] = *ptr
		case *int32:
			rv[fieldN] = *ptr
		case *int64:
			rv[fieldN] = *ptr
		case *uint:
			rv[fieldN] = *ptr
		case *uint8:
			rv[fieldN] = *ptr
		case *uint16:
			rv[fieldN] = *ptr
		case *uint32:
			rv[fieldN] = *ptr
		case *uint64:
			rv[fieldN] = *ptr
		case *float32:
			rv[fieldN] = *ptr
		case *float64:
			rv[fieldN] = *ptr
		case *string:
			rv[fieldN] = *ptr
		case *time.Time:
			rv[fieldN] = *ptr
		default:
			rv[fieldN] = v.Interface()
		}
	}
	return rv, nil
}

// Rebind will replace the currently bound value with the new variable v.
//
// v must have the same type as the original value used to create the BoundMapping
// otherwise a panic will occur.
//
// As a convenience Rebind allows v to be an instance of reflect.Value.  This prevents
// unnecessary calls to reflect.Value.Interface().
func (b *BoundMapping) Rebind(v interface{}) {
	if b.err != nil && errors.Is(b.err, ErrReadOnly) {
		return
	}
	//
	// Allow reflect.Value to be passed directly.
	var rv reflect.Value
	var T reflect.Type
	switch sw := v.(type) {
	case reflect.Value:
		rv = sw
		T = rv.Type()
	default:
		rv = reflect.ValueOf(v)
		T = rv.Type()
	}
	if b.top != T {
		panic(fmt.Sprintf("mismatching types during Rebind; have %v and got %T", b.top.String(), v)) // TODO ErrRebind maybe?
	}
	b.err = nil
	b.value, _ = Writable(rv)
}

// Set effectively sets V[field] = value.
func (b *BoundMapping) Set(field string, value interface{}) error {
	if b.err != nil && errors.Is(b.err, ErrReadOnly) {
		return b.err.(pkgerr).WithCallSite("BoundMapping.Set")
	}
	//
	step, ok := b.paths[field]
	if !ok {
		err := pkgerr{
			Err:      ErrUnknownField,
			CallSite: "BoundMapping.Set",
			Context:  "field [" + field + "] not found in type " + b.top.String(),
		}
		if b.err == nil {
			b.err = err
		}
		return err
	}
	v := b.value
	if step.HasPointer { // NB  Begin manual inline of path.ReflectPath.Value
		for _, n := range step.Index {
			v = v.Field(n)
			for ; v.Kind() == reflect.Ptr; v = v.Elem() {
				if v.IsNil() && v.CanSet() {
					v.Set(reflect.New(v.Type().Elem()))
				}
			}
		}
	} else {
		for _, n := range step.Index {
			v = v.Field(n)
		}
	}
	v = v.Field(step.Last) // NB  End manual inline of path.ReflectPath.Value
	//
	// If the types are directly equatable then we might be able to avoid creating a V(fieldValue),
	// which will cut down our allocations and increase speed.
	if v.Type() == reflect.TypeOf(value) {
		switch tt := value.(type) {
		case bool:
			v.SetBool(tt)
			return nil
		case int:
			v.SetInt(int64(tt))
			return nil
		case int8:
			v.SetInt(int64(tt))
			return nil
		case int16:
			v.SetInt(int64(tt))
			return nil
		case int32:
			v.SetInt(int64(tt))
			return nil
		case int64:
			v.SetInt(tt)
			return nil
		case uint:
			v.SetUint(uint64(tt))
			return nil
		case uint8:
			v.SetUint(uint64(tt))
			return nil
		case uint16:
			v.SetUint(uint64(tt))
			return nil
		case uint32:
			v.SetUint(uint64(tt))
			return nil
		case uint64:
			v.SetUint(tt)
			return nil
		case float32:
			v.SetFloat(float64(tt))
			return nil
		case float64:
			v.SetFloat(tt)
			return nil
		case string:
			v.SetString(tt)
			return nil
		}
	}
	//
	// If the type-switch above didn't hit then we'll coerce the
	// fieldValue to a Value and use our swiss-army knife Value.To().
	err := V(v).To(value)
	if err != nil && b.err == nil {
		b.err = err // TODO Possibly wrap with more information.
	}
	return err
}
