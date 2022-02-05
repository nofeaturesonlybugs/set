package coerce

import "reflect"

// Calculate and cache some common reflect.Type values needed
// during type coercion.
var (
	TypeBool    = reflect.TypeOf(false)
	TypeFloat32 = reflect.TypeOf(float32(0))
	TypeFloat64 = reflect.TypeOf(float64(0))
	TypeInt64   = reflect.TypeOf(int64(0))
	TypeUint64  = reflect.TypeOf(uint64(0))
	TypeString  = reflect.TypeOf("")
)
