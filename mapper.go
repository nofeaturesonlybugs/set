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

var mapperTreatAsScalar = map[reflect.Type]struct{}{
	reflect.TypeOf(time.Time{}):  {},
	reflect.TypeOf(&time.Time{}): {},
}

// Mapping is the result of traversing nested structures to generate a mapping of Key-to-Fields.
type Mapping struct {
	// Keys is a slice of key names in the order they were encountered during the mapping.
	Keys []string
	// Using a mapped name as a key the []int is the slice of indeces needed to traverse
	// the mapped struct to the nested field.
	Indeces map[string][]int
	// Using a mapped name as a key the StructField is the nested field.  Access to this
	// member is useful if your client package needs to inspect struct-fields-by-mapped-name; such
	// as in obtaining struct tags.
	StructFields map[string]reflect.StructField
	// CacheInfo map[string]struct { // TODO
	// 	Indeces  []int
	// 	SharedId int
	// }
}

// Mapper is used to traverse structures to create Mappings and then navigate the nested
// structure using string keys.
//
// Instantiate mappers as pointers:
//	myMapper := &set.Mapper{}
type Mapper struct {
	// If the types you wish to map contain embedded structs or interfaces you do not
	// want to map to string names include those types in the Ignored member.
	//
	// See also NewTypeList().
	Ignored TypeList
	//
	// Struct fields that are also structs or embedded structs will have their name
	// as part of the generated name unless it is included in the Elevated member.
	//
	// See also NewTypeList().
	Elevated TypeList
	//
	// Types in this list are treated as scalars when generating mappings; in other words
	// their exported fields are not mapped and the mapping created targets the type as
	// a whole.  This is useful when you want to create mappings for types such as sql.NullString
	// without traversing within the sql.NullString itself.
	TreatAsScalar TypeList
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
	//
	// The second argument, if non-nil, will be the first return value.  In other words pre-allocating
	// rv will cut down on memory allocations.  It is assumed that if rv is non-nil that len(fields) == len(rv)
	// and the lengths are not checked, meaning your program could panic.
	//
	// During traversal this method will instantiate struct fields that are themselves structs.
	//
	// An example use-case would be obtaining a slice of pointers for Rows.Scan() during database
	// query results.
	Assignables(fields []string, rv []interface{}) ([]interface{}, error)
	// Copy creates an exact copy of the BoundMapping.  Since a BoundMapping is implicitly tied to a single
	// value sometimes it may be useful to create and cache a BoundMapping early in a program's initialization and
	// then call Copy() to work with multiple instances of bound values simultaneously.
	Copy() BoundMapping
	// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
	// calls to Rebind()
	Err() error
	// Field returns the *Value for field.
	Field(field string) (*Value, error)
	// Fields returns a slice of interfaces by field names where each element is the field value.
	//
	// The second argument, if non-nil, will be the first return value.  In other words pre-allocating
	// rv will cut down on memory allocations.  It is assumed that if rv is non-nil that len(fields) == len(rv)
	// and the lengths are not checked, meaning your program could panic.
	//
	// During traversal this method will instantiate struct fields that are themselves structs.
	//
	// An example use-case would be obtaining a slice of query arguments by column name during
	// database queries.
	Fields(fields []string, rv []interface{}) ([]interface{}, error)
	// Rebind will replace the currently bound value with the new variable I.  If the underlying types do
	// not match a panic will occur.
	Rebind(I interface{})
	// Set effectively sets V[field] = value.
	Set(field string, value interface{}) error
}

// DefaultMapper joins names by "_" but performs no other modifications.
var DefaultMapper = &Mapper{
	Join: "_",
}

// Bind creates a Mapping bound to a specific instance I of a variable.
func (me *Mapper) Bind(I interface{}) BoundMapping {
	var v *Value
	if tv, ok := I.(*Value); ok {
		v = tv
	} else {
		v = V(I)
	}
	rv := newBoundMapping(v, me.Map(v))
	return rv
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
// Mappings that are returned are shared resources and should not be altered in any way.  If this is your
// use-case then create a copy of the Mapping with Mapping.Copy.
func (me *Mapper) Map(T interface{}) *Mapping {
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
		return rv.(*Mapping)
	}
	//
	rv := &Mapping{
		Keys:         []string{},
		Indeces:      map[string][]int{},
		StructFields: map[string]reflect.StructField{},
		// CacheInfo: map[string]struct { // TODO
		// 	Indeces  []int
		// 	SharedId int
		// }{},
	}
	//
	var scan func(typeInfo TypeInfo, indeces []int, prefix string)
	scan = func(typeInfo TypeInfo, indeces []int, prefix string) {
		for k, field := range typeInfo.StructFields {
			if field.PkgPath != "" {
				continue
			}
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
			if _, ok := mapperTreatAsScalar[fieldTypeInfo.Type]; ok {
				rv.Keys, rv.Indeces[name], rv.StructFields[name] = append(rv.Keys, name), nameIndeces, field
			} else if _, ok = me.TreatAsScalar[fieldTypeInfo.Type]; ok {
				rv.Keys, rv.Indeces[name], rv.StructFields[name] = append(rv.Keys, name), nameIndeces, field
			} else if fieldTypeInfo.IsStruct {
				scan(fieldTypeInfo, nameIndeces, name)
			} else if fieldTypeInfo.IsScalar {
				rv.Keys, rv.Indeces[name], rv.StructFields[name] = append(rv.Keys, name), nameIndeces, field
			}
		}
	}
	// Scan and assign the result to our known types.
	scan(typeInfo, []int{}, "")
	me.known.Store(typeInfo.Type, rv)
	//
	return rv
}

// Copy creates a copy of the Mapping.
func (me *Mapping) Copy() *Mapping {
	rv := &Mapping{
		Keys:         append([]string{}, me.Keys...),
		Indeces:      map[string][]int{},
		StructFields: map[string]reflect.StructField{},
		// CacheInfo: map[string]struct { // TODO
		// 	Indeces  []int
		// 	SharedId int
		// }{},
	}
	for _, key := range me.Keys {
		rv.Indeces[key] = append([]int{}, me.Indeces[key]...)
		rv.StructFields[key] = me.StructFields[key]
		// TODO Copy CacheInfo as well.
	}
	return rv
}

// Get returns the indeces associated with key in the mapping.  If no such key
// is found a nil slice is returned.
func (me *Mapping) Get(key string) []int {
	v, _ := me.Lookup(key)
	return v
}

// Lookup returns the value associated with key in the mapping.  If no such key is
// found a nil slice is returned and ok is false; otherwise ok is true.
func (me *Mapping) Lookup(key string) (indeces []int, ok bool) {
	if me == nil || me.Indeces == nil {
		return nil, false
	}
	indeces, ok = me.Indeces[key]
	return indeces, ok
}

// String returns the Mapping as a string value.
func (me *Mapping) String() string {
	if me == nil {
		return ""
	}
	parts := []string{}
	for str, indeces := range me.Indeces {
		parts = append(parts, fmt.Sprintf("%v\t\t%v", indeces, str))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}

// boundMapping is the implementation for BoundMapping.
type boundMapping struct {
	value *Value
	err   error
	//
	// NB: Fields below here are read-only.  boundMapping does not alter them in any way
	// thus during Copy() they do not need to be copied.  This is an effort to reduce
	// memory allocations.
	mapping *Mapping
}

// newBoundMapping creates a new boundMapping type.
func newBoundMapping(value *Value, mapping *Mapping) *boundMapping {
	return &boundMapping{
		value:   value,
		mapping: mapping,
	}
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
func (me *boundMapping) Assignables(fields []string, rv []interface{}) ([]interface{}, error) {
	// nil check is not necessary as boundMapping is only created within this package.
	if !me.value.CanWrite {
		return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldnum, name := range fields {
		indeces := me.mapping.Get(name)
		if len(indeces) == 0 {
			return nil, errors.Errorf("No mapping for field [%v]", name)
		}
		v := me.value.WriteValue
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
func (me *boundMapping) Copy() BoundMapping {
	return &boundMapping{
		value: me.value.Copy(),
		err:   me.err,
		// NB: Don't need to copy me.mapping since we never alter it.
		// mapping: me.mapping.Copy(),
		mapping: me.mapping,
	}
}

// Err returns an error that may have occurred during repeated calls to Set(); it is reset on
// calls to Rebind()
func (me *boundMapping) Err() error {
	// nil check is not necessary as boundMapping is only created within this package.
	return me.err
}

// Field returns the *Value for field.
func (me *boundMapping) Field(field string) (*Value, error) {
	// nil check is not necessary as boundMapping is only created within this package.
	var v reflect.Value
	var err error
	if v, err = me.value.FieldByIndex(me.mapping.Get(field)); err != nil {
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
func (me *boundMapping) Fields(fields []string, rv []interface{}) ([]interface{}, error) {
	// nil check is not necessary as boundMapping is only created within this package.
	if !me.value.CanWrite {
		return nil, errors.Errorf("Value in BoundMapping is not writable; pass the address of your destination value when binding.")
	}
	if rv == nil {
		rv = make([]interface{}, len(fields))
	}
	for fieldnum, name := range fields {
		indeces := me.mapping.Get(name)
		if len(indeces) == 0 {
			return nil, errors.Errorf("No mapping for field [%v]", name)
		}
		v := me.value.WriteValue
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
func (me *boundMapping) Rebind(I interface{}) {
	// nil check is not necessary as boundMapping is only created within this package.
	me.err = nil
	me.value.Rebind(I)
}

// Set effectively sets V[field] = value.
func (me *boundMapping) Set(field string, value interface{}) error {
	// nil check is not necessary as boundMapping is only created within this package.
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
