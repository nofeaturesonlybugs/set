package set

import (
	"reflect"
)

// CanPanic is a namespace for operations prioritizing speed over type safety or error checking.  Reach for this
// namespace when your usage of the `set` package is carefully crafted to ensure panics will not result from your
// actions.
//
// Methods within CanPanic will not validate that points are non-nil.
//
// It is strongly encouraged to create suitable `go tests` within your project when reaching for CanPanic.
//
// You do not need to create or instantiate this type; instead you can use the global `var Panics`.
type CanPanic struct{}

// Panics is a global instance of CanPanic; it is provided for convenience.
var Panics = CanPanic{}

// Append appends any number of Value types to the dest Value.  The safest way to use this method
// is when your code uses a variant of:
//	var records []*Record
//	v := set.V( &records )
//	for k := 0; k < 10; k++ {
//		elem := v.Elem.New()
//		set.Panics.Append( v, elem )
//	}
func (p CanPanic) Append(dest Value, values ...Value) {
	for _, value := range values {
		dest.WriteValue.Set(reflect.Append(dest.WriteValue, reflect.Indirect(value.TopValue)))
	}
}
