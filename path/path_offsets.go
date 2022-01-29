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
