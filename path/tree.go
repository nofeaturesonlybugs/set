package path

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Tree is the full mapping of a struct.
type Tree struct {
	// Leaves are struct fields that have no children.
	Leaves map[string]Path

	// Branches are struct fields that have children.
	Branches map[string]Path
}

// Slice returns all members of Leaves and Branches combined into
// a sorted slice.
func (t Tree) Slice() []Path {
	var keys []string
	var slice []Path
	for key := range t.Branches {
		keys = append(keys, key)
	}
	for key := range t.Leaves {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	//
	for _, key := range keys {
		if p, ok := t.Branches[key]; ok {
			slice = append(slice, p)
		} else if p, ok = t.Leaves[key]; ok {
			slice = append(slice, p)
		}
	}
	return slice
}

// String returns the Tree represented as a string.
func (t Tree) String() string {
	return t.StringIndent("\t")
}

// String returns the Tree represented as a string.
func (t Tree) StringIndent(indent string) string {
	rv := strings.Builder{}
	paths := make([]Path, 0, len(t.Leaves))
	//
	for _, br := range t.Branches {
		paths = append(paths, br)
	}
	sort.Sort(Paths(paths))

	rv.WriteString("Branches\n")
	for _, br := range paths {
		in := strings.Repeat(indent, strings.Count(br.PathwayName, ".")+1)
		rv.WriteString(fmt.Sprintf("%v%v = %v\n", in, br.PathwayName, br))
	}
	paths = paths[0:0]

	rv.WriteString("Leaves\n")
	for _, leaf := range t.Leaves {
		paths = append(paths, leaf)
	}
	sort.Sort(Paths(paths))
	for _, leaf := range paths {
		in := strings.Repeat(indent, strings.Count(leaf.PathwayName, ".")+1)
		rv.WriteString(fmt.Sprintf("%v%v = %v\n", in, leaf.PathwayName, leaf))
	}
	return rv.String()
}

// Stat inspects the incoming value to build a Tree which consists of two sets of Paths
// -- branches and leaves.
//
// The incoming value must be a struct, ptr-to-struct, or ptr-chain-to-struct.
//
// Leaves are final fields in the struct that can be traversed no further.  If a field
// is a struct with no exported fields then it is a leaf.
//
// Branches are fields that are structs or embedded structs that can be traversed deeper.
func Stat(v interface{}) Tree {
	t := Tree{
		Leaves:   map[string]Path{},
		Branches: map[string]Path{},
	}
	type Meta struct {
		PtrDeref   int
		PathwayPtr bool
		EndType    reflect.Type
		Path
	}
	//
	var stat func(v interface{}, parent Meta) int
	stat = func(v interface{}, parent Meta) int {
		var T reflect.Type
		switch vv := v.(type) {
		case reflect.Type:
			T = vv
		default:
			T = reflect.TypeOf(vv)
			for T.Kind() == reflect.Ptr {
				T = T.Elem()
			}
		}
		if T.Kind() != reflect.Struct {
			return 0
		}
		//
		fields := make([]Path, 0, T.NumField())
		for fieldIndex, size := 0, T.NumField(); fieldIndex < size; fieldIndex++ {
			F := T.Field(fieldIndex)
			if F.PkgPath != "" {
				continue // PkgPath is non-empty for private fields.
			}
			m := Meta{
				PtrDeref:   0,
				PathwayPtr: parent.PtrDeref > 0 || parent.PathwayPtr,
				EndType:    F.Type,
				Path: Path{
					Index:             fieldIndex,
					Offset:            F.Offset,
					Name:              F.Name,
					ParentPathwayName: parent.PathwayName,
					PathwayName:       strings.TrimPrefix(parent.PathwayName+"."+F.Name, "."),
					Type:              F.Type,
				},
			}
			//
			if len(parent.PathwayIndex) == 0 {
				m.PathwayIndex = [][]int{{fieldIndex}}
				m.PathwayOffsets = []PathOffsetSegment{{
					Offset:  m.Offset,
					Type:    F.Type,
					EndType: F.Type,
				}}
			} else {
				// Copy the parent's pathway index.
				m.PathwayIndex = make([][]int, len(parent.PathwayIndex))
				for s, slice := range parent.PathwayIndex {
					// ↓↓↓ TODO RM This seems to be a lefter over from when a nil could be present. ↓↓↓↓
					// if slice == nil { // Duplicate nils exactly.
					// 	m.PathwayIndex[s] = nil
					// 	continue
					// }
					// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
					m.PathwayIndex[s] = append([]int(nil), slice...)
				}
				// Copy the parent's pathway offsets.
				m.PathwayOffsets = append([]PathOffsetSegment(nil), parent.PathwayOffsets...)
				//
				// If the immediate parent has PtrDeref>0 then it is a pointer and our Pathway fields
				// representing segments from pointers need to be updated.
				if parent.PtrDeref > 0 {
					m.PathwayIndex = append(m.PathwayIndex, nil)
					m.PathwayOffsets = append(m.PathwayOffsets, PathOffsetSegment{})
				}
				//
				// Append our fieldIndex to the last element of the PathwayIndex.
				n := len(m.PathwayIndex)
				m.PathwayIndex[n-1] = append(m.PathwayIndex[n-1], fieldIndex)
				//
				// Update the last element of PathwayOffsets.
				// Add our Offset to existing Offset to represent total Offset from owning struct's memory address.
				n = len(m.PathwayOffsets)
				m.PathwayOffsets[n-1].Offset = m.PathwayOffsets[n-1].Offset + m.Offset
				m.PathwayOffsets[n-1].Type = F.Type
				m.PathwayOffsets[n-1].EndType = F.Type
			}
			//
			FT := F.Type
			for ; FT.Kind() == reflect.Ptr; FT = FT.Elem() {
				m.PtrDeref++
			}
			//
			// If we dereferenced a pointer chain
			if m.PtrDeref > 0 {
				// Update our EndType
				m.EndType = FT
				//
				// Within our Offsets segments update the dereference count and final reflect type.
				// We need this reflect type information to take an unsafe.Pointer and return it back into
				// Go's type system.
				n := len(m.PathwayOffsets)
				m.PathwayOffsets[n-1].IndirectionLevel = m.PtrDeref
				m.PathwayOffsets[n-1].EndType = FT
			}
			fields = append(fields, m.Path)
			//
			// Process this field type.
			children := stat(m.EndType, m)
			if children == 0 {
				t.Leaves[m.PathwayName] = m.Path
			} else {
				t.Branches[m.PathwayName] = m.Path
			}
		}
		return len(fields)
	}
	stat(v, Meta{})
	return t
}
