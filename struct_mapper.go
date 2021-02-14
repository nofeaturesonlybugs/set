package set

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// Mapping is the result of traversing nested structures to generate a mapping of Key-to-Indeces where:
//	Key is a string key representing a common or friendly name for the nested struct member.
//	Indeces is an []int that can be used to index to the proper struct member.
type Mapping map[string][]int

// Mapper is used to traverse structures to create Mappings and then navigate the nested
// structure using string keys.
type Mapper struct {
	// A slice of type instances that should be ignored during name generation.
	Ignored []interface{}
	//
	// During name generation types in the `Elevated` member will not have the name affected
	// by the struct field name or data type.  Use this for struct members or embedded structs
	// when you do not want their name or type affecting the generated name.
	Elevated []interface{}
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
	mut      sync.RWMutex
	ready    bool
	elevated map[reflect.Type]struct{} // Types that are elevated.
	ignored  map[reflect.Type]struct{} // Types that are ignored.
	known    map[reflect.Type]Mapping  // Types that are known -- i.e. we've already created the mapping.
}

var DefaultStructMapper = &Mapper{
	Join: "_",
}

// init initializes some internal members.
func (me *Mapper) init() {
	if me == nil {
		return
	}
	me.mut.RLock()
	if me.ready {
		me.mut.RUnlock()
		return
	} else {
		me.mut.RUnlock()
		me.mut.Lock()
		defer me.mut.Unlock()
		//
		me.ignored = map[reflect.Type]struct{}{}
		me.elevated = map[reflect.Type]struct{}{}
		me.known = make(map[reflect.Type]Mapping)
		for _, value := range me.Ignored {
			me.ignored[reflect.TypeOf(value)] = struct{}{}
		}
		for _, value := range me.Elevated {
			me.elevated[reflect.TypeOf(value)] = struct{}{}
		}
	}
}

// Register adds T to the StructMapper's list of known and recognized types.
func (me *Mapper) Register(T interface{}) Mapping {
	me.init()
	rv := make(Mapping)
	if me == nil {
		return rv
	}
	//
	v := V(T)
	//
	// Before we scan we can check if this is a known type; if so just return our previous scan.
	me.mut.RLock()
	if rv, ok := me.known[v.pt]; ok {
		me.mut.RUnlock()
		return rv
	}
	me.mut.RUnlock() // If we make it here it is not known; release the lock so we can acquire Write lock.
	//
	var scan func(v *Value, indeces []int, prefix string, indent int)
	scan = func(v *Value, indeces []int, prefix string, indent int) {
		for k, field := range v.Fields() {
			if _, found := me.ignored[field.Value.pt]; found {
				continue
			}
			//
			name := ""
			if _, found := me.elevated[field.Value.pt]; !found {
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
	// Our scan is complete so now we should assign the result to our known types.
	me.mut.Lock()
	defer me.mut.Unlock()
	me.known[v.pt] = rv
	//
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

// String returns the StructMapping as a string value.
func (me Mapping) String() string {
	parts := []string{}
	for str, indeces := range me {
		parts = append(parts, fmt.Sprintf("%v\t\t%v", indeces, str))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}
