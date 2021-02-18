// Package set is a small wrapper around the official reflect package that facilitates loose type conversion,
// assignment into native Go types, and utilities to populate deeply nested Go structs.
//
// Data Types
//
// In this package documentation float, int, and uint include the bit-specific types.  For example
// float includes float32 and float64, int includes int8, int16, int32, and int64, etc.
//
// Scalar Types
//
// This package considers the following types to be scalars:
//	bool, float, int, uint, & string
//	pointer to any of the above; e.g: *bool, *float, *int, etc
//	nested pointer to any of the above; e.g: ***bool, *****string, etc
//
// Package Name Collision
//
// I called this package `set` because I like typing short identifiers and I don't really use
// packages that implement logical set operations such as union or intersection.
//
// I also like the semantics of typing: set.V().To() // i.e. set value to
//
// If you find yourself dealing with name collision here are some alternate imports
// that are still short and keep the semantics (with varying success):
//	import (
//		assign "github.com/nofeaturesonlybugs/set"
//		accept "github.com/nofeaturesonlybugs/set"
//		coerce "github.com/nofeaturesonlybugs/set"
//		from "github.com/nofeaturesonlybugs/set"
//		make "github.com/nofeaturesonlybugs/set"
//		pin "github.com/nofeaturesonlybugs/set"
//		will "github.com/nofeaturesonlybugs/set"
//	)
//
// Basic Type Coercion
//
// A simple example with type coercion:
// 	b, i := true, 42
//	set.V(&b).To("False")		// Sets b to false
//	set.V(&i).To("3.14")		// Sets i to 3
//
// In general:
//	var t T					// T is a target data type, t is a variable of that type.
//	var s S					// S is a source data, s is a variable of that type.
//	set.V(&t).To(s)				// Sets s into t with a "best effort" approach.
//
// If t is not a pointer or is a pointer that is set to nil then pass the address of t:
//	var t bool		// Pass the address when t is not a pointer.
//	set.V(&t)
//
//	var t *int 		// Pass the address when t is a pointer but set to nil.
//	set.V(&t)
//
//	var i int
//	t := &i			// Do not pass the address when t is a pointer and not set to nil.
//	set.V(t)		// Note that values assigned go into i!
//
// Scalar to Scalar of Same Type
//
// When T and S are both scalars of the same type then assignment with this package is the same as direct
// assignment:
//	var a int			// Call this T
//	var b int			// Call this S
//	set.V(&a).To(b)			// Effectively the same as: a = b
//
//	// pointers to structs
//	x := &T{}				// Call this T
//	y := &T{"data", "more data"}		// Call this S
//	set.V(&x).To(y)				// Effectively the same as: x = y
//
// This package grants no benefit when used as such.
//
// Scalar to Scalar of Different Type
//
// When T and S are both scalars but with different types then assignments can be made with type coercion:
//	a := int(0)			// Call this T
//	b := uint(42)			// Call this S
//	set.V(&a).To(b)			// This coerces b into a if possible.
//	set.V(&a).To("-57")		// Also works.
//	set.V(&a).To("Hello")		// Returns an error.
//
// Pointers Are Tricky But Work Well
//
// If a pointer already contains a memory address then you do not need to pass the pointer's address to set:
//	var b bool
//	bp := &b 		// bp already contains an address
//	set.V(bp).To("1")	// b is set to true
//
// If a pointer does not contain an address then passing the pointer's address allows set to create the
// data type and assign it to the pointer:
//	var bp *bool
//	set.V(&bp).To("True")	// bp is now a pointer to bool and *bp is true.
//
// This works even if the pointer is multiple levels of indirection:
//	var bppp ***bool
//	set.V(&bppp).To("True")
//	fmt.Println(***bppp) // Prints true
//
// If S is a pointer it will be dereferenced until the final value; if the final value is a scalar
// then it will be coerced into T:
//	var ippp ***int
//	s := "42"
//	sp := &s
//	spp := &sp
//	set.V(&ippp).To(spp)
//	fmt.Println(***ippp) // Prints 42
//
//
// Scalars to Slices
//
// When T is a slice and S is a scalar then T is assigned a slice with S as its single element:
//	var b []bool
//	set.V(&b).To("True") // b is []bool{ true }
//
// Slices to Scalars
//
// When T is a scalar and S is a slice then the last element of S is assigned to T:
//	var b bool
//	set.V(&b).To([]bool{ false, false, true } ) // b is true, coercion not needed.
//	set.V(&b).To([]interface{}{ float32(1), uint(0) }) // b is false, coercion needed.
//
// If S is a nil or empty slice then T is set to an appropriate zero value.
//
// Slices to Slices
//
// When T and S are both slices set always creates a new slice []T and copies elements from []S into []T.
//	var t []bool
//	var s []interface{}
//	s = []interface{}{ "true", 0, float64(1) }
//	set.V(&t).To(s) // t is []bool{ true, false, true }
//
//	var t []bool
//	var s []bool
//	s = []bool{ true, false, true }
//	set.V(&t).To(s) // t is []bool{ true, false, true } and t != s
//
// If a single element within []S can not be coerced into an element of T then []T will be empty:
//	var t []int
// 	var s []string
//	s = []string{ "42", "24", "Hello!" }
//	set.V(&t).To(s) // t is []int{} because "Hello" can not coerce.
//
//
// Populating Structs with Value.Fill() and a Getter
//
// Structs can be populated by using Value.Fill() and a Getter; note the function is type casted to
// a set.GetterFunc.
//	// An example getter.
//	myGetter := set.GetterFunc(func( name string ) interface{} {
// 		switch name {
// 		case "Name":
// 			return "Bob"
// 		case "Age":
// 			return "42"
// 		default:
// 			return nil
//		}
//	})
//
// Populating a struct by field; i.e. the struct field names are the names passed to the Getter:
// 	type T struct {
// 		Name string
// 		Age uint
// 	}
// 	var t T
// 	set.V(&t).Fill(myGetter)
//
// Populating a struct by struct tag; i.e. if the struct tag exists on the field then the tag's value
// is passed to the Getter:
// 	type T struct {
// 		SomeField string `key:"Name"`
// 		OtherField uint `key:"Age"`
// 	}
// 	var t T
// 	set.V(&t).FillByTag("key", myGetter)
//
// Populating Nested Structs with Value.Fill() and a Getter
//
// To populate nested structs a Getter needs to return a Getter for the given name:
// 	myGetter := set.GetterFunc(func(key string) interface{} {
// 		switch key {
// 		case "name":
// 			return "Bob"
// 		case "age":
// 			return "42"
// 		case "address":
// 			return set.GetterFunc(func(key string) interface{} {
// 				switch key {
// 				case "street1":
// 					return "97531 Some Street"
// 				case "street2":
// 					return ""
// 				case "city":
// 					return "Big City"
// 				case "state":
// 					return "ST"
// 				case "zip":
// 					return "12345"
// 				default:
// 					return nil
// 				}
// 			})
// 		default:
// 			return nil
// 		}
// 	})
//
// Using the Getter above:
// 	type Address struct {
// 		Street1 string `key:"street1"`
// 		Street2 string `key:"street2"`
// 		City    string `key:"city"`
// 		State   string `key:"state"`
// 		Zip     string `key:"zip"`
// 	}
// 	type Person struct {
// 		Name    string  `key:"name"`
// 		Age     uint    `key:"age"`
// 		Address Address `key:"address"`
// 	}
// 	var t Person
// 	set.V(&t).FillByTag("key", myGetter)
//
// Maps as Getters
//
// A more practical source of data for a Getter might be a map.  To use a map as a Getter the map key has
// to be assignable to string; e.g: string or interface{}.  If the map contains nested maps that also meet
// the criteria for becoming a Getter then those maps can be used to populate nested structs.
//
// 	m := map[string]interface{}{
// 		"name": "Bob",
// 		"age":  42,
// 		"address": map[interface{}]string{
// 			"street1": "97531 Some Street",
// 			"street2": "",
// 			"city":    "Big City",
// 			"state":   "ST",
// 			"zip":     "12345",
// 		},
// 	}
// 	myGetter := set.MapGetter(m)
//
// 	type Address struct {
// 		Street1 string `key:"street1"`
// 		Street2 string `key:"street2"`
// 		City    string `key:"city"`
// 		State   string `key:"state"`
// 		Zip     string `key:"zip"`
// 	}
// 	type Person struct {
// 		Name    string  `key:"name"`
// 		Age     uint    `key:"age"`
// 		Address Address `key:"address"`
// 	}
// 	var t Person
// 	set.V(&t).FillByTag("key", myGetter)
//
// Populating Structs with Mapper, Mapping, and BoundMap
//
// If you need to populate or traverse structs using strings as lookups consider using a Mapper.  A Mapper traverses a type T
// and generates a Mapping, where a Mapping is currently implemented as a map[string][]int.
//
// When you index into a Mapping you will receive a slice of ints representing the indeces into the nested structure
// to the desired field.
//
// For convenience a Mapper can create a BoundMapping which binds the Mapping to an instance of T.  The BoundMapping
// can then be used to update the data within the instance.  See the BoundMapping examples.
//
// Examples Subdirectory
//
// The examples subdirectory contains multiple examples for this package; separating them keeps
// this package documentation a little cleaner.
//	var myVar bool
//	value := set.V(&myVar)
// To see what you can do with `value` in the above code find the `Bool` type in the examples package.
//	// Assuming SomeType is a struct.
//	var myVar SomeType
//	value := set.V(&myVar)
// To see what you can do with `value` in the above code find the `Struct` type in the examples package.
//
package set
