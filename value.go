package set

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// zeroV is created once and returned whenever an empty or invalid Value is needed.
var zeroV = V(nil)

// V returns a new Value.
//
// The returned Value must not be copied except via its Copy method.
func V(arg interface{}) Value {
	var v Value
	v.original = arg
	//
	if arg == nil {
		v.err = pkgerr{
			Err:     ErrUnsupported,
			Context: "nil value",
			Hint:    "set.V(nil) was called",
		}
		return v
	}
	//
	var V reflect.Value
	switch tt := arg.(type) {
	case reflect.Value:
		V = tt
	default:
		V = reflect.ValueOf(arg)
	}
	if V.IsValid() {
		v.TypeInfo = TypeCache.StatType(V.Type())
	}
	v.WriteValue, v.CanWrite = Writable(V)
	v.TopValue = V

	if v.IsMap || v.IsSlice {
		v.ElemTypeInfo = TypeCache.StatType(v.ElemType)
	}

	if !v.CanWrite {
		v.err = pkgerr{
			Err:     ErrReadOnly,
			Context: v.Type.String() + " is not writable",
			Hint:    "call to set.V(" + v.Type.String() + ") should have been set.V(*" + v.Type.String() + ")",
		}
	}
	return v
}

// Value wraps around a Go variable and performs magic.
//
// Once created a Value should only be copied via its Copy method.
type Value struct {
	// TypeInfo describes the type T in WriteValue.  When the value is created with a pointer P
	// this TypeInfo will describe the final type at the end of the pointer chain.
	//
	// To conserve memory and maintain speed this TypeInfo object may be shared with
	// other Value instances.  Altering the members within TypeInfo will most likely
	// crash your program with a panic.
	//
	// Treat this value as read only.
	TypeInfo

	// CanWrite specifies if WriteValue.CanSet() would return true.
	CanWrite bool

	// TopValue is the original value passed to V() but wrapped in a reflect.Value.
	TopValue reflect.Value

	// WriteValue is a reflect.Value representing the modifiable value wrapped within this Value.
	//
	// If you call V( &t ) then CanWrite will be true and WriteValue will be a usable reflect.Value.
	// If you call V( t ) where t is not a pointer or does not point to allocated memory then
	// CanWrite will be false and any attempt to set values on WriteValue will probably panic.
	//
	// All methods on this type that alter the value Append(), Fill*(), To(), etc work on this
	// value.  Generally you should avoid it but it's also present if you really know what you're doing.
	WriteValue reflect.Value

	// When IsMap or IsSlice are true then ElemTypeInfo is a TypeInfo struct describing the element type.
	ElemTypeInfo TypeInfo

	//
	err      error
	original interface{}
}

// Append appends the item(s) to the end of the Value assuming it is some type of slice and every
// item can be type-coerced into the slice's data type.  Either all items are appended without an error
// or no items are appended and an error is returned describing the type of the item that could not
// be appended.
func (v Value) Append(items ...interface{}) error {
	if v.err != nil {
		return v.err.(pkgerr).WithCallSite("Value.Append")
	} else if v.Kind != reflect.Slice {
		return pkgerr{Err: ErrUnsupported, CallSite: "Value.Append", Context: "can not append to " + v.Type.String()}
	}
	//
	zero := reflect.Zero(v.Type)
	for _, item := range items {
		elem := reflect.New(v.ElemType)
		elemAsValue := V(elem)
		if err := elemAsValue.To(item); err != nil {
			return err
		}
		zero = reflect.Append(zero, reflect.Indirect(elemAsValue.TopValue))
	}
	v.WriteValue.Set(reflect.AppendSlice(v.WriteValue, zero))
	//
	return nil
}

// Copy creates a clone of the Value and its internal members.
//
// If you need to create many Value for a type T in order to Rebind(T) in a goroutine
// architecture then consider creating and caching a V(T) early in your application
// and then calling Copy() on that cached copy before using Rebind().
func (v Value) Copy() Value {
	cp := Value{
		TypeInfo:     v.TypeInfo,
		CanWrite:     v.CanWrite,
		TopValue:     v.TopValue,
		WriteValue:   v.WriteValue,
		ElemTypeInfo: v.ElemTypeInfo,
		original:     v.original,
		err:          v.err,
	}
	return cp
}

// Fields returns a slice of Field structs when Value is wrapped around a struct; for all other values
// nil is returned.
//
// This function has some overhead because it creates a new Value for each struct field.  If you only need
// the reflect.StructField information consider using the public StructFields member.
func (v Value) Fields() []Field {
	if v.Kind != reflect.Struct {
		return nil
	}
	var fields []Field
	for k, max := 0, v.Type.NumField(); k < max; k++ {
		fV, fT := v.WriteValue.Field(k), v.Type.Field(k)
		fields = append(fields, Field{Value: V(fV), Field: fT})
	}
	return fields
}

// FieldByIndex returns the nested field corresponding to index.
//
// Key differences between this method and the built-in method on reflect.Value.FieldByIndex() are
// the built-in causes panics while this one will return errors and this method will instantiate nil struct
// members as it traverses.
func (v Value) FieldByIndex(index []int) (reflect.Value, error) {
	var V reflect.Value
	size := len(index)
	if v.err != nil {
		return V, v.err.(pkgerr).WithCallSite("Value.FieldByIndex")
	} else if size == 0 {
		return V, pkgerr{Err: ErrUnsupported, CallSite: "Value.FieldByIndex", Context: "empty index"}
	}
	V = v.WriteValue
	for k := 0; k < size; k++ {
		n := index[k] // n is the index (or field num) to consider
		if V.Kind() != reflect.Struct {
			return V, pkgerr{Err: ErrUnsupported, CallSite: "Value.FieldByIndex", Context: fmt.Sprintf("want struct but got %v", V.Type())}
		} else if n > V.NumField()-1 {
			return V, pkgerr{Err: ErrIndexOutOfBounds, CallSite: "Value.FieldByIndex", Context: fmt.Sprintf("index %v exceeds max %v", n, V.NumField()-1)}
		}
		V = V.Field(n)
		T, K := V.Type(), V.Kind()
		// Instantiate nil pointer chains.
		if K == reflect.Ptr {
			for K == reflect.Ptr {
				if V.IsNil() && V.CanSet() {
					ptr := reflect.New(T.Elem())
					V.Set(ptr)
					V = ptr
				}
				V = V.Elem()
				T, K = V.Type(), V.Kind()
			}
		}
	}
	return V, nil
}

// FieldByIndexAsValue calls into FieldByIndex and if there is no error the resulting reflect.Value is
// wrapped within a call to V() to return a Value.
func (v Value) FieldByIndexAsValue(index []int) (Value, error) {
	var RV reflect.Value
	var err error
	if RV, err = v.FieldByIndex(index); err != nil {
		return zeroV, err
	}
	return V(RV), nil
}

// FieldsByTag is the same as Fields() except only Fields with the given struct-tag are returned and the
// TagValue member of Field will be set to the tag's value.
func (v Value) FieldsByTag(key string) []Field {
	var fields []Field
	for _, f := range v.Fields() {
		if value, ok := f.Field.Tag.Lookup(key); ok {
			f.TagValue = value
			fields = append(fields, f)
		}
	}
	return fields
}

// fill is the underlying function that powers Fill() and FillByTag().
//
// getter is the original Getter passed to Fill() or FillByTag().
//
// Fill() and FillByTag() have essentially the same complicated logic except where they get the string/key to pass
// to getter() and how they sub-fill nested structures.  The keyFunc and fillFunc arguments allow them to
// cascade the appropriate logic into this function.
func (v Value) fill(getter Getter, fields []Field, keyFunc func(Field) string, fillFunc func(Value, Getter) error) error {
	var err error
	for _, field := range fields {
		getName := keyFunc(field)
		switch got := getter.Get(getName).(type) {

		case Getter:
			// What was returned from the Getter is itself a Getter; therefore we expect field.Value
			// to be either a struct or []struct that we can sub-fill.
			if field.Value.IsStruct {
				if err = fillFunc(field.Value, got); err != nil {
					return err
				}
			} else if field.Value.IsSlice && field.Value.ElemTypeInfo.IsStruct {
				if err = field.Value.Zero(); err != nil {
					return err
				}
				elem := V(reflect.New(field.Value.ElemTypeInfo.Type))
				if err = fillFunc(elem, got); err != nil {
					return err
				}
				_ = field.Value.Append(elem.WriteValue.Interface()) // It is impossible for this to return an error here.
			} else {
				return pkgerr{Err: ErrUnsupported, CallSite: "Value.fill", Context: fmt.Sprintf("value is Getter but field %v is %v", getName, field.Value.Type)}
			}

		case []Getter:
			// What was returned from the Getter is a []Getter; therefore we expect field.Value to
			// be a []struct or struct that we can sub-fill.
			if field.Value.IsSlice && field.Value.ElemTypeInfo.IsStruct {
				// Zero out the existing slice.
				if err = field.Value.Zero(); err != nil {
					return err
				}
				for _, elemGetter := range got {
					elem := V(reflect.New(field.Value.ElemTypeInfo.Type))
					if err = fillFunc(elem, elemGetter); err != nil {
						return err
					}
					_ = field.Value.Append(elem.WriteValue.Interface()) // It is impossible for this to return an error here.
				}
			} else if field.Value.IsStruct {
				size := len(got)
				if size > 0 {
					if err = fillFunc(field.Value, got[size-1]); err != nil {
						return err
					}
				}
			} else {
				return pkgerr{Err: ErrUnsupported, CallSite: "Value.fill", Context: fmt.Sprintf("value is []Getter but field %v is %v", getName, field.Value.Type)}
			}

		default:
			if err = field.Value.To(got); err != nil {
				return err
			}
		}
	}
	return nil
}

// Fill iterates a struct's fields and calls To() on each one by passing the field name to the Getter.
// Fill stops and returns on the first error encountered.
func (v Value) Fill(getter Getter) error {
	fields := v.Fields()
	keyFunc := func(field Field) string {
		return field.Field.Name
	}
	fillFunc := func(value Value, getter Getter) error {
		return value.Fill(getter)
	}
	return v.fill(getter, fields, keyFunc, fillFunc)
}

// FillByTag is the same as Fill() except the argument passed to Getter is the value of the struct-tag.
func (v Value) FillByTag(key string, getter Getter) error {
	fields := v.FieldsByTag(key)
	keyFunc := func(field Field) string {
		return field.TagValue
	}
	fillFunc := func(value Value, getter Getter) error {
		return value.FillByTag(key, getter)
	}
	return v.fill(getter, fields, keyFunc, fillFunc)
}

// Rebind will swap the underlying original value used to create Value with the incoming
// value if:
//	Type(Original) == Type(Incoming).
//
// If Rebind succeeds the following public members will have been replaced appropriately:
//	CanWrite
//	TopValue
//	WriteValue
//
// Reach for this function to translate:
//	var slice []T
//	// populate slice
//	for _, item := range slice {
//		v := set.V( item ) // Creates new Value every iteration -- can be expensive!
//		// manipulate v in order to affect item
//	}
// to:
//	var slice []T
//	v := set.V( T{} ) // Create a single Value for the type T
//	// populate slice
//	for _, item := range slice {
//		v.Rebind( item ) // Reuse the existing Value -- will be faster!
//		// manipulate v in order to affect item
//	}
//
func (v *Value) Rebind(arg interface{}) {
	var V reflect.Value
	switch tt := arg.(type) {
	case reflect.Value:
		V = tt
	case Value:
		V = tt.TopValue
	default:
		V = reflect.ValueOf(arg)
	}
	if v.TopValue.Type() != V.Type() {
		panic(fmt.Sprintf("mismatching types during Rebind; have %T and got %T", v.original, arg))
	}
	v.original, v.TopValue = arg, V
	v.WriteValue, v.CanWrite = Writable(V)
}

// Zero sets the Value to the Zero value of the appropriate type.
func (v Value) Zero() error {
	if v.err != nil {
		return v.err.(pkgerr).WithCallSite("Value.Zero")
	}
	v.WriteValue.Set(reflect.Zero(v.Type))
	return nil
}

// NewElem instantiates and returns a Value that can be Panics.Append()'ed to this type; only valid
// if Value.ElemType describes a valid type.
func (v Value) NewElem() (Value, error) {
	if v.ElemTypeInfo.Kind == reflect.Invalid {
		return zeroV, pkgerr{Err: ErrUnsupported, CallSite: "Value.NewElem", Context: fmt.Sprintf("%T is not an element container", v.original)}
	}
	return V(reflect.New(v.ElemType)), nil
}

// To attempts to assign the argument into Value.
//
// If Value is wrapped around an unwritable reflect.Value or the type is reflect.Invalid an
// error will be returned.  You probably forgot to call set.V() with an address to your type.
//
// If the assignment can not be made but the wrapped value is writable then the wrapped
// value will be set to an appropriate zero type to overwrite any existing data.
//
// 	set.V(&T).To(S)
//
//	T is scalar, S is scalar, same type
//		-> direct assignment
//
//	If S is a pointer then dereference until final S value and continue...
//
//	T is scalar, S is scalar, different types
//		-> assignment with attempted type coercion
//	T is scalar, S is slice []S
//		-> T is assigned S[ len( S ) - 1 ]; i.e. last element in S if length greater than 0.
//	T is slice []T, S is scalar
//		-> T is set to []T{ S }; i.e. a slice of T with S as the only element.
//	T is slice []T, S is slice []S
//		-> T is set to []T{ S... }; i.e. a new slice with elements from S copied.
//		-> Note: T != S; they are now different slices; changes to T do not affect S and vice versa.
//		-> Note: If the elements themselves are pointers then, for example, T[0] and S[0] point
//			at the same memory and will see changes to whatever is pointed at.
func (v Value) To(arg interface{}) error {
	if v.err != nil {
		return v.err.(pkgerr).WithCallSite("Value.To")
	} else if arg == nil {
		return v.Zero()
	}
	//
	// This typeswitch handles the case where we are a scalar.
	switch v.Kind {
	case reflect.Bool:
		c, err := coerce.Bool(arg)
		v.WriteValue.SetBool(c)
		return err
	case reflect.Float32:
		c, err := coerce.Float32(arg)
		v.WriteValue.SetFloat(float64(c))
		return err
	case reflect.Float64:
		c, err := coerce.Float64(arg)
		v.WriteValue.SetFloat(c)
		return err
	case reflect.Int:
		c, err := coerce.Int(arg)
		v.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int8:
		c, err := coerce.Int8(arg)
		v.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int16:
		c, err := coerce.Int16(arg)
		v.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int32:
		c, err := coerce.Int32(arg)
		v.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int64:
		c, err := coerce.Int64(arg)
		v.WriteValue.SetInt(int64(c))
		return err
	case reflect.Uint:
		c, err := coerce.Uint(arg)
		v.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint8:
		c, err := coerce.Uint8(arg)
		v.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint16:
		c, err := coerce.Uint16(arg)
		v.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint32:
		c, err := coerce.Uint32(arg)
		v.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint64:
		c, err := coerce.Uint64(arg)
		v.WriteValue.SetUint(uint64(c))
		return err
	case reflect.String:
		c, err := coerce.String(arg)
		v.WriteValue.SetString(c)
		return err
	}
	//
	// TODO Code below here could be optimized better.
	rv := reflect.ValueOf(arg)
	for {
		//
		// If arg is any kind of pointer dereference to final value or nil
		for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
			if rv.IsNil() {
				return v.Zero()
			}
		}
		T := rv.Type()
		//
		if (T == v.Type || T.AssignableTo(v.Type)) && v.Kind != reflect.Slice {
			// NB  We checked that me.Kind is not a slice because this package always makes a copy of a slice!
			v.WriteValue.Set(rv)
			return nil
		}
		//
		if v.IsSlice {
			_ = v.Zero() // Zero only returns errors on nil receiver, invalid kind, or !CanWrite -- which are already checked above.
			if rv.Kind() != reflect.Slice {
				arg = []interface{}{arg}
			}
			slice := reflect.ValueOf(arg)
			for k, size := 0, slice.Len(); k < size; k++ {
				elem := V(reflect.New(v.ElemType))
				if err := elem.To(slice.Index(k).Interface()); err != nil {
					_ = v.Zero()
					return err
				}
				v.WriteValue.Set(reflect.Append(v.WriteValue, elem.WriteValue))
			}
			return nil
		} else if rv.Kind() == reflect.Slice {
			// When incoming value is a slice we use the last value and try again.
			if n := rv.Len(); n > 0 {
				rv = rv.Index(n - 1)
				continue
			}
			return v.Zero()
		}
		return v.Zero()
	}
}
