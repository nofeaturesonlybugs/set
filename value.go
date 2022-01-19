package set

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/errors"
)

var (
	// ErrReadOnly is returned by methods requring writable access to a variable but it was not
	// passed by address and is readonly.  It is commonly wrapped with a %T verb to describe the
	// type that was passed.
	ErrReadOnly = fmt.Errorf("set: value is readonly; pass by address")
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

	// When IsMap or IsSlice are true then ElemTypeInfo is a TypeInfo struct describing the element types.
	ElemTypeInfo TypeInfo

	//
	original interface{}
}

// errorUnsupported returns a string that can be used in an error message to indicate the underlying original type
// does not support the requested operation.
func (me *Value) errorUnsupported(method string) string {
	return fmt.Sprintf("%v is unsupported for original type [%T]", method, me.original)
}

// Append appends the item(s) to the end of the Value assuming it is some type of slice and every
// item can be type-coerced into the slice's data type.  Either all items are appended without an error
// or no items are appended and an error is returned describing the type of the item that could not
// be appended.
func (me *Value) Append(items ...interface{}) error {
	if me == nil {
		return errors.NilReceiver()
	} else if me.Kind != reflect.Slice {
		return errors.Errorf(me.errorUnsupported("Append"))
	}
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("%v", r)
			}
		}()
		zero := reflect.Zero(me.Type)
		for _, item := range items {
			elem := reflect.New(me.ElemType)
			elemAsValue := V(elem)
			if err = elemAsValue.To(item); err != nil {
				err = errors.Go(err)
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
	if me != nil && me.IsStruct {
		for k, max := 0, me.Type.NumField(); k < max; k++ {
			v, f := me.WriteValue.Field(k), me.Type.Field(k)
			rv = append(rv, Field{Value: V(v), Field: f})
		}
	}
	return rv
}

// FieldByIndex returns the nested field corresponding to index.
//
// Key differences between this method and the built-in method on reflect.Value.FieldByIndex() are
// the built-in causes panics while this one will return errors and this method will instantiate nil struct
// members as it traverses.
func (me *Value) FieldByIndex(index []int) (reflect.Value, error) {
	v := reflect.Value{}
	size := len(index)
	if me == nil {
		return v, errors.NilReceiver()
	} else if !me.CanWrite {
		return v, errors.Errorf(me.errorUnsupported("FieldByIndex"))
	} else if size == 0 {
		return v, errors.Errorf("Zero length index provided to FieldByIndex()")
	}
	v = me.WriteValue
	for k := 0; k < size; k++ {
		n := index[k] // n is the index (or field num) to consider
		if v.Kind() != reflect.Struct {
			return v, errors.Errorf("FieldByIndex requires type to be a struct; type is %v", v.Type())
		} else if n > v.NumField() {
			return v, errors.Errorf("Index out of bounds; field is len %v and index is %v", v.NumField(), n)
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
		return nil, errors.Go(err)
	}
	return V(v), nil
}

// FieldsByTag is the same as Fields() except only Fields with the given struct-tag are returned and the
// TagValue member of Field will be set to the tag's value.
func (me *Value) FieldsByTag(key string) []Field {
	if me == nil || me.Kind != reflect.Struct {
		return nil
	}
	var rv []Field
	all := me.Fields()
	for _, f := range all {
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
					return errors.Go(err)
				}
			} else if field.Value.IsSlice && field.Value.ElemTypeInfo.IsStruct {
				if err = field.Value.Zero(); err != nil {
					return errors.Go(err)
				}
				elem := V(reflect.New(field.Value.ElemTypeInfo.Type))
				if err = fillFunc(elem, got); err != nil {
					return errors.Go(err)
				}
				field.Value.Append(elem.WriteValue.Interface()) // This can return an error but it _should_be_ impossible.
			} else {
				return errors.Errorf("Getter.Get( %v ) returned a Getter for field %v and field is not fillable.", getName, field.Field.Name)
			}

		case []Getter:
			// What was returned from the Getter is a []Getter; therefore we expect field.Value to
			// be a []struct or struct that we can sub-fill.
			if field.Value.IsSlice && field.Value.ElemTypeInfo.IsStruct {
				// Zero out the existing slice.
				if err = field.Value.Zero(); err != nil {
					return errors.Go(err)
				}
				for _, elemGetter := range got {
					elem := V(reflect.New(field.Value.ElemTypeInfo.Type))
					if err = fillFunc(elem, elemGetter); err != nil {
						return errors.Go(err)
					}
					field.Value.Append(elem.WriteValue.Interface()) // This can return an error but it _should_be impossible.
				}
			} else if field.Value.IsStruct {
				size := len(got)
				if size > 0 {
					if err = fillFunc(field.Value, got[size-1]); err != nil {
						return errors.Go(err)
					}
				}
			} else {
				return errors.Errorf("Getter.Get( %v ) returned a []Getter for field %v and field is not fillable.", getName, field.Field.Name)
			}

		default:
			if err = field.Value.To(got); err != nil {
				return errors.Go(err)
			}
		}
	}
	return nil
}

// Fill iterates a struct's fields and calls Set() on each one by passing the field name to the Getter.
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
	if me == nil {
		return errors.NilReceiver()
	} else if !me.CanWrite || me.Kind == reflect.Invalid {
		return errors.Errorf(me.errorUnsupported("Zero"))
	}
	me.WriteValue.Set(reflect.Zero(me.Type))
	return nil
}

// NewElem instantiates and returns a *Value that can be Panics.Append()'ed to this type; only valid
// if Value.ElemType describes a valid type.
func (me *Value) NewElem() (*Value, error) {
	if me == nil {
		return nil, errors.NilReceiver()
	} else if me.ElemTypeInfo.Kind == reflect.Invalid {
		return nil, errors.Errorf(me.errorUnsupported("NewElem"))
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
//	T is pointer, S is pointer, same type and level of indirection
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
	// Performance note(s):
	//	Early versions of this called me.Zero() and then simply returned on error or for incompatible types.
	//	It turns out the call to Zero() can be relatively expensive in terms of ns/op and memory allocations.
	//	We now explicitly call me.Zero only on those conditions where we are returning without actually
	//	changing me.WriteValue.
	//
	if me == nil {
		return errors.NilReceiver()
	} else if me.original == nil || !me.CanWrite || me.Kind == reflect.Invalid {
		return errors.Errorf(me.errorUnsupported("To"))
	}
	T := reflect.TypeOf(arg)
	if arg == nil || T == nil {
		return me.Zero()
	} else if (T == me.Type || T.AssignableTo(me.Type)) && me.Kind != reflect.Slice {
		// N.B: We checked that me.Kind is not a slice because this package always makes a copy of a slice!
		//
		// Performance note(s):
		//	(T == me.Type || T.AssignableTo(me.Type)) will short-circuit the call to T.AssignableTo() for
		//		basic types, further increasing performance (from 4.20% of Total down to 1.79%)
		//
		//	Early versions of this simply did:
		//		me.WriteValue.Set(reflect.ValueOf(arg))
		//		For basic built-in types this is relatively expensive, hence the type switch.
		//		Pre-bench: 		210ms within To() (9.50% of Total), 140ms in original statement.
		//		Post-bench:		50ms within To() (4.20% of Total), 10ms spread across calls to me.WriteValue.SetT()
		switch tt := arg.(type) {
		case bool:
			me.WriteValue.SetBool(tt)
		case int:
			me.WriteValue.SetInt(int64(tt))
		case int8:
			me.WriteValue.SetInt(int64(tt))
		case int16:
			me.WriteValue.SetInt(int64(tt))
		case int32:
			me.WriteValue.SetInt(int64(tt))
		case int64:
			me.WriteValue.SetInt(tt)
		case uint:
			me.WriteValue.SetUint(uint64(tt))
		case uint8:
			me.WriteValue.SetUint(uint64(tt))
		case uint16:
			me.WriteValue.SetUint(uint64(tt))
		case uint32:
			me.WriteValue.SetUint(uint64(tt))
		case uint64:
			me.WriteValue.SetUint(tt)
		case float32:
			me.WriteValue.SetFloat(float64(tt))
		case float64:
			me.WriteValue.SetFloat(tt)
		case string:
			me.WriteValue.SetString(tt)
		default:
			me.WriteValue.Set(reflect.ValueOf(arg))
		}
		return nil
	} else if me.IsScalar && T.Kind() == me.Kind && T.ConvertibleTo(me.Type) {
		// This catches scenarios where the types differ but the underlying kinds do not; for example:
		//		type NewString string
		//		var dst NewString
		//		src := "A regular string."
		//		set.V(&dst).To(src)
		me.WriteValue.Set(reflect.ValueOf(arg).Convert(me.Type))
		return nil
	}
	//
	// If arg/data represents any type of pointer we want to get to the final value:
	dataValue := reflect.ValueOf(arg)
	for ; dataValue.Kind() == reflect.Ptr; dataValue = reflect.Indirect(dataValue) {
		if dataValue.IsNil() { // If arg is a pointer and eventually nil we're done because we're already zero value.
			return me.Zero()
		}
	}
	dataTypeInfo := TypeCache.StatType(dataValue.Type())
	//
	if me.IsSlice {
		me.Zero() // Zero only returns errors on nil receiver, invalid kind, or !CanWrite -- which are already checked above.
		if !dataTypeInfo.IsSlice {
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
	} else if dataTypeInfo.Kind == reflect.Slice {
		// If the incoming type is slice but ours is not then we call set again using the last element in the slice.
		if dataValue.Len() > 0 {
			return me.To(dataValue.Index(dataValue.Len() - 1).Interface())
		}
	} else if me.IsScalar {
		if err := coerce(me.WriteValue, dataValue); err != nil {
			return errors.Go(err)
		}
		return nil
	}
	return me.Zero()
}
