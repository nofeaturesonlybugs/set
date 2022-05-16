// Package coerce provides loose type coercion and assignment into native Go types.
//
// Primitive Types
//
// This package considers the following types to be primitive types:
//	bool
//	float32 float64
//	int8 int16 int32 int64 int
//	uint8 uint16 uint32 uint64 uint
//	string
//
// Coercion Logic
//
// All coercion functions use the same basic logic to coerce an incoming value v:
//	for {
//		v is primitive
//			return coerced value or error
//		v's underlying kind is primitive
//			// convert v to underlying kind and try again
//			v = UnderlyingKind(v); continue
//		v is a pointer
//			// dereference v and try again
//			v = *v; continue
//		v is a slice
//			// try again with last slice element
//			v = v[len(v)-1]; continue
//		ErrUnsupported
//	}
//
// If v is a primitive type it will be coerced via type coercion, one or more calls to the
// strconv package, or a small amount of custom logic.
//
// If v's underlying kind is primitive it is converted to the underlying primitive and the loop
// restarts with a continue statement.
//
// If v is a pointer or any pointer chain it is followed until the final value and the loop
// restarts with a continue statement.  A nil pointer shortcuts and returns an appropriate zero value.
//
// If v is a slice then v is reassigned to the last element and the loop starts again with
// a continue statement.  An empty slice shortcuts and returns an appropriate zero value.
//
// All other types for v (e.g. chan, map, func, etc) return a zero value and ErrUnsupported.
//
// Overflow
//
// During numeric coercions this package checks incoming values against the minimum and maximum value
// for the target type.  If the incoming value is out of range for the target type ErrOverflow is returned.
//
// Otherwise the coercion is made with type conversion.
//
// String Parsing
//
// Where necessary this package will call into strconv to parse values for bool, int, float, or uint.
//
// If a call to strconv returns error this package will return the zero value for the target type and
// either ErrOverflow or ErrInvalid depending on the error received from strconv.
//
// This package may make multiple calls to strconv to parse an incoming string.  For example it may first
// try strconv.ParseInt, followed by strconv.ParseUint, and finally strconv.ParseFloat to parse a string.
// If any of the calls succeed then type coercion continues with the parsed value.
package coerce
