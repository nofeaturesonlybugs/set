package set

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/nofeaturesonlybugs/set/path"
)

// PreparedMapping is returned from Mapper's Prepare method.
//
// A PreparedMapping must not be copied except via its Copy method.
//
// PreparedMappings should be used in iterative code that needs to read or mutate
// many instances of the same struct.  PreparedMappings do not allow for indeterminate
// field access between instances -- every struct instance must have the same fields
// accessed in the same order.  This behavior is akin to prepared statements in
// a database engine; if you need adhoc or indeterminate access use a BoundMapping.
//	var a, b T
//
//	p := myMapper.Prepare(&a)
//	_ = p.Plan("Field", "Other") // check err in production
//
//	p.Set(10)          // a.Field = 10
//	p.Set("Hello")     // a.Other = "Hello"
//
//	p.Rebind(&b)       // resets internal plan counter
//	p.Set(27)          // b.Field = 27
//	p.Set("World")     // b.Other = "World"
//
// All methods that return an error will return ErrPlanInvalid until Plan is called
// specifying an access plan.  Methods that do not return an error can be called before
// a plan has been specified.
type PreparedMapping struct {
	// top is the original type used to create the value; it is needed to ensure type compatibility
	// when calling Rebind.
	// value is the bound value after passing through Writable to get to the end of any pointer chain.
	top   reflect.Type
	value reflect.Value

	// valid=true means PreparedMapping is valid and bound to a writable value.
	// valid=false means PreparedMapping is invalid and most calls should return err.
	valid bool
	err   error

	// plan is the slice of steps created by Plan and
	// k is the index into plan for the next step.
	k    int
	plan []path.ReflectPath

	// NB	paths is obtained from Mapping.Paths and is not a copy.
	//      Treat as read only.
	paths map[string]path.ReflectPath
}

// Assignables returns a slice of pointers to the fields in the currently bound
// struct in the order specified by the last call to Plan.
//
// To alleviate pressure on the garbage collector the return slice can be pre-allocated and
// passed as the argument to Assignables.  If non-nil it is assumed len(plan) == len(rv)
// and failure to provide an appropriately sized non-nil slice will cause a panic.
//
// During traversal this method will allocate struct fields that are nil pointers.
//
// An example use-case would be obtaining a slice of pointers for Rows.Scan() during database
// query results.
func (p PreparedMapping) Assignables(rv []interface{}) ([]interface{}, error) {
	if !p.valid {
		return rv, p.err.(pkgerr).WithCallSite("PreparedMapping.Assignables")
	}
	if rv == nil {
		rv = make([]interface{}, len(p.plan))
	}
	for fieldN, step := range p.plan {
		v := p.value
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

// Copy creates an exact copy of the PreparedMapping.
//
// One use case for Copy is to create a set of PreparedMappings early in a program's
// init phase.  During later execution when a PreparedMapping is needed for type T
// it can be obtained by calling Copy on the cached PreparedMapping for that type.
func (p PreparedMapping) Copy() PreparedMapping {
	return PreparedMapping{
		top:   p.top,
		value: p.value,
		valid: p.valid,
		err:   p.err,
		k:     p.k,
		plan:  append([]path.ReflectPath(nil), p.plan...),
		paths: p.paths,
	}
}

// Err returns an error that may have occurred during repeated calls to Set.
//
// Err is reset on calls to Plan or Rebind.
func (p PreparedMapping) Err() error {
	return p.err
}

// Field returns the *Value for the next field.
//
// Each call to Field advances the internal access pointer in order to traverse the
// fields in the same order as the last call to Plan.
//
// ErrPlanInvalid is returned if Plan has not been called.  If this call to Field
// exceeds the length of the plan then ErrPlanExceeded is returned.  Other
// errors from this package or standard library may also be returned.
func (p *PreparedMapping) Field() (*Value, error) {
	if !p.valid {
		return nil, p.err.(pkgerr).WithCallSite("PreparedMapping.Field")
	}
	//
	if err := p.next(); err != nil {
		return nil, err.(pkgerr).WithCallSite("PreparedMapping.Field")
	}
	//
	step := p.plan[p.k]
	v := p.value
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
// specified by the last call to Plan.
//
// To alleviate pressure on the garbage collector the return slice can be pre-allocated and
// passed as the argument to Fields.  If non-nil it is assumed len(plan) == len(rv)
// and failure to provide an appropriately sized non-nil slice will cause a panic.
//
// During traversal this method will allocate struct fields that are nil pointers.
//
// An example use-case would be obtaining a slice of query arguments by column name during
// database queries.
func (p PreparedMapping) Fields(rv []interface{}) ([]interface{}, error) {
	if !p.valid {
		return rv, p.err.(pkgerr).WithCallSite("PreparedMapping.Fields")
	}
	if rv == nil {
		rv = make([]interface{}, len(p.plan))
	}
	for fieldN, step := range p.plan {
		v := p.value
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

// Plan builds the field access plan and must be called before any other methods that
// return an error.
//
// Each call to plan:
//	1. Resets any internal error to nil
//	2. Resets the internal plan-step counter.
//
// If an unknown field is specified then ErrUnknownField is wrapped with the field name
// and the internal error is set to ErrPlanInvalid.
func (p *PreparedMapping) Plan(fields ...string) error {
	if p.err != nil && errors.Is(p.err, ErrReadOnly) {
		return p.err.(pkgerr).WithCallSite("PreparedMapping.Plan")
	}
	//
	// Ensure p.plan has enough space and then reslice to empty.
	if max := len(fields); max > cap(p.plan) {
		p.plan = make([]path.ReflectPath, 0, max)
	}
	p.plan = p.plan[0:0]
	// p.valid=false // TODO+NB Needs test case
	//
	for _, field := range fields {
		path, ok := p.paths[field]
		if !ok {
			context := "field [" + field + "] not found in type " + p.top.String()
			p.err = pkgerr{
				Err:      ErrNoPlan,
				CallSite: "PreparedMapping.Plan",
				Context:  context + " during plan creation",
			}
			return pkgerr{
				Err:      ErrUnknownField,
				Context:  context,
				CallSite: "PreparedMapping.Plan",
			}
		}
		p.plan = append(p.plan, path)
	}
	//
	p.err = nil
	p.k = -1
	p.valid = true
	//
	return nil
}

// Rebind will replace the currently bound value with the new variable v.
//
// v must have the same type as the original value used to create the PreparedMapping
// otherwise a panic will occur.
//
// As a convenience Rebind allows v to be an instance of reflect.Value.  This prevents
// unnecessary calls to reflect.Value.Interface().
func (p *PreparedMapping) Rebind(v interface{}) {
	if p.err != nil && errors.Is(p.err, ErrReadOnly) {
		return
	}
	//
	// Allow reflect.Value to be passed directly.
	var rv reflect.Value
	switch sw := v.(type) {
	case reflect.Value:
		rv = sw
	default:
		rv = reflect.ValueOf(v)
	}
	T := rv.Type()
	if p.top != T {
		panic(fmt.Sprintf("mismatching types during Rebind; have %T and got %T", p.value.Interface(), v)) // TODO ErrRebind maybe?
	}
	//
	if p.valid {
		// Only clear previous error if we are valid.
		p.err = nil
	}
	p.value, _ = Writable(rv)
	//
	p.k = -1
}

// next attempts to advance the internal counter by one.  If advancing the counter
// exceeds len(plan) then ErrPlanExceeded is returned.
func (p *PreparedMapping) next() error {
	p.k++
	if p.k == len(p.plan) {
		p.k--
		err := pkgerr{Err: ErrPlanOutOfBounds, Context: "value of " + p.top.String()}
		if p.err == nil {
			p.err = err
		}
		return err
	}
	return nil
}

// Set effectively sets V[field] = value.
//
// Each call to Set advances the internal access pointer in order to traverse the
// fields in the same order as the last call to Plan.
//
// ErrPlanInvalid is returned if Plan has not been called.  If this call to Set
// exceeds the length of the plan then ErrPlanExceeded is returned.  Other
// errors from this package or standard library may also be returned.
func (p *PreparedMapping) Set(value interface{}) error {
	if !p.valid {
		return p.err.(pkgerr).WithCallSite("PreparedMapping.Set")
	}
	//
	var err error
	//
	if err = p.next(); err != nil {
		return err.(pkgerr).WithCallSite("PreparedMapping.Set")
	}
	//
	v := p.value
	step := p.plan[p.k]
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
	// fieldValue to a *Value and use our swiss-army knife Value.To().
	err = V(v).To(value)
	if err != nil && p.err == nil {
		p.err = pkgerr{
			Err:      err,
			CallSite: "PreparedMapping.Set",
		}
	}
	return err
}
