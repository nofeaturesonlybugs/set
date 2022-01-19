package set

import (
	"reflect"

	"github.com/nofeaturesonlybugs/errors"
)

// BoundMapping is returned from Mapper's Bind() method.
//
// A BoundMapping must not be copied except via its Copy method.
type BoundMapping struct {
	value *Value
	err   error

	// NB   indeces is obtained from Mapping.Indeces and is not a copy.
	//      Treat as read only.
	indeces map[string][]int
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
func (b BoundMapping) Assignables(fields []string, rv []interface{}) ([]interface{}, error) {
	if !b.value.CanWrite {
		return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldnum, name := range fields {
		indeces := b.indeces[name]
		if len(indeces) == 0 {
			return nil, errors.Errorf("No mapping for field [%v]", name)
		}
		v := b.value.WriteValue
		for k, size := 0, len(indeces); k < size; k++ {
			n := indeces[k] // n is the specific FieldIndex into v
			if k < size-1 {
				// n is not the final index; therefore it is a nested/embedded struct and we might need to instantiate it.
				v, _ = Writable(v.Field(n))
			} else {
				// n is the final index; it represents a scalar/leaf at the end of the nested struct hierarchy.
				rv[fieldnum] = v.Field(n).Addr().Interface()
			}
		}
	}
	return rv, nil
}

// Copy creates an exact copy of the BoundMapping.  Since a BoundMapping is implicitly tied to a single
// value sometimes it may be useful to create and cache a BoundMapping early in a program's initialization and
// then call Copy() to work with multiple instances of bound values simultaneously.
func (b BoundMapping) Copy() BoundMapping {
	return BoundMapping{
		value: b.value.Copy(),
		err:   b.err,
		// NB   Don't need to copy me.mapping since we never alter it.
		indeces: b.indeces,
	}
}

// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
// calls to Rebind()
func (b BoundMapping) Err() error {
	return b.err
}

// Field returns the *Value for field.
func (b BoundMapping) Field(field string) (*Value, error) {
	var v reflect.Value
	var err error
	if v, err = b.value.FieldByIndex(b.indeces[field]); err != nil {
		return nil, errors.Go(err)
	}
	return V(v), nil
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
func (b BoundMapping) Fields(fields []string, rv []interface{}) ([]interface{}, error) {
	if !b.value.CanWrite {
		return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldnum, name := range fields {
		indeces := b.indeces[name]
		if len(indeces) == 0 {
			return nil, errors.Errorf("No mapping for field [%v]", name)
		}
		v := b.value.WriteValue
		for k, size := 0, len(indeces); k < size; k++ {
			n := indeces[k] // n is the specific FieldIndex into v
			if k < size-1 {
				// n is not the final index; therefore it is a nested/embedded struct and we might need to instantiate it.
				v, _ = Writable(v.Field(n))
			} else {
				// n is the final index; it represents a scalar/leaf at the end of the nested struct hierarchy.
				rv[fieldnum] = v.Field(n).Interface()
			}
		}
	}
	return rv, nil
}

// Rebind will replace the currently bound value with the new variable I.
func (b *BoundMapping) Rebind(I interface{}) {
	// ^^ me.value is a pointer so a value receiver is ok.
	b.err = nil
	b.value.Rebind(I)
}

// Set effectively sets V[field] = value.
func (b *BoundMapping) Set(field string, value interface{}) error {
	// nil check is not necessary as boundMapping is only created within this package.
	indeces := b.indeces[field]
	fieldValue, err := b.value.FieldByIndex(indeces)
	if err != nil {
		if b.err == nil {
			b.err = errors.Go(err)
			return b.err
		}
		return errors.Go(err)
	}
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
	if err != nil && b.err == nil {
		b.err = errors.Errorf("While setting [%v]: %v", field, err.Error())
	}
	return err
}
