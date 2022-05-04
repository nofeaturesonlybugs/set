package set

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// V returns a new Value.
//
// Memory is possibly created when calling this function:
//	// No memory is created because b is a local variable and we pass its address.
//	var b bool
//	v := set.V(&b)
//
//	// No memory is created because bp points at an existing local variable and we pass the pointer bp.
//	var b bool
//	bp := &b
//	v := set.V(bp)
//
//	// Memory is created because the local variable is an unallocated pointer AND we pass its address.
//	var *bp bool
//	v := set.V(&bp) // bp now contains allocated memory.
func V(arg interface{}) *Value {
	rv := &Value{}
	rv.original = arg
	//
	if arg == nil {
		rv.err = pkgerr{
			Err:     ErrUnsupported,
			Context: "nil value",
			Hint:    "set.V(nil) was called",
		}
		return rv
	}
	//
	var v reflect.Value
	switch tt := arg.(type) {
	case reflect.Value:
		v = tt
	default:
		v = reflect.ValueOf(arg)
	}
	if v.IsValid() {
		rv.TypeInfo = TypeCache.StatType(v.Type())
	}
	rv.WriteValue, rv.CanWrite = Writable(v)
	rv.TopValue = v

	if rv.IsMap || rv.IsSlice {
		rv.ElemTypeInfo = TypeCache.StatType(rv.ElemType)
	}

	if !rv.CanWrite {
		rv.err = pkgerr{
			Err:     ErrReadOnly,
			Context: rv.Type.String() + " is not writable",
			Hint:    "call to set.V(" + rv.Type.String() + ") should have been set.V(*" + rv.Type.String() + ")",
		}
	}
	return rv
}

// Value wraps around a Go variable and performs magic.
type Value struct {
	// TypeInfo describes the type T in WriteValue.  When the value is created with a pointer P
	// this TypeInfo will describe the final type at the end of the pointer chain.
	//
	// To conserve memory and maintain speed this TypeInfo object may be shared with
	// other *Value instances.  Altering the members within TypeInfo will most likely
	// crash your program with a panic.
	//
	// Treat this value as read only.
	TypeInfo

	// CanWrite specifies if WriteValue.CanSet() would return true.
	CanWrite bool

	// TopValue is the original value passed to V() but wrapped in a reflect.Value.
	TopValue reflect.Value

	// WriteValue is a reflect.Value representing the modifiable value wrapped within this *Value.
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
func (me *Value) Append(items ...interface{}) error {
	if me.err != nil {
		return me.err.(pkgerr).WithCallSite("Value.Append")
	} else if me.Kind != reflect.Slice {
		return pkgerr{Err: ErrUnsupported, CallSite: "Value.Append", Context: "can not append to " + me.Type.String()}
	}
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r) // TODO+NB Update this
			}
		}()
		zero := reflect.Zero(me.Type)
		for _, item := range items {
			elem := reflect.New(me.ElemType)
			elemAsValue := V(elem)
			if err = elemAsValue.To(item); err != nil {
				return
			}
			zero = reflect.Append(zero, reflect.Indirect(elemAsValue.TopValue))
		}
		me.WriteValue.Set(reflect.AppendSlice(me.WriteValue, zero))
	}()
	return err
}

// Copy creates a clone of the *Value and its internal members.
//
// If you need to create many *Value for a type T in order to Rebind(T) in a goroutine
// architecture then consider creating and caching a V(T) early in your application
// and then calling Copy() on that cached copy before using Rebind().
func (me *Value) Copy() *Value {
	rv := &Value{
		TypeInfo:     me.TypeInfo,
		CanWrite:     me.CanWrite,
		TopValue:     me.TopValue,
		WriteValue:   me.WriteValue,
		ElemTypeInfo: me.ElemTypeInfo,
		original:     me.original,
		err:          me.err,
	}
	return rv
}

// Fields returns a slice of Field structs when Value is wrapped around a struct; for all other values
// nil is returned.
//
// This function has some overhead because it creates a new *Value for each struct field.  If you only need
// the reflect.StructField information consider using the public StructFields member.
func (me *Value) Fields() []Field {
	if me == nil || me.Kind != reflect.Struct {
		return nil
	}
	var rv []Field
	for k, max := 0, me.Type.NumField(); k < max; k++ {
		v, f := me.WriteValue.Field(k), me.Type.Field(k)
		rv = append(rv, Field{Value: V(v), Field: f})
	}
	return rv
}

// FieldByIndex returns the nested field corresponding to index.
//
// Key differences between this method and the built-in method on reflect.Value.FieldByIndex() are
// the built-in causes panics while this one will return errors and this method will instantiate nil struct
// members as it traverses.
func (me *Value) FieldByIndex(index []int) (reflect.Value, error) {
	var v reflect.Value
	size := len(index)
	if me.err != nil {
		return v, me.err.(pkgerr).WithCallSite("Value.FieldByIndex")
	} else if size == 0 {
		return v, pkgerr{Err: ErrUnsupported, CallSite: "Value.FieldByIndex", Context: "empty index"}
	}
	v = me.WriteValue
	for k := 0; k < size; k++ {
		n := index[k] // n is the index (or field num) to consider
		if v.Kind() != reflect.Struct {
			return v, pkgerr{Err: ErrUnsupported, CallSite: "Value.FieldByIndex", Context: fmt.Sprintf("want struct but got %v", v.Type())}
		} else if n > v.NumField()-1 {
			return v, pkgerr{Err: ErrIndexOutOfBounds, CallSite: "Value.FieldByIndex", Context: fmt.Sprintf("index %v exceeds max %v", n, v.NumField()-1)}
		}
		v = v.Field(n)
		t, k := v.Type(), v.Kind()
		// Instantiate nil pointer chains.
		if k == reflect.Ptr {
			for k == reflect.Ptr {
				if v.IsNil() && v.CanSet() {
					ptr := reflect.New(t.Elem())
					v.Set(ptr)
					v = ptr
				}
				v = v.Elem()
				t, k = v.Type(), v.Kind()
			}
		}
	}
	return v, nil
}

// FieldByIndexAsValue calls into FieldByIndex and if there is no error the resulting reflect.Value is
// wrapped within a call to V() to return a *Value.
func (me *Value) FieldByIndexAsValue(index []int) (*Value, error) {
	var v reflect.Value
	var err error
	if v, err = me.FieldByIndex(index); err != nil {
		return nil, err
	}
	return V(v), nil
}

// FieldsByTag is the same as Fields() except only Fields with the given struct-tag are returned and the
// TagValue member of Field will be set to the tag's value.
func (me *Value) FieldsByTag(key string) []Field {
	var rv []Field
	for _, f := range me.Fields() {
		if value, ok := f.Field.Tag.Lookup(key); ok {
			f.TagValue = value
			rv = append(rv, f)
		}
	}
	return rv
}

// fill is the underlying function that powers Fill() and FillByTag().
//
// getter is the original Getter passed to Fill() or FillByTag().
//
// Fill() and FillByTag() have essentially the same complicated logic except where they get the string/key to pass
// to getter() and how they sub-fill nested structures.  The keyFunc and fillFunc arguments allow them to
// cascade the appropriate logic into this function.
func (me *Value) fill(getter Getter, fields []Field, keyFunc func(Field) string, fillFunc func(*Value, Getter) error) error {
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
				field.Value.Append(elem.WriteValue.Interface()) // This can return an error but it _should_be_ impossible.
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
					field.Value.Append(elem.WriteValue.Interface()) // This can return an error but it _should_be impossible.
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
func (me *Value) Fill(getter Getter) error {
	fields := me.Fields()
	keyFunc := func(field Field) string {
		return field.Field.Name
	}
	fillFunc := func(value *Value, getter Getter) error {
		return value.Fill(getter)
	}
	return me.fill(getter, fields, keyFunc, fillFunc)
}

// FillByTag is the same as Fill() except the argument passed to Getter is the value of the struct-tag.
func (me *Value) FillByTag(key string, getter Getter) error {
	fields := me.FieldsByTag(key)
	keyFunc := func(field Field) string {
		return field.TagValue
	}
	fillFunc := func(value *Value, getter Getter) error {
		return value.FillByTag(key, getter)
	}
	return me.fill(getter, fields, keyFunc, fillFunc)
}

// Rebind will swap the underlying original value used to create *Value with the incoming
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
//		v := set.V( item ) // Creates new *Value every iteration -- can be expensive!
//		// manipulate v in order to affect item
//	}
// to:
//	var slice []T
//	v := set.V( T{} ) // Create a single *Value for the type T
//	// populate slice
//	for _, item := range slice {
//		v.Rebind( item ) // Reuse the existing *Value -- will be faster!
//		// manipulate v in order to affect item
//	}
//
func (me *Value) Rebind(arg interface{}) {
	var v reflect.Value
	switch tt := arg.(type) {
	case reflect.Value:
		v = tt
	case *Value:
		v = tt.TopValue
	default:
		v = reflect.ValueOf(arg)
	}
	if me.TopValue.Type() != v.Type() {
		panic(fmt.Sprintf("Rebind expects same underlying type: original %T not compatible with incoming %T", me.WriteValue.Interface(), arg))
	}
	me.original, me.TopValue = arg, v
	me.WriteValue, me.CanWrite = Writable(v)
}

// Zero sets the Value to the Zero value of the appropriate type.
func (me *Value) Zero() error {
	if me.err != nil {
		return me.err.(pkgerr).WithCallSite("Value.Zero")
	}
	me.WriteValue.Set(reflect.Zero(me.Type))
	return nil
}

// NewElem instantiates and returns a *Value that can be Panics.Append()'ed to this type; only valid
// if Value.ElemType describes a valid type.
func (me *Value) NewElem() (*Value, error) {
	if me.ElemTypeInfo.Kind == reflect.Invalid {
		return nil, pkgerr{Err: ErrUnsupported, CallSite: "Value.NewElem", Context: fmt.Sprintf("%T is not an element container", me.original)}
	}
	return V(reflect.New(me.ElemType)), nil
}

// To attempts to assign the argument into Value.
//
// If *Value is wrapped around an unwritable reflect.Value or the type is reflect.Invalid an
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
func (me *Value) To(arg interface{}) error {
	if me.err != nil {
		return me.err.(pkgerr).WithCallSite("Value.To")
	} else if arg == nil {
		return me.Zero()
	}
	//
	// This typeswitch handles the case where we are a scalar.
	// TODO Possibly benchmark switching on reflect.Kind vs type switching on me.original
	switch me.Kind {
	case reflect.Bool:
		c, err := coerce.Bool(arg)
		me.WriteValue.SetBool(c)
		return err
	case reflect.Float32:
		c, err := coerce.Float32(arg)
		me.WriteValue.SetFloat(float64(c))
		return err
	case reflect.Float64:
		c, err := coerce.Float64(arg)
		me.WriteValue.SetFloat(c)
		return err
	case reflect.Int:
		c, err := coerce.Int(arg)
		me.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int8:
		c, err := coerce.Int8(arg)
		me.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int16:
		c, err := coerce.Int16(arg)
		me.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int32:
		c, err := coerce.Int32(arg)
		me.WriteValue.SetInt(int64(c))
		return err
	case reflect.Int64:
		c, err := coerce.Int64(arg)
		me.WriteValue.SetInt(int64(c))
		return err
	case reflect.Uint:
		c, err := coerce.Uint(arg)
		me.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint8:
		c, err := coerce.Uint8(arg)
		me.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint16:
		c, err := coerce.Uint16(arg)
		me.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint32:
		c, err := coerce.Uint32(arg)
		me.WriteValue.SetUint(uint64(c))
		return err
	case reflect.Uint64:
		c, err := coerce.Uint64(arg)
		me.WriteValue.SetUint(uint64(c))
		return err
	case reflect.String:
		c, err := coerce.String(arg)
		me.WriteValue.SetString(c)
		return err
	}
	//
	rv := reflect.ValueOf(arg)
	for {
		//
		// If arg is any kind of pointer dereference to final value or nil
		for ; rv.Kind() == reflect.Ptr; rv = rv.Elem() {
			if rv.IsNil() {
				return me.Zero()
			}
		}
		T := rv.Type()
		//
		if (T == me.Type || T.AssignableTo(me.Type)) && me.Kind != reflect.Slice {
			// NB  We checked that me.Kind is not a slice because this package always makes a copy of a slice!
			me.WriteValue.Set(rv)
			return nil
		}
		//
		if me.IsSlice {
			me.Zero() // Zero only returns errors on nil receiver, invalid kind, or !CanWrite -- which are already checked above.
			if rv.Kind() != reflect.Slice {
				arg = []interface{}{arg}
			}
			slice := reflect.ValueOf(arg)
			for k, size := 0, slice.Len(); k < size; k++ {
				elem := V(reflect.New(me.ElemType).Interface())
				if err := elem.To(slice.Index(k).Interface()); err != nil {
					me.Zero()
					return err
				}
				me.WriteValue.Set(reflect.Append(me.WriteValue, elem.WriteValue))
			}
			return nil
		} else if rv.Kind() == reflect.Slice {
			// When incoming value is a slice we use the last value and try again.
			if n := rv.Len(); n > 0 {
				rv = rv.Index(n - 1)
				continue
			}
			return me.Zero()
		}
		return me.Zero()
	}
}
