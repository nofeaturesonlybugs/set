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
	var keys []string
	if a, b := len(t.Branches), len(t.Leaves); a > b {
		keys = make([]string, 0, a)
	} else {
		keys = make([]string, 0, b)
	}
	rv := strings.Builder{}

	rv.WriteString("Branches\n")
	for key := range t.Branches {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		rv.WriteString(fmt.Sprintf("\t%v = %v\n", key, t.Branches[key]))
	}

	rv.WriteString("Leaves\n")
	keys = keys[0:0]
	for key := range t.Leaves {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		rv.WriteString(fmt.Sprintf("\t%v = %v\n", key, t.Leaves[key]))
	}
	return rv.String()
}

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
	var stat func(v interface{}, depth int, parent Meta) int
	stat = func(v interface{}, depth int, parent Meta) int {
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
		// if depth == 0 { //TODO RM
		// 	fmt.Println(T) //TODO RM
		// 	depth++        //TODO RM
		// } //TODO RM
		if T.Kind() != reflect.Struct {
			return 0
		}
		// indent := strings.Repeat("\t", depth) // TODO RM for printing
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
					PathwayName:       strings.TrimPrefix(parent.PathwayName+"_"+F.Name, "_"), // TODO Better way
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
					if slice == nil { // Duplicate nils exactly.
						m.PathwayIndex[s] = nil
						continue
					}
					m.PathwayIndex[s] = make([]int, len(slice))
					for k, v := range slice {
						m.PathwayIndex[s][k] = v
					}
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
			// Print for debugging
			// fmt.Printf("%v%v\n", indent, m) //TODO RM
			//
			// Process this field type.
			children := stat(m.EndType, depth+1, m)
			if children == 0 {
				t.Leaves[m.PathwayName] = m.Path
			} else {
				t.Branches[m.PathwayName] = m.Path
			}
		}
		return len(fields)
	}
	stat(v, 0, Meta{})
	return t
}
