package set

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nofeaturesonlybugs/errors"
)

var mapper_TreatAsScalar = map[reflect.Type]struct{}{
	reflect.TypeOf(time.Time{}):  {},
	reflect.TypeOf(&time.Time{}): {},
}

// Mapping is the result of traversing nested structures to generate a mapping of Key-to-Indeces where:
//	Key is a string key representing a common or friendly name for the nested struct member.
//	Indeces is an []int that can be used to index to the proper struct member.
type Mapping map[string][]int

// Mapper is used to traverse structures to create Mappings and then navigate the nested
// structure using string keys.
type Mapper struct {
	// A slice of type instances that should be ignored during name generation.
	Ignored TypeList
	//
	// During name generation types in the `Elevated` member will not have the name affected
	// by the struct field name or data type.  Use this for struct members or embedded structs
	// when you do not want their name or type affecting the generated name.
	Elevated TypeList
	//
	// A list of struct tags that will be used for name generation in order of preference.
	// An example would be using this feature for both JSON and DB field name specification.
	// If most of your db and json names match but you occasionally want to override the json
	// struct tag value with the db struct tag value you could set this member to:
	//	[]string{ "db", "json" } // struct tag `db` used before struct tag `json`
	Tags []string
	//
	// Join specifies the string used to join generated names as nesting increases.
	Join string
	//
	// If set this function is called when the struct field name is being used as
	// the generated name.  This function can perform string alteration to force all
	// names to lowercase, string replace, etc.
	Transform func(string) string

	//
	mut   sync.RWMutex
	known map[reflect.Type]Mapping // Types that are known -- i.e. we've already created the mapping.
}

// BoundMapping binds a Mapping to a specific variable instance V.
type BoundMapping interface {
	// Assignables returns a slice of interfaces by field names where each element is a pointer
	// to the field in the bound variable.
	Assignables(fields []string) ([]interface{}, error)
	// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
	// calls to Rebind()
	Err() error
	// Field returns the *Value for field.
	Field(field string) (*Value, error)
	// Rebind will replace the currently bound value with the new variable I.  If the underlying types are
	// not the same then an error is returned.
	Rebind(I interface{}) error
	// Set effectively sets V[field] = value.
	Set(field string, value interface{}) error
}

var DefaultMapper = &Mapper{
	Join: "_",
}

// Bind creates a Mapping bound to a specific instance I of a variable.
func (me *Mapper) Bind(I interface{}) (BoundMapping, error) {
	var v *Value
	if tv, ok := I.(*Value); ok {
		v = tv
	} else {
		v = V(I)
	}
	mapping, err := me.Map(v) // It's this call to Map() that performs the nil receiver check.
	if err != nil {
		return nil, errors.Go(err)
	}
	rv := &bound_mapping_t{
		value:   v,
		mapping: mapping,
	}
	return rv, nil
}

// Map adds T to the Mapper's list of known and recognized types.
//
// Map is written to be goroutine safe in that multiple goroutines can call it.  If multiple goroutines call
// Map simultaneously it is not guaranteed that each goroutine receives the same Mapping instance;
// however it is guaranteed that each instance will behave identically.
//
// If you depend on Map returning the same instance of Mapping for a type T every time it is called then
// you should call it once for each possible type T synchronously before using the Mapper in your goroutines.
//
// Mappings that are returned should not be written to or altered in any way.  If this is your use-case then
// create a copy of the Mapping with Mapping.Copy.
func (me *Mapper) Map(T interface{}) (Mapping, error) {
	if me == nil {
		return nil, errors.NilReceiver()
	}
	var v *Value
	if tv, ok := T.(*Value); ok {
		v = tv
	} else {
		v = V(T)
	}
	//
	me.mut.RLock()
	if rv, ok := me.known[v.Type]; ok {
		me.mut.RUnlock()
		return rv, nil
	}
	me.mut.RUnlock()
	//
	rv := make(Mapping)
	//
	var scan func(v *Value, indeces []int, prefix string, indent int)
	scan = func(v *Value, indeces []int, prefix string, indent int) {
		for k, field := range v.Fields() {
			if me.Ignored.Has(field.Value.Type) {
				continue
			}
			//
			name := ""
			if !me.Elevated.Has(field.Value.Type) {
				for _, tagName := range append(me.Tags, "") {
					if tagValue, ok := field.Field.Tag.Lookup(tagName); ok {
						name = tagValue
						break
					} else if tagName == "" {
						name = field.Field.Name
						if me.Transform != nil {
							name = me.Transform(name)
						}
						break
					}
				}
			}
			if prefix != "" && name != "" {
				name = prefix + me.Join + name
			} else if prefix != "" {
				name = prefix
			}
			nameIndeces := append(indeces, k)
			if _, ok := mapper_TreatAsScalar[field.Value.Type]; ok {
				rv[name] = nameIndeces
			} else if field.Value.IsStruct {
				scan(field.Value, nameIndeces, name, indent+1)
			} else if field.Value.IsScalar {
				rv[name] = nameIndeces
			}
		}
	}
	scan(v, []int{}, "", 0)
	//
	// Acquiring the write lock is delayed until after scanning because the desired use case is faster
	// reads; the following lines are the smallest set in which the lock *MUST* be held.
	//
	// Our scan is complete so now we should assign the result to our known types.
	me.mut.Lock()
	defer me.mut.Unlock()
	if me.known == nil {
		me.known = make(map[reflect.Type]Mapping)
	}
	me.known[v.Type] = rv
	//
	return rv, nil
}

// Copy creates a copy of the Mapping.
func (me Mapping) Copy() Mapping {
	if me == nil {
		return nil
	}
	rv := make(Mapping)
	for k, v := range me {
		rv[k] = append([]int{}, v...)
	}
	return rv
}

// Get returns the indeces associated with key in the mapping.  If no such key
// is found a nil slice is returned.
func (me Mapping) Get(key string) []int {
	v, _ := me.Lookup(key)
	return v
}

// Lookup returns the value associated with key in the mapping.  If no such key is
// found a nil slice is returned and ok is false; otherwise ok is true.
func (me Mapping) Lookup(key string) (indeces []int, ok bool) {
	if me == nil {
		return nil, false
	}
	indeces, ok = me[key]
	return indeces, ok
}

// String returns the Mapping as a string value.
func (me Mapping) String() string {
	parts := []string{}
	for str, indeces := range me {
		parts = append(parts, fmt.Sprintf("%v\t\t%v", indeces, str))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}

// bound_mapping_t is the implementation for BoundMapping.
type bound_mapping_t struct {
	value   *Value
	mapping Mapping
	err     error
}

// Assignables returns a slice of interfaces by field names where each element is a pointer
// to the field in the bound variable.
//
// An example use-case would be obtaining a slice of pointers for Rows.Scan() during database
// query results.
func (me *bound_mapping_t) Assignables(fields []string) ([]interface{}, error) {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	rv := []interface{}{}
	for _, name := range fields {
		if field, err := me.Field(name); err != nil {
			return nil, errors.Errorf("%v while accessing field [%v]", err.Error(), name)
		} else {
			rv = append(rv, field.WriteValue.Addr().Interface())
		}
	}
	return rv, nil
}

// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
// calls to Rebind()
func (me *bound_mapping_t) Err() error {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	return me.err
}

// Field returns the *Value for field.
func (me *bound_mapping_t) Field(field string) (*Value, error) {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	return me.value.FieldByIndex(me.mapping.Get(field))
}

// Rebind will replace the currently bound value with the new variable I.  If the underlying types are
// not the same then an error is returned.
func (me *bound_mapping_t) Rebind(I interface{}) error {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	var v *Value
	if tv, ok := I.(*Value); ok {
		v = tv
	} else {
		v = V(I)
	}
	if v.Type != me.value.Type {
		return errors.Errorf("Rebind expects same underlying type; had %T and got %T", me.value.WriteValue.Interface(), v.WriteValue.Interface())
	}
	me.err = nil
	me.value = v
	return nil
}

// Set effectively sets V[field] = value.
func (me *bound_mapping_t) Set(field string, value interface{}) error {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	v, err := me.value.FieldByIndex(me.mapping.Get(field))
	if err != nil {
		if me.err == nil {
			me.err = errors.Go(err)
		}
		return errors.Go(err)
	}
	err = v.To(value)
	if err != nil && me.err == nil {
		me.err = errors.Errorf("While setting [%v]: %v", field, err.Error())
	}
	return err
}
