package path

import (
	"fmt"
	"reflect"
)

// Paths is a slice of Path.
type Paths []Path

// Len returns the length of Paths.
func (p Paths) Len() int {
	return len(p)
}

// Less returns true if the value at a is less than b.
func (p Paths) Less(a, b int) bool {
	var m, n []int
	for _, slice := range p[a].PathwayIndex {
		m = append(m, slice...)
	}
	for _, slice := range p[b].PathwayIndex {
		n = append(n, slice...)
	}
	N := len(n)
	for k, v := range m {
		if k < N {
			if v == n[k] {
				continue
			}
			return v < n[k]
		}
	}
	return false
}

// Swap swaps element a with element b.
func (p Paths) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
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

	// PathwayIndex represents the pathway when traversed by field index.
	//
	// The inner slice represents field indexes when calling reflect.Value.Field(n).
	// The outer slice represents breaks in the indexes where a pointer occurs.
	//
	// len(PathwayIndex)=1 means there are no pointers in this pathway and
	// reflect.Value.FieldByIndex(PathwayIndex[0]) can reach the field directly.
	PathwayIndex [][]int

	// PathwayOffsets represents the pathway when traversed by pointer arithmetic.
	//
	// Each PathOffsetSegment describes the memory offset from the struct's beginning
	// memory, followed by the number of pointer-indirections required to reach the
	// final type.
	//
	// len(PathwayOffsets)=1 means there are pointers and a single pointer addition
	// from the beginning of the struct's memory will yield the field.
	PathwayOffsets []PathOffsetSegment

	// The full pathway name with nested structs or fields joined by a DOT or PERIOD.
	PathwayName string

	// ParentPathwayName is the pathway name of the parent.
	ParentPathwayName string
}

// ReflectPath returns a condensed representation of Path purpose-built to
// navigate the Path via reflect.
//
// Path contains a moderately excessive amount of information.  Discarding it
// in favor of ReflectPath can lower the memory requirements of packages needing
// this information and functionality.
func (p Path) ReflectPath() ReflectPath {
	var index []int
	for _, indeces := range p.PathwayIndex {
		index = append(index, indeces...)
	}
	n := len(index)
	if n == 0 {
		// If index is empty then this is an in valid value; Last is set to -1
		// to ensure any attempt to call reflect.Value.Field(Last) panics.
		return ReflectPath{Last: -1}
	}
	return ReflectPath{
		HasPointer: len(p.PathwayIndex) > 1,
		Index:      append([]int(nil), index[0:n-1]...),
		Last:       index[n-1],
	}
}

// String returns Path represented as a string.
func (p Path) String() string {
	return fmt.Sprintf("%v %v %v Type=%v Pathway[%v]%v Parent[%v] Offsets= %v", p.Name, p.Index, p.Offset, p.Type, p.PathwayName, PathIndeces(p.PathwayIndex), p.ParentPathwayName, PathOffsets(p.PathwayOffsets))
}

// Value returns the reflect.Value for the path from the origin value.
func (p Path) Value(v reflect.Value) reflect.Value {
	for ; v.Kind() == reflect.Ptr; v = v.Elem() {
		// Walk pointer chain if origin is a pointer
	}
	for _, index := range p.PathwayIndex {
		v = v.FieldByIndex(index)
		for ; v.Kind() == reflect.Ptr; v = v.Elem() {
			if v.IsNil() && v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
	}
	return v
}
