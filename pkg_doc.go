// Package set is a performant reflect wrapper supporting loose type conversion,
// struct mapping and population, and slice building.
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
// Populating Structs by Lookup Function
//
// See examples for GetterFunc.
//
// Populating Structs by Map
//
// See example for MapGetter.
//
// Struct Mapping
//
// Struct and struct hierarchies can be mapped to a flat list of string keys.  This is useful
// for deserializers and unmarshalers that need to convert a friendly string such as a column name or
// environment variable and use it to locate a target field within a struct or its hierarchy.
//
// See examples for Mapper.
//
// Mapping, BoundMapping, and PreparedMapping
//
// Once an instance of Mapper is created it can be used to create Mapping, BoundMapping, and
// PreparedMapping instances that facilitate struct traversal and population.
//
// BoundMapping and PreparedMapping are specialized types that bind to an instance of T
// and allow performant access to T's fields or values.
//
// See examples for Mapper.Bind and Mapper.Prepare.
//
// In tight-loop scenarios an instance of BoundMapping or PreparedMapping can be bound
// to a new instance T with the Rebind method.
//
// See examples for BoundMapping.Rebind and PreparedMapping.Rebind.
//
// If neither BoundMapping nor PreparedMapping are suitable for your use case you can
// call Mapper.Map to get a general collection of data in a Mapping.  The data in a Mapping
// may be helpful for creating your own traversal or population algorithm without having
// to dive into all the complexities of the reflect package.
//
// BoundMapping vs PreparedMapping
//
// BoundMapping allows adhoc access to the bound struct T.  You can set or retrieve fields in
// any order.  Conceptually a BoundMapping is similar to casting a struct and its hierarchy into
// a map of string keys that target fields within the hierarchy.
//
// PreparedMapping requires an access plan to be set by calling the Plan method.  Once set
// the bound value's fields must be set or accessed in the order described by the plan.  A
// PreparedMapping is similar to a prepared SQL statement.
//
// Of the two PreparedMapping yields better performance.  You should use PreparedMapping
// when you know every bound value will have its fields accessed in a determinate order.  If
// fields will not be accessed in a determinate order then you should use a BoundMapping.
//
// BoundMapping methods require the field(s) as arguments; in some ways
// this can help with readability as your code will read:
//	b.Set("FooField", "Hello")
//	b.Set("Number", 100)
//
// whereas code using PreparedMapping will read:
//	p.Set("Hello")
//	p.Set(100)
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
