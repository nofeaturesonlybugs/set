package path

import (
	"fmt"
	"strings"
)

// PathIndeces is a [][]int, where the inner []int would be appropriate
// for passing to reflect.Value.FieldByIndex([]int) to select a nested
// field.
//
// A caveat of reflect.Value.FieldByIndex() is that it doesn't automatically
// instantiate intermediate fields that may be pointers.  Any field that is
// a pointer effectively breaks or segments the original []int into a pair of []int
// with one on each side of the pointer.  As a result this package stores the PathIndeces
// as [][]int.
//
// If len(PathIndeces)==1 then the pathway is contiguous as described in the
// documentation for PathOffsetSegment and a single call of
// reflect.Value.FieldByIndex(Indeces[0]) is enough to obtain the field.
type PathIndeces [][]int

// String returns the string description.
func (i PathIndeces) String() string {
	n := len(i)
	if n == 0 {
		return "<empty>"
	}
	parts := make([]string, n)
	for k, slice := range i {
		parts[k] = fmt.Sprintf("%v", slice)
	}
	return strings.Join(parts, " ∪ ")
}
