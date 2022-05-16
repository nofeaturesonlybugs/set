package coerce

import "unsafe"

const sizeWord uintptr = unsafe.Sizeof(uintptr(0))

// KindBool unpacks v whose underlying kind must be bool.
// NB Ended up not using this but leaving here for posterity.
// func KindBool(v interface{}) bool {
// 	var e interface{} = false
// 	vp := (*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
// 	ep := (*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
// 	*ep = *vp
// 	return e.(bool)
// }

// KindFloat32 unpacks v whose underlying kind must be float32.
func KindFloat32(v interface{}) float32 {
	var e interface{} = float32(0)
	vp := (*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
	ep := (*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
	*ep = *vp
	return e.(float32)
	// reflect.ValueOf(v).Convert(TypeFloat32).Interface()
}

// KindFloat64 unpacks v whose underlying kind must be float64.
func KindFloat64(v interface{}) float64 {
	var e interface{} = float64(0)
	vp := (*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
	ep := (*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
	*ep = *vp
	return e.(float64)
	//	reflect.ValueOf(v).Convert(TypeFloat64).Interface()
}

// KindInt unpacks v whose underlying kind must be int, int8, int16, int32, or int64.
// NB Ended up not using this but leaving here for posterity.
// func KindInt(v interface{}) int64 {
// 	var e interface{} = int64(0)
// 	vp := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
// 	ep := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
// 	*ep = *vp
// 	return e.(int64)
// 	// NB ~300% slower, 5x more memory
// 	// v = reflect.ValueOf(v).Convert(TypeInt64).Interface()
// }

// KindUint unpacks v whose underlying kind must be uint, uint8, uint16, uint32, or uint64.
// NB Ended up not using this but leaving here for posterity.
// func KindUint(v interface{}) uint64 {
// 	var e interface{} = uint64(0)
// 	vp := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
// 	ep := (*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
// 	*ep = *vp
// 	return e.(uint64)
// 	// NB ~300% slower, 5x more memory
// 	// v = reflect.ValueOf(v).Convert(TypeUint64).Interface()
// }

// KindString unpacks v whose underlying kind must be string.
// NB Ended up not using this but leaving here for posterity.
// func KindString(v interface{}) string {
// 	var e interface{} = ""
// 	vp := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + sizeWord))
// 	ep := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&e)) + sizeWord))
// 	*ep = *vp
// 	return e.(string)
// 	// v = reflect.ValueOf(v).Convert(TypeString).Interface()
// }
