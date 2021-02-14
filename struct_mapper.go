package set

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/nofeaturesonlybugs/errors"
)

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

var DefaultMapper = &Mapper{
	Join: "_",
}

// Map adds T to the StructMapper's list of known and recognized types.
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
	if rv, ok := me.known[v.pt]; ok {
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
			if me.Ignored.Has(field.Value.pt) {
				continue
			}
			//
			name := ""
			if !me.Elevated.Has(field.Value.pt) {
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
			if field.Value.IsStruct {
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
	me.known[v.pt] = rv
	//
	return rv, nil
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

// String returns the StructMapping as a string value.
func (me Mapping) String() string {
	parts := []string{}
	for str, indeces := range me {
		parts = append(parts, fmt.Sprintf("%v\t\t%v", indeces, str))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}
