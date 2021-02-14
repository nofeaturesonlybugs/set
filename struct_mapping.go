package set

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// StructMapping traverses nested structures and generates a mapping of Key-to-Indeces where:
//	Key is a string key representing a common or friendly name for the nested struct member.
//	Indeces is an []int that can be used to index to the proper struct member.
type StructMapping struct {
	mapping map[string][]int
}

type StructMappingOptions struct {
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
}

var DefaultStructMappingOptions = &StructMappingOptions{
	Join: "_",
}

func NewStructMapping(v *Value, options *StructMappingOptions) *StructMapping {
	if options == nil {
		options = DefaultStructMappingOptions
	}
	rv := &StructMapping{
		mapping: map[string][]int{},
	}
	//
	ignored := map[reflect.Type]struct{}{}
	elevated := map[reflect.Type]struct{}{}
	for _, value := range options.Ignored {
		ignored[reflect.TypeOf(value)] = struct{}{}
	}
	for _, value := range options.Elevated {
		elevated[reflect.TypeOf(value)] = struct{}{}
	}
	//
	var scan func(v *Value, indeces []int, prefix string, indent int)
	scan = func(v *Value, indeces []int, prefix string, indent int) {
		for k, field := range v.Fields() {
			if _, found := ignored[field.Value.pt]; found {
				continue
			}
			//
			name := ""
			if _, found := elevated[field.Value.pt]; !found {
				for _, tagName := range append(options.Tags, "") {
					if tagValue, ok := field.Field.Tag.Lookup(tagName); ok {
						name = tagValue
						break
					} else if tagName == "" {
						name = field.Field.Name
						if options.Transform != nil {
							name = options.Transform(name)
						}
						break
					}
				}
			}
			if prefix != "" && name != "" {
				name = prefix + options.Join + name
			} else if prefix != "" {
				name = prefix
			}
			nameIndeces := append(indeces, k)
			if field.Value.IsStruct {
				scan(field.Value, nameIndeces, name, indent+1)
			} else if field.Value.IsScalar {
				rv.mapping[name] = nameIndeces
			}
		}
	}
	scan(v, []int{}, "", 0)
	return rv
}

// Get returns the indeces associated with key in the mapping.  If no such key
// is found a nil slice is returned.
func (me *StructMapping) Get(key string) []int {
	v, _ := me.Lookup(key)
	return v
}

// Lookup returns the value associated with key in the mapping.  If no such key is
// found a nil slice is returned and ok is false; otherwise ok is true.
func (me *StructMapping) Lookup(key string) (indeces []int, ok bool) {
	if me == nil {
		return nil, false
	}
	indeces, ok = me.mapping[key]
	return indeces, ok
}

// String returns the StructMapping as a string value.
func (me *StructMapping) String() string {
	parts := []string{}
	for str, indeces := range me.mapping {
		parts = append(parts, fmt.Sprintf("%v\t\t%v", indeces, str))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}
