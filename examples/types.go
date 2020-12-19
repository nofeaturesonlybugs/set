package examples

//
// SCALARS
//

// Bool shows how to use set.V(&boolType).
type Bool struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a bool.
func (b Bool) Set() {}

// Float shows how to use set.V(&floatType)
type Float struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a float.
func (f Float) Set() {}

// Int shows how to use set.V(&intType)
type Int struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into an int.
func (i Int) Set() {}

// Uint shows how to use set.V(&uintType)
type Uint struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a uint.
func (i Uint) Set() {}

// String shows how to use set.V(&stringType)
type String struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a string.
func (s String) Set() {}

// Scalar shows how to use set.V(&scalarType)
type Scalar struct{}

// Set is a stub.  The proceeding examples demonstrate setting non-scalar values on scalar destinations.
func (s Scalar) Set() {}

// Struct shows how to use set.V(&structType)
type Struct struct{}

// Fill is a stub.  The proceeding examples demonstrate calling Fill on a struct.
func (s Struct) Fill() {}

//
// SLICES
//

// BoolSlice shows how to use set.V(&sliceBoolType)
type BoolSlice struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a []bool.
func (b BoolSlice) Set() {}

// FloatSlice shows how to use set.V(&sliceFloatType)
type FloatSlice struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a []float.
func (f FloatSlice) Set() {}

// IntSlice shows how to use set.V(&sliceIntType)
type IntSlice struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into an []int.
func (i IntSlice) Set() {}

// UintSlice shows how to use set.V(&sliceUintType)
type UintSlice struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a []uint.
func (i UintSlice) Set() {}

// StringSlice shows how to use set.V(&sliceStringType)
type StringSlice struct{}

// Set is a stub.  The proceeding examples demonstrate coercing types into a []string.
func (s StringSlice) Set() {}
