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
	// Performance note:
	//	Initially known was map[reflect.Type]Mapping and required a sync.RWMutex to protect access.
	//	Similarly to the change in type_info_cache_t we changing both in favor of the sync.Map.
	known sync.Map
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
	rv := new_bound_mapping_t(v, mapping)
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
	var typeInfo TypeInfo
	switch tt := T.(type) {
	case *Value:
		typeInfo = tt.TypeInfo
	case reflect.Type:
		typeInfo = TypeCache.StatType(tt)
	default:
		typeInfo = TypeCache.Stat(T)
	}
	//
	if rv, ok := me.known.Load(typeInfo.Type); ok {
		return rv.(Mapping), nil
	}
	//
	rv := make(Mapping)
	//
	var scan func(typeInfo TypeInfo, indeces []int, prefix string)
	scan = func(typeInfo TypeInfo, indeces []int, prefix string) {
		for k, field := range typeInfo.StructFields {
			fieldTypeInfo := TypeCache.StatType(field.Type)
			if me.Ignored.Has(fieldTypeInfo.Type) {
				continue
			}
			//
			name := ""
			if !me.Elevated.Has(fieldTypeInfo.Type) {
				for _, tagName := range append(me.Tags, "") {
					if tagValue, ok := field.Tag.Lookup(tagName); ok {
						name = tagValue
						break
					} else if tagName == "" {
						name = field.Name
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
			if _, ok := mapper_TreatAsScalar[fieldTypeInfo.Type]; ok {
				rv[name] = nameIndeces
			} else if fieldTypeInfo.IsStruct {
				scan(fieldTypeInfo, nameIndeces, name)
			} else if fieldTypeInfo.IsScalar {
				rv[name] = nameIndeces
			}
		}
	}
	// Scan and assign the result to our known types.
	scan(typeInfo, []int{}, "")
	me.known.Store(typeInfo.Type, rv)
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

// new_bound_mapping_t creates a new bound_mapping_t type.
func new_bound_mapping_t(value *Value, mapping Mapping) *bound_mapping_t {
	return &bound_mapping_t{
		value:   value,
		mapping: mapping,
	}
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
	if v, err := me.value.FieldByIndex(me.mapping.Get(field)); err != nil {
		return nil, errors.Go(err)
	} else {
		return V(v), nil
	}
}

// Rebind will replace the currently bound value with the new variable I.  If the underlying types are
// not the same then an error is returned.
func (me *bound_mapping_t) Rebind(I interface{}) error {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	me.err = nil
	return me.value.Rebind(I)
}

// Set effectively sets V[field] = value.
func (me *bound_mapping_t) Set(field string, value interface{}) error {
	// nil check is not necessary as bound_mapping_t is only created within this package.
	fieldValue, err := me.value.FieldByIndex(me.mapping.Get(field))
	if err != nil {
		if me.err == nil {
			me.err = errors.Go(err)
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
	if err != nil && me.err == nil {
		me.err = errors.Errorf("While setting [%v]: %v", field, err.Error())
	}
	return err
}
