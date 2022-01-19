package set

import (
	"reflect"

	"github.com/nofeaturesonlybugs/errors"

	"github.com/nofeaturesonlybugs/set/path"
)

// preparedMappingField is a prepare field access.
type preparedMappingField struct {
	path.Path
}

// PreparedMapping is returned from Mapper's Prepare method.
//
// Do not create PreparedMapping types any other way.
//
// PreparedMappings should be used in iterative code that needs to read or mutate
// many instances of the same struct.  PreparedMappings do not allow for indeterminate
// field access between instances -- every struct instance must have the same fields
// accessed in the same order.  This behavior is akin to prepared statements in
// a database engine; if you need more adhoc or indeterminate access use a BoundMapping.
//
// The Plan() method must be called with the field names that are intended to be
// accessed before any other method calls are valid, with the exception of Rebind().
type PreparedMapping struct {
	value *Value
	err   error

	// prepared plan values
	//
	// k is the index into the plan and must be incremented at the
	// beginning of every method call except Rebind() where it shall
	// be reset to k=-1
	k    int
	plan []preparedMappingField

	// NB   indeces is obtained from Mapping.Indeces and is not a copy.
	//      Treat as read only.
	indeces map[string][]int
	// NB	paths is obtained from Mapping.Paths and is not a copy.
	//      Treat as read only.
	paths map[string]path.Path
}

// Assignables returns a slice of interfaces by field names where each element is a pointer
// to the field in the bound variable.
//
// The second argument, if non-nil, will be the first return value.  In other words pre-allocating
// rv will cut down on memory allocations.  It is assumed that if rv is non-nil that len(fields) == len(rv)
// and the lengths are not checked, meaning your program could panic.
//
// During traversal this method will instantiate struct fields that are themselves structs.
//
// An example use-case would be obtaining a slice of pointers for Rows.Scan() during database
// query results.
func (p PreparedMapping) Assignables(fields []string, rv []interface{}) ([]interface{}, error) {
	panic("PreparedMapping.Assignables not imlemented") // TODO RM
	// TODO
	// if !p.value.CanWrite {
	// 	return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	// }
	// if rv == nil {
	// 	rv = make([]interface{}, len(fields))
	// }
	// for fieldnum, name := range fields {
	// 	indeces := p.indeces[name]
	// 	if len(indeces) == 0 {
	// 		return nil, errors.Errorf("No mapping for field [%v]", name)
	// 	}
	// 	v := p.value.WriteValue
	// 	for k, size := 0, len(indeces); k < size; k++ {
	// 		n := indeces[k] // n is the specific FieldIndex into v
	// 		if k < size-1 {
	// 			// n is not the final index; therefore it is a nested/embedded struct and we might need to instantiate it.
	// 			v, _ = Writable(v.Field(n))
	// 		} else {
	// 			// n is the final index; it represents a scalar/leaf at the end of the nested struct hierarchy.
	// 			rv[fieldnum] = v.Field(n).Addr().Interface()
	// 		}
	// 	}
	// }
	// return rv, nil
}

// Copy creates an exact copy of the BoundMapping.  Since a BoundMapping is implicitly tied to a single
// value sometimes it may be useful to create and cache a BoundMapping early in a program's initialization and
// then call Copy() to work with multiple instances of bound values simultaneously.
func (p PreparedMapping) Copy() PreparedMapping {
	panic("PreparedMapping.Copy not implemented") // TODO RM
	// TODO
	// return BoundMapping{
	// 	value: p.value.Copy(),
	// 	err:   p.err,
	// 	// NB   Don't need to copy me.mapping since we never alter it.
	// 	indeces: p.indeces,
	// }
}

// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
// calls to Rebind()
func (p PreparedMapping) Err() error {
	return p.err
}

// Field returns the *Value for field.
func (p PreparedMapping) Field(field string) (*Value, error) {
	panic("PreparedMapping.Field not yet implemented") // TODO RM
	// TODO
	// var v reflect.Value
	// var err error
	// if v, err = p.value.FieldByIndex(p.indeces[field]); err != nil {
	// 	return nil, errors.Go(err)
	// }
	// return V(v), nil
}

// Fields returns a slice of interfaces by field names where each element is the field value.
//
// The second argument, if non-nil, will be the first return value.  In other words pre-allocating
// rv will cut down on memory allocations.  It is assumed that if rv is non-nil that len(fields) == len(rv)
// and the lengths are not checked, meaning your program might panic if rv is not the correct length.
//
// During traversal this method will instantiate struct fields that are themselves structs.
//
// An example use-case would be obtaining a slice of query arguments by column name during
// database queries.
func (p PreparedMapping) Fields(fields []string, rv []interface{}) ([]interface{}, error) {
	panic("PreparedMapping.Fields not yet implemented") // TODO RM
	// TODO
	// if !p.value.CanWrite {
	// 	return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	// }
	// if rv == nil {
	// 	rv = make([]interface{}, len(fields))
	// }
	// for fieldnum, name := range fields {
	// 	indeces := p.indeces[name]
	// 	if len(indeces) == 0 {
	// 		return nil, errors.Errorf("No mapping for field [%v]", name)
	// 	}
	// 	v := p.value.WriteValue
	// 	for k, size := 0, len(indeces); k < size; k++ {
	// 		n := indeces[k] // n is the specific FieldIndex into v
	// 		if k < size-1 {
	// 			// n is not the final index; therefore it is a nested/embedded struct and we might need to instantiate it.
	// 			v, _ = Writable(v.Field(n))
	// 		} else {
	// 			// n is the final index; it represents a scalar/leaf at the end of the nested struct hierarchy.
	// 			rv[fieldnum] = v.Field(n).Interface()
	// 		}
	// 	}
	// }
	// return rv, nil
}

// Plan builds the field access plan and must be called before calling other methods
// with the exception of Rebind().  In other words Rebind() is the only method that can
// be called while a plan is not in place.
func (p *PreparedMapping) Plan(fields ...string) error {
	var path path.Path
	var ok bool
	//
	p.k = -1
	p.plan = nil
	for _, field := range fields {
		if path, ok = p.paths[field]; !ok {
			return errors.Errorf("unknown field %v", field)
		}
		p.plan = append(p.plan, preparedMappingField{path})
	}
	// TODO Rebuild cache for current object
	return nil
}

// Rebind will replace the currently bound value with the new variable I.
func (p *PreparedMapping) Rebind(I interface{}) {
	// ^^ me.value is a pointer so a value receiver is ok.
	p.err = nil
	p.value.Rebind(I)
	p.k = -1
}

// Set effectively sets V[field] = value.
func (p *PreparedMapping) Set(value interface{}) error {
	var fieldValue reflect.Value
	var err error
	// Increment plan counter and compare with length.
	if p.k = p.k + 1; !(p.k < len(p.plan)) {
		err = errors.Errorf("access beyond planned steps") // TODO Use an ErrPlanAccess or similar
		if p.err != nil {
			p.err = err
		}
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
		p.err = errors.Errorf("While setting [k=%v]: %v", p.k, err.Error()) // TODO Better meta than p.k?
	}
	return err
}
