package path

import (
	"fmt"
	"reflect"
	"strings"
)

// PathOffsetSegment describes a segment of a Path when traversed via pointer arithmetic.
type PathOffsetSegment struct {
	// Offset is the total offset from the memory address of the
	// struct at the top of the struct hierarchy for this specific
	// segment.
	//
	//                      Address
	//	A struct            0xabab                # Top
	//	    B struct        0xabab + Offset(B)    # A is top.
	//	        C struct    0xabab + Offset(C)    # A is top.
	//
	// A and its members are contiguous regions of memory.  All Offsets
	// are calculated from A's memory address.
	//
	//                      Address
	//	X struct            0xeded                                     # Top
	//	    *Y struct       0xeded + Offset(Y) -> 0xeeee               # A is top; pointer begins new hierarchy at Y.
	//	        M                                 0xeeee + Offset(M)   #   Y is top.
	//	        N                                 0xeeee + Offset(N)   #   Y is top.
	//
	// X's hierarchy is segmented by the pointer at *Y.  Therefore Y's Offset
	// is calculated from X's memory address **but** the Offsets for
	// M and N are added to Y's address to obtain their locations.
	Offset uintptr

	// IndirectionLevel determines if the field at Offset is a pointer or not.
	//
	// IndirectionLevel=0 means field is not a pointer.
	// IndirectionLevel=1 means field is a pointer with a single level of indirection.
	// IndirectionLevel>2 means field is a pointer with multiple levels of indirection.
	//
	//                  IndirectionLevel
	//	A struct                0        *(&A + Offset) -> is not a pointer
	//      B struct            0        *(&A + Offset) -> is not a pointer
	//          C struct        0        *(&A + Offset) -> is not a pointer.
	//
	//	X struct                0        *(&A + Offset) -> is not a pointer.
	//      *Y struct           1        *(&A + Offset) -> is *Y, pointer with single level of indirection.
	//     **Z struct           2        *(&A + Offset) -> is **Z, pointer with 2 levels of indirection.
	IndirectionLevel int

	// Type and EndType describe the type(s) located at Offset.
	//
	// IndirectionLevel=0 means Type and EndType are the same type.
	// IndirectionLevel>0 means Type is the pointer's type and EndType
	// is the type at the end of the pointer chain.
	//
	//                          Type     EndType
	//	A struct                A        A
	//      B struct            B        B
	//          C struct        C        C
	//
	//	X struct                X        X
	//      *Y struct          *Y        Y
	//     **Z struct         **Z        Z
	Type    reflect.Type
	EndType reflect.Type
}

// String returns the string description.
func (s PathOffsetSegment) String() string {
	if s.IndirectionLevel == 0 {
		return fmt.Sprintf("+%v ∴ %v", s.Offset, s.Type)
	}
	return fmt.Sprintf("+%v ↬ %v ∴ %v", s.Offset, s.IndirectionLevel, s.EndType)
}

// PathOffsets is a slice of PathOffsetSegment.
type PathOffsets []PathOffsetSegment

// String returns the string description.
func (o PathOffsets) String() string {
	n := len(o)
	if n == 0 {
		return "<empty>"
	}
	parts := make([]string, n)
	for k, part := range o {
		parts[k] = part.String()
	}
	return strings.Join(parts, " ∪ ")
}

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

// A Path represents a traversal from an originating struct to an accessible
// public field within the struct or any of its descendents, where the descendents
// are limited to being structs, pointer-to-struct, or pointer chains that end
// in struct.
//
// Paths are either contiguous or segmented.  A contiguous path is one where
// the origin and destination occupy a single unbroken region of memory.  A segmented
// path is one where the origin and destination may be in separate regions of
// memory; segmented paths are caused when intermediary fields along the path are
// pointers.
//  A struct
//      B struct
//          C struct
//              M, N
// The Paths A ↣ M (start at A and traverse to M) or A ↣ N are contiguous
// because no intermediary field between the origin and destination is a pointer.
// In terms of pointer arithmetic M or N can be reached by adding their total offset
// to A's memory address.
//
//  A struct
//      B struct
//         *C struct
//              M, N
// The Paths A ↣ B and A ↣ C are still contiguous because no section of the path
// is segmented by a pointer.  Note that even though C itself is a pointer
// it can be reached by adding its total offset to A's memory address.
//
// However the paths A ↣ M and A ↣ N are no longer contiguous.  In order to
// start at A and reach either M or N we must first traverse to C, perform
// pointer dereferencing (and possibly allocate a C), and then from &(*C)
// we can reach M or N by adding their offsets to C's address.
type Path struct {
	// Name is the field's name.
	Name string

	// Index is the field's index into its owning struct.
	Index int

	// Offset is the field's uintptr offset into its owning struct.
	Offset uintptr

	// Type is the field's reflect.Type.
	Type reflect.Type

	// TODO+NB Information describing the pathway to this struct field.
	PathwayOffsets []PathOffsetSegment // TODO+NB Describe me.
	PathwayIndex   [][]int             // Slice of indeces leading to this path; appropriate for reflect.Value.FieldByIndex()
	PathwayName    string              // The joined named of the pathway.

	// ParentPathwayName is the pathway name of the parent.
	ParentPathwayName string
}

// String returns Path represented as a string.
func (p Path) String() string {
	return fmt.Sprintf("%v %v %v Type=%v Pathway[%v]%v Parent[%v] Offsets= %v", p.Name, p.Index, p.Offset, p.Type, p.PathwayName, PathIndeces(p.PathwayIndex), p.ParentPathwayName, PathOffsets(p.PathwayOffsets))
}

// Value returns the reflect.Value for the path from the origin value.
func (p Path) Value(origin reflect.Value) reflect.Value {
	for ; origin.Kind() == reflect.Ptr; origin = origin.Elem() {
		// Walk pointer chain if origin is a pointer
	}
	for _, index := range p.PathwayIndex {
		origin = origin.FieldByIndex(index)
		for ; origin.Kind() == reflect.Ptr; origin = origin.Elem() {
			if origin.IsNil() && origin.CanSet() {
				origin.Set(reflect.New(origin.Type().Elem()))
			}
		}
	}
	return origin
}
