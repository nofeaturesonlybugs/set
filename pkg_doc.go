// Package set is a performant reflect wrapper supporting loose type conversion,
// struct mapping and population, and slice building.
//
// Type Coercion
//
// Value and its methods provide a generous facility for type coercion:
//	var t T					// T is a target data type, t is a variable of that type.
//	var s S					// S is a source data type, s is a variable of that type.
//	set.V(&t).To(s)				// Sets s into t with a "best effort" approach.
//
// See documentation and examples for Value.To.
//
// The examples subdirectory contains additional examples for the Value type.
//
// Finally you may wish to work directly with the coerce subpackage, which is the workhorse
// underneath Value.To.
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
// What is Rebind
//
// The BoundMapping, PreparedMapping, and Value types internally contain meta data about the
// types they are working with.  Most of this meta data is obtained with calls to reflect and
// calls to reflect can be expensive.
//
// In one-off scenarios the overhead of gathering meta data is generally not a concern.  But in
// tight-loop situations this overhead begins to add up and is a key reason reflect has gained
// a reputation for being slow in the Go community.
//
// Where appropriate types in this package have a Rebind method.  Rebind swaps the current value
// being worked on with a new incoming value without regathering reflect meta data.  When
// used appropriately with Rebind the BoundMapping, PreparedMapping, and Value types become
// much more performant.
//
// A Note About Package Examples
//
// Several examples ignore errors for brevity:
//	_ = p.Plan(...) // error ignored for brevity
//	_ = b.Set(...) // error ignored for brevity
//
// This is a conscious decision because error checking is not the point of the examples.  However
// in production code you should check errors appropriately.
package set
