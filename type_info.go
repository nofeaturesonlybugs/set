package set

import (
	"reflect"
	"sync"
)

// TypeInfo describes information about a type that is pertinent to this package.
type TypeInfo struct {
	// True if the Value is a scalar type:
	//	bool, float32, float64, string
	//	int, int8, int16, int32, int64
	//	uint, uint8, uint16, uint32, uint64
	IsScalar bool

	// True if the Value is a map.
	IsMap bool

	// True if the Value is a slice.
	IsSlice bool

	// True if the Value is a struct.
	IsStruct bool

	// Kind is the reflect.Kind; when Stat() or StatType() were called with a pointer this will be the final
	// kind at the end of the pointer chain.  Otherwise it will be the original kind.
	Kind reflect.Kind

	// Type is the reflect.Type; when Stat() or StatType() were called with a poiner this will be the final
	// type at the end of the pointer chain.  Otherwise it will be the original type.
	Type reflect.Type

	// When IsMap or IsSlice are true then ElemType will be the reflect.Type for elements that can be directly
	// inserted into the map or slice; it is not the type at the end of the chain if the element type is a pointer.
	ElemType reflect.Type

	// When IsStruct is true then StructFields will contain the reflect.StructField values for the struct.
	StructFields []reflect.StructField
}

// TypeInfoCache builds a cache of TypeInfo types; when requesting TypeInfo for a type T that is a pointer
// the TypeInfo returned will describe the type *T (or ******T) at the end of the pointer chain.
//
// When Stat() or StatType() are called with nil or an Interface(nil) a zero TypeInfo is returned; essentially
// nothing useful can be done with the type needed to be described.
type TypeInfoCache interface {
	// Stat accepts an arbitrary variable and returns the associated TypeInfo structure.
	Stat(T interface{}) TypeInfo
	// StatType is the same as Stat() except it expects a reflect.Type.
	StatType(T reflect.Type) TypeInfo
}

var TypeCache = NewTypeInfoCache()

// NewTypeInfoCache creates a new TypeInfoCache.
func NewTypeInfoCache() TypeInfoCache {
	return &type_info_cache_t{
		cache: &sync.Map{},
	}
}

// type_info_cache_t is the implementation of a TypeInfoCache for this package.
type type_info_cache_t struct {
	// Performance note:
	//	Initially this was a map[reflect.Type]TypeInfo and we used a sync.RWMutex to control
	//	access.  Switching to sync.Map removes the need for the RWMutex and changed
	//	performance stats for StatType()
	//		from:		360ms, 17.82% of Total
	//		to:			120ms, 11.21% of Total
	cache *sync.Map
}

// Stat accepts an arbitrary variable and returns the associated TypeInfo structure.
func (me *type_info_cache_t) Stat(T interface{}) TypeInfo {
	t := reflect.TypeOf(T)
	return me.StatType(t)
}

// StatType is the same as Stat() except it expects a reflect.Type.
func (me *type_info_cache_t) StatType(T reflect.Type) TypeInfo {
	if T == nil {
		return TypeInfo{}
	}
	if rv, ok := me.cache.Load(T); ok {
		return rv.(TypeInfo)
	}
	//
	T_orig := T
	//
	rv := TypeInfo{}
	V := reflect.New(T)
	T = V.Type()
	K := V.Kind()
	//
	for K == reflect.Ptr {
		if V.IsNil() && V.CanSet() {
			ptr := reflect.New(T.Elem())
			V.Set(ptr)
		}
		K = T.Elem().Kind()
		T = T.Elem()
		V = V.Elem()
	}
	//
	rv.IsMap = K == reflect.Map
	rv.IsSlice = K == reflect.Slice
	rv.IsStruct = K == reflect.Struct
	rv.IsScalar = K == reflect.Bool ||
		K == reflect.Int || K == reflect.Int8 || K == reflect.Int16 || K == reflect.Int32 || K == reflect.Int64 ||
		K == reflect.Uint || K == reflect.Uint8 || K == reflect.Uint16 || K == reflect.Uint32 || K == reflect.Uint64 ||
		K == reflect.Float32 || K == reflect.Float64 ||
		K == reflect.String
	if rv.IsMap || rv.IsSlice {
		rv.ElemType = T.Elem()
	} else if rv.IsStruct {
		for k, size := 0, T.NumField(); k < size; k++ {
			rv.StructFields = append(rv.StructFields, T.Field(k))
		}
	}
	rv.Type, rv.Kind = T, K
	//
	me.cache.Store(T_orig, rv)
	//
	return rv
}
