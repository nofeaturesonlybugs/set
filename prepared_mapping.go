package set

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set/path"
)

var (
	// ErrPlanExceeded is returned when the next access to a PreparedMapping exceeds the
	// fields specified by the earlier call to Plan.
	ErrPlanExceeded = fmt.Errorf("set: prepared mapping: attempted access extends plan")

	// ErrPlanInvalid is returned when a PreparedMapping does not have a valid access plan.
	ErrPlanInvalid = fmt.Errorf("set: prepared mapping: plan invalid")
)

// preparedMappingField is a prepare field access.
type preparedMappingField struct {
	path.Path
}

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
	value *Value
	err   error

	// plan is the slice of steps created by Plan and
	// k is the index into plan for the next step.
	k    int
	plan []preparedMappingField

	// NB   indeces is obtained from Mapping.Indeces and is not a copy.
	//      Treat as read only.
	indeces map[string][]int
	// NB	paths is obtained from Mapping.Paths and is not a copy.
	//      Treat as read only.
	paths map[string]path.Path
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
	if p.err == ErrPlanInvalid || p.err == ErrReadOnly {
		return rv, p.err
	}
	var v reflect.Value
	if rv == nil {
		rv = make([]interface{}, len(p.plan))
	}
	for n, path := range p.plan {
		v = path.Value(p.value.WriteValue)
		rv[n] = v.Addr().Interface()
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
		value:   p.value.Copy(),
		err:     p.err,
		k:       p.k,
		plan:    append([]preparedMappingField(nil), p.plan...),
		indeces: p.indeces,
		paths:   p.paths,
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
	if p.err == ErrPlanInvalid || p.err == ErrReadOnly {
		return nil, p.err
	}
	//
	var fieldValue reflect.Value
	var err error
	//
	if err = p.next(); err != nil {
		return nil, err
	}
	//
	path := p.plan[p.k].Path
	fieldValue = path.Value(p.value.WriteValue)
	return V(fieldValue), nil
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
	if p.err == ErrPlanInvalid || p.err == ErrReadOnly {
		return rv, p.err
	}
	var v reflect.Value
	if rv == nil {
		rv = make([]interface{}, len(p.plan))
	}
	for n, path := range p.plan {
		v = path.Value(p.value.WriteValue)
		rv[n] = v.Interface()
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
	if p.err == ErrReadOnly {
		return p.err
	}
	var path path.Path
	var max int
	var ok bool
	//
	p.err = nil
	p.k = -1
	//
	// Ensure p.plan has enough space and then reslice to empty.
	if max = len(fields); max > cap(p.plan) {
		p.plan = make([]preparedMappingField, 0, max)
	}
	p.plan = p.plan[0:0]
	//
	for _, field := range fields {
		if path, ok = p.paths[field]; !ok {
			p.err = ErrPlanInvalid
			return fmt.Errorf("%w: %v", ErrUnknownField, field)
		}
		p.plan = append(p.plan, preparedMappingField{Path: path})
	}
	// TODO Rebuild cache for current object
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
	if p.err == ErrReadOnly {
		return
	}
	p.err = nil
	p.k = -1
	// TODO When refactoring this to remove p.value copy+paste the implementation in BoundMapping.
	p.value.Rebind(v)
}

// next attempts to advance the internal counter by one.  If advancing the counter
// exceeds len(plan) then ErrPlanExceeded is returned.
func (p *PreparedMapping) next() error {
	p.k++
	if p.k == len(p.plan) {
		p.k--
		if p.err == nil {
			p.err = ErrPlanExceeded
		}
		return ErrPlanExceeded
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
	if p.err == ErrPlanInvalid || p.err == ErrReadOnly {
		return p.err
	}
	//
	var fieldValue reflect.Value
	var err error
	//
	if err = p.next(); err != nil {
		return err
	}
	//
	path := p.plan[p.k].Path
	fieldValue = path.Value(p.value.WriteValue)
	//
	// If the types are directly equatable then we might be able to avoid creating a V(fieldValue),
	// which will cut down our allocations and increase speed.
	if fieldValue.Type() == reflect.TypeOf(value) {
		switch tt := value.(type) {
		case bool:
			fieldValue.SetBool(tt)
			return nil
		case int:
			fieldValue.SetInt(int64(tt))
			return nil
		case int8:
			fieldValue.SetInt(int64(tt))
			return nil
		case int16:
			fieldValue.SetInt(int64(tt))
			return nil
		case int32:
			fieldValue.SetInt(int64(tt))
			return nil
		case int64:
			fieldValue.SetInt(tt)
			return nil
		case uint:
			fieldValue.SetUint(uint64(tt))
			return nil
		case uint8:
			fieldValue.SetUint(uint64(tt))
			return nil
		case uint16:
			fieldValue.SetUint(uint64(tt))
			return nil
		case uint32:
			fieldValue.SetUint(uint64(tt))
			return nil
		case uint64:
			fieldValue.SetUint(tt)
			return nil
		case float32:
			fieldValue.SetFloat(float64(tt))
			return nil
		case float64:
			fieldValue.SetFloat(tt)
			return nil
		case string:
			fieldValue.SetString(tt)
			return nil
		}
	}
	//
	// If the type-switch above didn't hit then we'll coerce the
	// fieldValue to a *Value and use our swiss-army knife Value.To().
	err = V(fieldValue).To(value)
	if err != nil && p.err == nil {
		p.err = err // TODO Possibly wrap with more information.
	}
	return err
}
