package set

import (
	"reflect"

	"github.com/nofeaturesonlybugs/errors"
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
//	v := set.V(&bp)
//	// bp now contains allocated memory.
func V(arg interface{}) *Value {
	var v reflect.Value
	var t reflect.Type
	var k reflect.Kind
	//
	rv := &Value{}
	rv.original = arg
	//
	if argReflectValue, ok := arg.(reflect.Value); ok {
		v, t, k = argReflectValue, argReflectValue.Type(), argReflectValue.Kind()
	} else {
		v = reflect.ValueOf(arg)
		func() {
			defer func() {
				recover()
			}()
			// This can panic if arg == nil so wrap cleanly.
			t, k = v.Type(), v.Kind()
		}()
	}
	rv.v, rv.t, rv.k = v, t, k
	rv.pv, rv.pt, rv.pk = v, t, k
	//
	// If the incoming type is a pointer then we will follow the pointer trail to the end element;
	// along the way we will instantiate any new pointers necessary.
	func() {
		defer func() {
			recover()
		}()
		if k == reflect.Ptr {
			for k == reflect.Ptr {
				if v.IsNil() && v.CanSet() {
					ptr := reflect.New(t.Elem())
					v.Set(ptr)
					v = ptr
				}
				v = v.Elem()
				// Note that if the original arg was a nil-pointer and unsettable that we will panic
				// here; thus the anonymous function wrapping and defer-recover.
				t, k = v.Type(), v.Kind()
			}
			rv.pv, rv.pt, rv.pk = v, t, k
		}
	}()
	//
	rv.IsMap = rv.pk == reflect.Map
	rv.IsSlice = rv.pk == reflect.Slice
	rv.IsStruct = rv.pk == reflect.Struct
	rv.IsScalar = rv.pk == reflect.Bool ||
		rv.pk == reflect.Int || rv.pk == reflect.Int8 || rv.pk == reflect.Int16 || rv.pk == reflect.Int32 || rv.pk == reflect.Int64 ||
		rv.pk == reflect.Uint || rv.pk == reflect.Uint8 || rv.pk == reflect.Uint16 || rv.pk == reflect.Uint32 || rv.pk == reflect.Uint64 ||
		rv.pk == reflect.Float32 || rv.pk == reflect.Float64 ||
		rv.pk == reflect.String
		//
	if rv.IsMap || rv.IsSlice {
		ptr := reflect.New(rv.pt.Elem())
		rv.Elem = V(ptr)
	}
	//
	return rv
}

// Value wraps around a Go variable and performs magic.
type Value struct {
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

	// When IsMap or IsSlice are true then Elem is a *Value of a the zero-type contained in the map or slice.
	// Otherwise Elem is a nil pointer.
	Elem *Value

	original interface{}
	// Basic reflect information for the initial value.
	v reflect.Value
	t reflect.Type
	k reflect.Kind
	//
	// If original is a pointer then we want to know about the type pointed to.
	pv reflect.Value // Pointed-to-reflect.Value
	pt reflect.Type  // Pointed-to-reflect.Type
	pk reflect.Kind  // Pointed-to-reflect.Kind
}

// Append appends the item(s) to the end of the Value assuming it is some type of slice and every
// item can be type-coerced into the slice's data type.  Either all items are appended without an error
// or no items are appended and an error is returned describing the type of the item that could not
// be appended.
func (me *Value) Append(items ...interface{}) error {
	if !me.IsSlice {
		return nil
	} else {
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("%v", r)
				}
			}()
			zero := reflect.Zero(me.pt)
			for _, item := range items {
				elem := reflect.New(me.pt.Elem())
				elemAsValue := V(elem)
				if err = elemAsValue.To(item); err != nil {
					err = errors.Go(err)
					return
				}
				zero = reflect.Append(zero, reflect.Indirect(elemAsValue.v))
			}
			me.pv.Set(reflect.AppendSlice(me.pv, zero))
		}()
		return err
	}
}

// Fields returns a slice of Field structs when Value is wrapped around a struct; for all other values
// nil is returned.
func (me *Value) Fields() []Field {
	var rv []Field
	if me != nil && me.IsStruct {
		for k, max := 0, me.pt.NumField(); k < max; k++ {
			v, f := me.pv.Field(k), me.pt.Field(k)
			rv = append(rv, Field{Value: V(v), Field: f})
		}
	}
	return rv
}

// FieldsByTag is the same as Fields() except only Fields with the given struct-tag are returned and the
// TagValue member of Field will be set to the tag's value.
func (me *Value) FieldsByTag(key string) []Field {
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
			} else if field.Value.IsSlice && field.Value.Elem != nil && field.Value.Elem.IsStruct {
				if err = field.Value.Zero(); err != nil {
					return errors.Go(err)
				}
				elem := V(reflect.New(field.Value.Elem.pt))
				if err = fillFunc(elem, got); err != nil {
					return errors.Go(err)
				}
				if err = field.Value.Append(elem.pv.Interface()); err != nil {
					return errors.Go(err)
				}
			} else {
				return errors.Errorf("Getter.Get( %v ) returned a Getter for field %v and field is not fillable.", getName, field.Field.Name)
			}

		case []Getter:
			// What was returned from the Getter is a []Getter; therefore we expect field.Value to
			// be a []struct or struct that we can sub-fill.
			if field.Value.IsSlice && field.Value.Elem != nil && field.Value.Elem.IsStruct {
				// Zero out the existing slice.
				if err = field.Value.Zero(); err != nil {
					return errors.Go(err)
				}
				for _, elemGetter := range got {
					elem := V(reflect.New(field.Value.Elem.pt))
					if err = fillFunc(elem, elemGetter); err != nil {
						return errors.Go(err)
					}
					if err = field.Value.Append(elem.pv.Interface()); err != nil {
						return errors.Go(err)
					}
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

// Zero sets the Value to the Zero value of the appropriate type.
func (me *Value) Zero() error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("%v", r)
			}
		}()
		me.pv.Set(reflect.Zero(me.pt))
	}()
	return err
}

// To attempts to assign the argument into Value; Value is always set to the Zero value for its type before
// any other assignment ensuring if an assignment fails for any reason that any old data is overwritten.
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
	data := V(arg)
	if me.pv.CanSet() == false {
		return errors.Errorf("Set expects an assignable type; %T is not.", me.original)
	} else if me.pv.Set(reflect.Zero(me.pt)); false {
		// A non-entry if to set the type to its zero value but after checking if its settable.
	} else if data.original == nil {
		return nil
	} else if data.v.IsValid() && data.t.AssignableTo(me.pt) && me.pk != reflect.Slice {
		// N.B: We checked that me.pk is not a slice because this package always makes a copy of a slice!
		me.pv.Set(data.v)
		return nil
	}
	//
	// If arg/data represents any type of pointer we want to get to the final value:
	for ; data.k == reflect.Ptr; data = V(reflect.Indirect(data.v)) {
	}
	//
	if me.IsSlice {
		return me.setSlice(arg)
	} else if data.k == reflect.Slice {
		// If the incoming type is slice but ours is not then we call set again using the last element in the slice.
		if data.v.Len() > 0 {
			return me.To(data.v.Index(data.v.Len() - 1).Interface())
		}
	} else if err := coerce(me.pv, data.v); err != nil {
		return errors.Go(err)
	}
	return nil
}

// setSlice is a sub-feature of Set() where Value represents a slice, type irrelevant, and we
// want to zero it out and append everything passed to Set() to the new slice.
func (me *Value) setSlice(arg interface{}) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.Errorf("%v", r)
			}
		}()
		//
		// Initially we set ourselves to a Zero value.
		if err = me.Zero(); err != nil {
			return
		}
		// We make a zero-value of the underlying slice but wrapped in our Value interface so we can
		// append-with-type-coercion.
		zero := V(reflect.New(me.pt))
		//
		// Here we coerce the incoming argument to a slice if it isn't already.
		if reflect.TypeOf(arg).Kind() != reflect.Slice {
			arg = []interface{}{arg}
		}
		slice := reflect.ValueOf(arg)
		// With slice representing a reflected slice we can iterate the elements and append to the zero variable
		// which also handles the type coercions for us.
		for k, size := 0, slice.Len(); k < size; k++ {
			elem := slice.Index(k)
			if err = zero.Append(elem.Interface()); err != nil {
				return
			}
		}
		//
		// If we make it here all elements were coerced or appended and we can set our internal reflected
		// slice to the zero one we created.
		me.pv.Set(zero.pv)
	}()
	return err
}
