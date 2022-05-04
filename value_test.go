package set_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestValueOnNil(t *testing.T) {
	chk := assert.New(t)
	//
	{
		value := set.V(nil)
		chk.Equal(false, value.WriteValue.IsValid())
	}
	{
		var err error
		value := set.V(err)
		chk.Equal(false, value.WriteValue.IsValid())
		value = set.V(&err)
		chk.Equal(true, value.WriteValue.IsValid())
	}
	{ // Test case that pointer that is eventually nil sets slice to nil
		var pptr **bool
		ppptr := &pptr
		slice := []int{0, 1, 2, 3}
		value := set.V(&slice)
		err := value.To(ppptr)
		chk.NoError(err)
		chk.Nil(slice)
	}
}
func TestValue_fields(t *testing.T) {
	chk := assert.New(t)
	//
	{
		var b bool
		value := set.V(&b)
		chk.NotNil(value)
		fields := value.Fields()
		chk.Nil(fields)
	}
	{
		var b **bool
		value := set.V(&b)
		chk.NotNil(value)
		fields := value.Fields()
		chk.Nil(fields)
	}
	{
		// Can't set unexported fields.
		type T struct {
			a string
			b string
		}
		var t T
		value := set.V(&t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.Error(err)
		}
		chk.Equal("", t.a)
		chk.Equal("", t.b)
	}
	{
		// Address-of t not passed; can not set fields.
		type T struct {
			a string
			b string
		}
		var t T
		value := set.V(t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.Error(err)
		}
		chk.Equal("", t.a)
		chk.Equal("", t.b)
	}
	{
		// Settable.
		type T struct {
			A string
			B string
		}
		var t T
		value := set.V(&t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.NoError(err)
		}
		chk.Equal(".", t.A)
		chk.Equal(".", t.B)
	}
}

func TestValue_rebind(t *testing.T) {
	chk := assert.New(t)
	{
		var a, b string
		v := set.V(&a)
		v.To("Hello")
		v.Rebind(reflect.ValueOf(&b))
		v.To("Goodbye")
		chk.Equal("Hello", a)
		chk.Equal("Goodbye", b)
		v.Rebind(reflect.ValueOf(&a))
		v.To("Hello x2")
		chk.Equal("Hello x2", a)
	}
}

func TestValue_copy(t *testing.T) {
	chk := assert.New(t)
	//
	var b bool
	v := set.V(&b)
	err := v.To(true)
	chk.NoError(err)
	//
	v2 := v.Copy()
	chk.NotNil(v2)
	chk.True(v2.WriteValue.Bool())
}

func TestValue_set(t *testing.T) {
	chk := assert.New(t)
	//
	{ // Only addressable values can be set; passing local variable fails.
		var v bool
		err := set.V(v).To(true)
		chk.Error(err)
	}
	{ // Only addressable values can be set; passing address-of local variable works.
		var v bool
		chk.Equal(false, v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.Equal(true, v)
	}
	{ // The local variable is a pointer but its address is not passed; this fails.
		var v *bool
		chk.Equal((*bool)(nil), v)
		err := set.V(v).To(true)
		chk.Error(err)
		chk.Nil(v)
	}
	{ // The local variable is a pointer and its address is passed; a new pointer is created and assignable.
		var o, v *bool
		chk.Equal((*bool)(nil), v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		o = v // Save o so we can compare if a new pointer is created or existing is reused.
		//
		err = set.V(&v).To(false)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(false, *v)
		chk.Equal(o, v)
	}
	{ // The local variable is an instantiated pointer and we pass its address.
		v := new(bool)
		o := v
		chk.NotNil(v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		chk.Equal(o, v)
	}
	{ // The local variable is an instantiated pointer and we do not pass its address.
		v := new(bool)
		o := v
		chk.NotNil(v)
		err := set.V(v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		chk.Equal(o, v)
	}
}

func TestValueTo_FloatToFloat(t *testing.T) {
	chk := assert.New(t)
	var F32 float32
	var F64 float64
	err := set.V(&F64).To(float32(8))
	chk.NoError(err)
	chk.Equal(float64(8), F64)
	err = set.V(&F64).To(float32(16))
	chk.NoError(err)
	chk.Equal(float64(16), F64)
	err = set.V(&F64).To(float32(32))
	chk.NoError(err)
	chk.Equal(float64(32), F64)
	err = set.V(&F64).To(float64(-1))
	chk.NoError(err)
	chk.Equal(float64(-1), F64)
	//
	err = set.V(&F32).To(float64(math.MaxFloat32))
	chk.NoError(err)
	chk.Equal(float32(math.MaxFloat32), F32)
	err = set.V(&F32).To(float64(math.MaxFloat64))
	chk.Error(err)
	chk.Equal(float32(0), F32)
}

func TestValueTo_IntToInt(t *testing.T) {
	chk := assert.New(t)
	var I8 int8
	var I64 int64
	err := set.V(&I64).To(int8(8))
	chk.NoError(err)
	chk.Equal(int64(8), I64)
	err = set.V(&I64).To(int16(16))
	chk.NoError(err)
	chk.Equal(int64(16), I64)
	err = set.V(&I64).To(int32(32))
	chk.NoError(err)
	chk.Equal(int64(32), I64)
	err = set.V(&I64).To(int(-1))
	chk.NoError(err)
	chk.Equal(int64(-1), I64)
	//
	err = set.V(&I8).To(int64(math.MaxInt8))
	chk.NoError(err)
	chk.Equal(int8(math.MaxInt8), I8)
	err = set.V(&I8).To(int64(math.MaxInt64))
	chk.Error(err)
	chk.Equal(int8(0), I8)
}

func TestValueTo_UintToUint(t *testing.T) {
	chk := assert.New(t)
	var U8 uint8
	var U64 uint64
	err := set.V(&U64).To(uint8(8))
	chk.NoError(err)
	chk.Equal(uint64(8), U64)
	err = set.V(&U64).To(uint16(16))
	chk.NoError(err)
	chk.Equal(uint64(16), U64)
	err = set.V(&U64).To(uint32(32))
	chk.NoError(err)
	chk.Equal(uint64(32), U64)
	err = set.V(&U64).To(uint(math.MaxInt32))
	chk.NoError(err)
	chk.Equal(uint64(math.MaxInt32), U64)
	//
	err = set.V(&U8).To(uint64(math.MaxUint8))
	chk.NoError(err)
	chk.Equal(uint8(math.MaxUint8), U8)
	err = set.V(&U8).To(uint64(math.MaxUint64))
	chk.Error(err)
	chk.Equal(uint8(0), U8)
}

func TestValueToFast(t *testing.T) {
	chk := assert.New(t)
	var (
		B   bool
		I   int
		I8  int8
		I16 int16
		I32 int32
		I64 int64
		U   uint
		U8  uint8
		U16 uint16
		U32 uint32
		U64 uint64
		F32 float32
		F64 float64
		S   string
	)
	{
		var err error
		err = set.V(&B).To(true)
		chk.NoError(err)
		err = set.V(&I).To(int(-42))
		chk.NoError(err)
		err = set.V(&I8).To(int8(-8))
		chk.NoError(err)
		err = set.V(&I16).To(int16(-16))
		chk.NoError(err)
		err = set.V(&I32).To(int32(-32))
		chk.NoError(err)
		err = set.V(&I64).To(int64(-64))
		chk.NoError(err)
		err = set.V(&U).To(uint(42))
		chk.NoError(err)
		err = set.V(&U8).To(uint8(8))
		chk.NoError(err)
		err = set.V(&U16).To(uint16(16))
		chk.NoError(err)
		err = set.V(&U32).To(uint32(32))
		chk.NoError(err)
		err = set.V(&U64).To(uint64(64))
		chk.NoError(err)
		err = set.V(&F32).To(float32(3.14))
		chk.NoError(err)
		err = set.V(&F64).To(float64(6.28))
		chk.NoError(err)
		err = set.V(&S).To("string")
		chk.NoError(err)
		//
		chk.Equal(true, B)
		chk.Equal(-42, I)
		chk.Equal(int8(-8), I8)
		chk.Equal(int16(-16), I16)
		chk.Equal(int32(-32), I32)
		chk.Equal(int64(-64), I64)
		chk.Equal(uint(42), U)
		chk.Equal(uint8(8), U8)
		chk.Equal(uint16(16), U16)
		chk.Equal(uint32(32), U32)
		chk.Equal(uint64(64), U64)
		chk.Equal(float32(3.14), F32)
		chk.Equal(float64(6.28), F64)
		chk.Equal("string", S)
	}
}

func TestValue_setPointers(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		bp := &b
		err = set.V(bp).To("1")
		chk.NoError(err)
		chk.Equal(true, *bp)
		chk.Equal(true, b)
	}
	{
		var bp *bool
		err = set.V(&bp).To("True")
		chk.NoError(err)
		chk.Equal(true, *bp)
	}
	{
		var bppp ***bool
		err = set.V(&bppp).To("True")
		chk.NoError(err)
		chk.Equal(true, ***bppp)
	}
	{
		var ippp ***int
		s := "42"
		sp := &s
		spp := &sp
		err = set.V(&ippp).To(spp)
		chk.NoError(err)
		chk.Equal(42, ***ippp)
	}
	{
		s := "True"
		var bpp **bool
		err = set.V(&bpp).To(&s)
		chk.NoError(err)
		chk.Equal(true, **bpp)
	}
	{
		s := "True"
		var bp *bool
		err = set.V(&bp).To(&s)
		chk.NoError(err)
		chk.Equal(true, *bp)
	}
	{
		b, s := false, "True"
		bp, sp := &b, &s
		bpp, spp := &bp, &sp
		err = set.V(bpp).To(spp)
		chk.NoError(err)
		chk.Equal(true, b)
	}
}

func TestValue_setSlice(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(b).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(3, len(b))
	}
	{
		i := []int{2, 4, 6}
		chk.Equal(3, len(i))
		err = set.V(&i).To([]interface{}{"Hi"})
		chk.Error(err)
		chk.Equal(0, len(i))
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]interface{}{false, true, false, true})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]interface{}{"false", 1, 0, "True"})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To("True")
		chk.NoError(err)
		chk.Equal(1, len(b))
		chk.Equal(true, b[0])
		err = set.V(&b).To(0)
		chk.NoError(err)
		chk.Equal(1, len(b))
		chk.Equal(false, b[0])
	}
}

func TestValue_setSliceToBool(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		err = set.V(b).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(false, b)
	}
	{
		var b bool
		err = set.V(&b).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(true, b)
	}
	{
		b := true
		err = set.V(&b).To([]bool{})
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := true
		err = set.V(&b).To(([]bool)(nil))
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := true
		err = set.V(&b).To([]interface{}{true, 1, 0})
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := false
		err = set.V(&b).To([]interface{}{true, float32(0), "True"})
		chk.NoError(err)
		chk.Equal(true, b)
	}
}

func TestValue_setSliceToInt(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var i int
		err = set.V(i).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(0, i)
	}
	{
		var i int
		err = set.V(&i).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(1, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]int{})
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To(([]int)(nil))
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]interface{}{true, float32(32), float64(0)})
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]interface{}{true, float32(0), "3.14"})
		chk.NoError(err)
		chk.Equal(3, i)
	}
}

func TestValue_setSliceToString(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var s string
		err = set.V(s).To([]string{"a", "b", "c"})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal("", s)
	}
	{
		var s string
		err = set.V(&s).To([]string{"a", "b", "c"})
		chk.NoError(err)
		chk.Equal("c", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]string{})
		chk.NoError(err)
		chk.Equal("", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To(([]string)(nil))
		chk.NoError(err)
		chk.Equal("", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]interface{}{true, float32(32), float64(0), false})
		chk.NoError(err)
		chk.Equal("false", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]interface{}{true, float32(0), 64})
		chk.NoError(err)
		chk.Equal("64", s)
	}
}
func TestValue_setSliceCreatesCopies(t *testing.T) {
	chk := assert.New(t)
	//
	{
		slice := []bool{true, true, true}
		dest := []bool{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = false
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []float32{2, 4, 6}
		dest := []float32{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = 42
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []int{2, 4, 6}
		dest := []int{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = -4
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []uint{2, 4, 6}
		dest := []uint{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = 42
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []string{"Hello", "World", "foo"}
		dest := []string{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = "bar"
		chk.NotEqual(dest[1], slice[1])
	}
}

func TestValue_setStruct(t *testing.T) {
	type S struct {
		Num int
		Str string
	}
	type Wrong struct {
		Num int
		Str string
	}
	type StructTest struct {
		Name     string
		Dest     interface{}
		To       interface{}
		Error    error
		AssertFn func(interface{}, *testing.T)
	}
	tests := []StructTest{
		{
			Name: "nil",
			Dest: &S{Num: 42, Str: "Hello!"},
			To:   nil,
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(0, d.Num)
				chk.Equal("", d.Str)
			},
		},
		{
			Name: "struct",
			Dest: &S{},
			To:   S{Num: 42, Str: "Hello!"},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(42, d.Num)
				chk.Equal("Hello!", d.Str)
			},
		},
		{
			Name: "wrong struct",
			Dest: &S{Num: 10, Str: "Wrong incoming!"},
			To:   Wrong{Num: 42, Str: "Hello!"},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(0, d.Num)
				chk.Equal("", d.Str)
			},
		},
		{
			Name: "nil ptr",
			Dest: &S{Num: 42, Str: "Hello!"},
			To:   (*S)(nil),
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(0, d.Num)
				chk.Equal("", d.Str)
			},
		},
		{
			Name: "slice",
			Dest: &S{},
			To: []S{
				{Num: 1, Str: "First"},
				{Num: 2, Str: "Second"},
				{Num: 3, Str: "Third"},
			},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(3, d.Num)
				chk.Equal("Third", d.Str)
			},
		},
		{
			Name: "nil slice",
			Dest: &S{Num: 1, Str: "First"},
			To:   []S(nil),
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(0, d.Num)
				chk.Equal("", d.Str)
			},
		},
		{
			Name: "slice of ptr",
			Dest: &S{},
			To: []*S{
				{Num: 1, Str: "First"},
				{Num: 2, Str: "Second"},
				{Num: 3, Str: "Third"},
			},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(3, d.Num)
				chk.Equal("Third", d.Str)
			},
		},
		{
			Name: "slice of ptr nil element",
			Dest: &S{Num: 2, Str: "Second"},
			To: []*S{
				{Num: 1, Str: "First"},
				{Num: 2, Str: "Second"},
				nil,
			},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(0, d.Num)
				chk.Equal("", d.Str)
			},
		},
		{
			Name: "ptr to slice",
			Dest: &S{},
			To: &[]S{
				{Num: 1, Str: "First"},
				{Num: 2, Str: "Second"},
				{Num: 3, Str: "Third"},
			},
			AssertFn: func(dst interface{}, t *testing.T) {
				chk := assert.New(t)
				d := dst.(*S)
				chk.Equal(3, d.Num)
				chk.Equal("Third", d.Str)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			chk := assert.New(t)
			err := set.V(test.Dest).To(test.To)
			chk.ErrorIs(err, test.Error)
			test.AssertFn(test.Dest, t)
		})
	}
}

func TestValue_zero(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		v := 42
		value := set.V(v)
		err = value.Zero()
		chk.Error(err)
		chk.Equal(42, v)
	}
	{
		v := 42
		value := set.V(&v)
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, v)
	}
	{
		v := []int{42, -42}
		value := set.V(v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.Error(err)
		chk.Equal(2, len(v))
	}
	{
		v := []int{42, -42}
		value := set.V(&v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(v))
	}
	{
		// test zero followed by append.
		v := []int{42, -42}
		value := set.V(&v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(v))
		v = append(v, 1000)
		v = append(v, 10000, 100000)
		chk.Equal(3, len(v))
	}
	{
		type Test struct {
			S []string
		}
		instance := &Test{[]string{"Hello", "World"}}
		value := set.V(instance.S)
		chk.Equal(2, len(instance.S))
		err = value.Zero()
		chk.Error(err)
		chk.Equal(2, len(instance.S))
	}
	{
		type Test struct {
			S []string
		}
		instance := &Test{[]string{"Hello", "World"}}
		value := set.V(&instance.S)
		chk.Equal(2, len(instance.S))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(instance.S))
	}
}

func TestValue_append(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		err = set.V(&b).Append(true, false)
		chk.Error(err)
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
	}
	{
		var b []*bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, *b[0])
		chk.Equal(false, *b[1])
	}
	{
		var b []****bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, ****b[0])
		chk.Equal(false, ****b[1])
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false, "false", "1")
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false, "false", "1")
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
		// None of the following are appended.
		err = set.V(&b).Append(true, "asdf", false)
		chk.Error(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
}

func TestValue_fill(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type T struct {
		Name string
		Age  uint
	}
	type Tags struct {
		String string `key:"Name"`
		Number uint   `key:"Age"`
	}
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		default:
			return nil
		}
	})
	{
		var t T
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
	}
	{
		var t Tags
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.String)
		chk.Equal(uint(42), t.Number)
	}
}

func TestValue_fillNonStruct(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		default:
			return nil
		}
	})
	{
		var t bool
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
	}
	{
		var t int
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
	}
}

func TestValue_fillNested(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string
		Street2 string
		City    string
		State   string
		Zip     string
	}
	type Person struct {
		Name    string
		Age     uint
		Address Address
	}
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		case "Address":
			return set.GetterFunc(func(key string) interface{} {
				switch key {
				case "Street1":
					return "97531 Some Street"
				case "Street2":
					return ""
				case "City":
					return "Big City"
				case "State":
					return "ST"
				case "Zip":
					return "12345"
				default:
					return nil
				}
			})
		default:
			return nil
		}
	})
	{
		var t Person
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedByTag(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "name":
			return "Bob"
		case "age":
			return "42"
		case "address":
			return set.GetterFunc(func(key string) interface{} {
				switch key {
				case "street1":
					return "97531 Some Street"
				case "street2":
					return ""
				case "city":
					return "Big City"
				case "state":
					return "ST"
				case "zip":
					return "12345"
				default:
					return nil
				}
			})
		default:
			return nil
		}
	})
	{
		var t Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedByMap(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	m := map[string]interface{}{
		"name": "Bob",
		"age":  42,
		"address": map[interface{}]string{
			"street1": "97531 Some Street",
			"street2": "",
			"city":    "Big City",
			"state":   "ST",
			"zip":     "12345",
		},
	}
	getter := set.MapGetter(m)
	{
		var t Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedStructSlices(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	type Company struct {
		Name         string   `key:"name"`
		Employees    []Person `key:"employees"`
		LastEmployee Person   `key:"employees"`
		Slice        []Person `key:"slice"`
	}
	m := map[string]interface{}{
		"name": "Some Company",
		"slice": map[string]interface{}{
			"name": "Slice",
			"age":  2,
			"address": map[interface{}]string{
				"street1": "Slice Street",
				"street2": "",
				"city":    "Slice City",
				"state":   "SL",
				"zip":     "99999",
			},
		},
		"employees": []map[string]interface{}{
			{
				"name": "Bob",
				"age":  42,
				"address": map[interface{}]string{
					"street1": "97531 Some Street",
					"street2": "",
					"city":    "Big City",
					"state":   "ST",
					"zip":     "12345",
				},
			},
			{
				"name": "Sally",
				"age":  48,
				"address": map[interface{}]string{
					"street1": "555 Small Lane",
					"street2": "",
					"city":    "Other City",
					"state":   "OO",
					"zip":     "54321",
				},
			},
		},
	}
	getter := set.MapGetter(m)
	{
		var t Company
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		// chk.Equal("Some Company", t.Name)
		// //
		// chk.Equal(2, len(t.Employees))
		// //
		// chk.Equal("Bob", t.Employees[0].Name)
		// chk.Equal(uint(42), t.Employees[0].Age)
		// chk.Equal("97531 Some Street", t.Employees[0].Address.Street1)
		// chk.Equal("", t.Employees[0].Address.Street2)
		// chk.Equal("Big City", t.Employees[0].Address.City)
		// chk.Equal("ST", t.Employees[0].Address.State)
		// chk.Equal("12345", t.Employees[0].Address.Zip)
		// //
		// chk.Equal("Sally", t.Employees[1].Name)
		// chk.Equal(uint(48), t.Employees[1].Age)
		// chk.Equal("555 Small Lane", t.Employees[1].Address.Street1)
		// chk.Equal("", t.Employees[1].Address.Street2)
		// chk.Equal("Other City", t.Employees[1].Address.City)
		// chk.Equal("OO", t.Employees[1].Address.State)
		// chk.Equal("54321", t.Employees[1].Address.Zip)
		// //
		// chk.Equal("Sally", t.LastEmployee.Name)
		// chk.Equal(uint(48), t.LastEmployee.Age)
		// chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
		// chk.Equal("", t.LastEmployee.Address.Street2)
		// chk.Equal("Other City", t.LastEmployee.Address.City)
		// chk.Equal("OO", t.LastEmployee.Address.State)
		// chk.Equal("54321", t.LastEmployee.Address.Zip)
		// //
		// chk.Equal(1, len(t.Slice))
		// chk.Equal("Slice", t.Slice[0].Name)
		// chk.Equal("Slice Street", t.Slice[0].Address.Street1)
		// chk.Equal("", t.Slice[0].Address.Street2)
		// chk.Equal("Slice City", t.Slice[0].Address.City)
		// chk.Equal("SL", t.Slice[0].Address.State)
		// chk.Equal("99999", t.Slice[0].Address.Zip)
	}
}

func TestValue_fillNestedStructSlicesAsPointers(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string   `key:"name"`
		Age     uint     `key:"age"`
		Address *Address `key:"address"`
	}
	type Company struct {
		Name         string    `key:"name"`
		Employees    []*Person `key:"employees"`
		LastEmployee *Person   `key:"employees"`
		Slice        []*Person `key:"slice"`
	}
	m := map[string]interface{}{
		"name": "Some Company",
		"slice": map[string]interface{}{
			"name": "Slice",
			"age":  2,
			"address": map[interface{}]string{
				"street1": "Slice Street",
				"street2": "",
				"city":    "Slice City",
				"state":   "SL",
				"zip":     "99999",
			},
		},
		"employees": []map[string]interface{}{
			{
				"name": "Bob",
				"age":  42,
				"address": map[interface{}]string{
					"street1": "97531 Some Street",
					"street2": "",
					"city":    "Big City",
					"state":   "ST",
					"zip":     "12345",
				},
			},
			{
				"name": "Sally",
				"age":  48,
				"address": map[interface{}]string{
					"street1": "555 Small Lane",
					"street2": "",
					"city":    "Other City",
					"state":   "OO",
					"zip":     "54321",
				},
			},
		},
	}
	getter := set.MapGetter(m)
	{
		var t *Company
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Some Company", t.Name)
		//
		chk.Equal(2, len(t.Employees))
		//
		chk.Equal("Bob", t.Employees[0].Name)
		chk.Equal(uint(42), t.Employees[0].Age)
		chk.Equal("97531 Some Street", t.Employees[0].Address.Street1)
		chk.Equal("", t.Employees[0].Address.Street2)
		chk.Equal("Big City", t.Employees[0].Address.City)
		chk.Equal("ST", t.Employees[0].Address.State)
		chk.Equal("12345", t.Employees[0].Address.Zip)
		//
		chk.Equal("Sally", t.Employees[1].Name)
		chk.Equal(uint(48), t.Employees[1].Age)
		chk.Equal("555 Small Lane", t.Employees[1].Address.Street1)
		chk.Equal("", t.Employees[1].Address.Street2)
		chk.Equal("Other City", t.Employees[1].Address.City)
		chk.Equal("OO", t.Employees[1].Address.State)
		chk.Equal("54321", t.Employees[1].Address.Zip)
		//
		chk.Equal("Sally", t.LastEmployee.Name)
		chk.Equal(uint(48), t.LastEmployee.Age)
		chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
		chk.Equal("", t.LastEmployee.Address.Street2)
		chk.Equal("Other City", t.LastEmployee.Address.City)
		chk.Equal("OO", t.LastEmployee.Address.State)
		chk.Equal("54321", t.LastEmployee.Address.Zip)
		//
		chk.Equal(1, len(t.Slice))
		chk.Equal("Slice", t.Slice[0].Name)
		chk.Equal("Slice Street", t.Slice[0].Address.Street1)
		chk.Equal("", t.Slice[0].Address.Street2)
		chk.Equal("Slice City", t.Slice[0].Address.City)
		chk.Equal("SL", t.Slice[0].Address.State)
		chk.Equal("99999", t.Slice[0].Address.Zip)
	}
}

func TestValue_fillNestedStructPointerToSlicesAsPointers(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string   `key:"name"`
		Age     uint     `key:"age"`
		Address *Address `key:"address"`
	}
	type Company struct {
		Name         string     `key:"name"`
		Employees    *[]*Person `key:"employees"`
		LastEmployee *Person    `key:"employees"`
		Slice        *[]*Person `key:"slice"`
	}
	m := map[string]interface{}{
		"name": "Some Company",
		"slice": map[string]interface{}{
			"name": "Slice",
			"age":  2,
			"address": map[interface{}]string{
				"street1": "Slice Street",
				"street2": "",
				"city":    "Slice City",
				"state":   "SL",
				"zip":     "99999",
			},
		},
		"employees": []map[string]interface{}{
			{
				"name": "Bob",
				"age":  42,
				"address": map[interface{}]string{
					"street1": "97531 Some Street",
					"street2": "",
					"city":    "Big City",
					"state":   "ST",
					"zip":     "12345",
				},
			},
			{
				"name": "Sally",
				"age":  48,
				"address": map[interface{}]string{
					"street1": "555 Small Lane",
					"street2": "",
					"city":    "Other City",
					"state":   "OO",
					"zip":     "54321",
				},
			},
		},
	}
	getter := set.MapGetter(m)
	{
		var t *Company
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Some Company", t.Name)
		//
		chk.Equal(2, len(*t.Employees))
		//
		chk.Equal("Bob", (*t.Employees)[0].Name)
		chk.Equal(uint(42), (*t.Employees)[0].Age)
		chk.Equal("97531 Some Street", (*t.Employees)[0].Address.Street1)
		chk.Equal("", (*t.Employees)[0].Address.Street2)
		chk.Equal("Big City", (*t.Employees)[0].Address.City)
		chk.Equal("ST", (*t.Employees)[0].Address.State)
		chk.Equal("12345", (*t.Employees)[0].Address.Zip)
		//
		chk.Equal("Sally", (*t.Employees)[1].Name)
		chk.Equal(uint(48), (*t.Employees)[1].Age)
		chk.Equal("555 Small Lane", (*t.Employees)[1].Address.Street1)
		chk.Equal("", (*t.Employees)[1].Address.Street2)
		chk.Equal("Other City", (*t.Employees)[1].Address.City)
		chk.Equal("OO", (*t.Employees)[1].Address.State)
		chk.Equal("54321", (*t.Employees)[1].Address.Zip)
		//
		chk.Equal("Sally", t.LastEmployee.Name)
		chk.Equal(uint(48), t.LastEmployee.Age)
		chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
		chk.Equal("", t.LastEmployee.Address.Street2)
		chk.Equal("Other City", t.LastEmployee.Address.City)
		chk.Equal("OO", t.LastEmployee.Address.State)
		chk.Equal("54321", t.LastEmployee.Address.Zip)
		//
		chk.Equal(1, len(*t.Slice))
		chk.Equal("Slice", (*t.Slice)[0].Name)
		chk.Equal("Slice Street", (*t.Slice)[0].Address.Street1)
		chk.Equal("", (*t.Slice)[0].Address.Street2)
		chk.Equal("Slice City", (*t.Slice)[0].Address.City)
		chk.Equal("SL", (*t.Slice)[0].Address.State)
		chk.Equal("99999", (*t.Slice)[0].Address.Zip)
	}
}

func TestValue_fillNestedPointersByMap(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string   `key:"name"`
		Age     uint     `key:"age"`
		Address *Address `key:"address"`
	}
	m := map[string]interface{}{
		"name": "Bob",
		"age":  42,
		"address": map[interface{}]string{
			"street1": "97531 Some Street",
			"street2": "",
			"city":    "Big City",
			"state":   "ST",
			"zip":     "12345",
		},
	}
	getter := set.MapGetter(m)
	{
		var t *Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.NotNil(t.Address)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedPointersByMapWithNils(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string   `key:"name"`
		Age     uint     `key:"age"`
		Address *Address `key:"address"`
	}
	m := map[string]interface{}{
		"name": "Bob",
		"age":  42,
	}
	getter := set.MapGetter(m)
	{
		var t *Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.NotNil(t.Address)
	}
}

func TestValue_fieldByIndex(t *testing.T) {
	type Foo struct {
		Str string
		Num int
	}
	type Bar struct {
		A Foo
		B Foo
	}
	type Ptr struct {
		B **Bar
	}
	//
	type IndexValue struct {
		Index       []int
		Value       interface{}
		Expect      interface{}
		IndexError  error
		AssignError error
	}
	type FieldByIndexTest struct {
		Name       string
		Dest       interface{}
		Index      []int
		IndexError error
		To         interface{}
		Expect     interface{}
		ToError    error
	}
	//
	tests := []FieldByIndexTest{
		{
			Name:   "bar",
			Dest:   &Bar{},
			Index:  []int{0, 0}, // A.Str
			To:     10,
			Expect: "10",
		},
		{
			Name:   "bar",
			Dest:   &Bar{},
			Index:  []int{0, 1}, // A.Num
			To:     "-10",
			Expect: -10,
		},
		{
			Name:   "bar",
			Dest:   &Bar{},
			Index:  []int{1, 0}, // B.Str
			To:     20,
			Expect: "20",
		},
		{
			Name:   "bar",
			Dest:   &Bar{},
			Index:  []int{1, 1}, // B.Num
			To:     "-20",
			Expect: -20,
		},
		// ptr
		{
			Name:   "ptr 0,0,0",
			Dest:   &Ptr{},
			Index:  []int{0, 0, 0},
			To:     -20,
			Expect: "-20",
		},
		{
			Name:   "ptr 0,0,1",
			Dest:   &Ptr{},
			Index:  []int{0, 0, 1},
			To:     "-20",
			Expect: -20,
		},
		{
			Name:   "ptr 0,1,0",
			Dest:   &Ptr{},
			Index:  []int{0, 1, 0},
			To:     -20,
			Expect: "-20",
		},
		{
			Name:   "ptr 0,1,1",
			Dest:   &Ptr{},
			Index:  []int{0, 1, 1},
			To:     "-20",
			Expect: -20,
		},
	}
	for _, test := range tests {
		t.Run(test.Name+" "+fmt.Sprintf("%v", test.Index), func(t *testing.T) {
			chk := assert.New(t)
			v := set.V(test.Dest)
			f, err := v.FieldByIndex(test.Index)
			chk.ErrorIs(err, test.IndexError)
			if test.IndexError == nil {
				fv := set.V(f)
				err = fv.To(test.To)
				chk.ErrorIs(err, test.ToError)
				chk.Equal(test.Expect, fv.WriteValue.Interface())
			}
		})
	}
}

func TestValue_fieldByIndex_outOfRange(t *testing.T) {
	chk := assert.New(t)
	type S struct {
		A int
	}
	var s S
	_, err := set.V(&s).FieldByIndex([]int{0})
	chk.NoError(err)
	_, err = set.V(&s).FieldByIndex([]int{1})
	chk.ErrorIs(err, set.ErrIndexOutOfBounds)
}

func TestValue_fieldByIndexCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	var err error
	var value, field *set.Value
	type A struct {
		A string
	}
	type B struct {
		B string
		A
	}
	var a A

	value = set.V(map[string]string{})
	field, err = value.FieldByIndexAsValue(nil)
	chk.Error(err)
	chk.Nil(field)
	//
	value = set.V(a)
	field, err = value.FieldByIndexAsValue([]int{0})
	chk.Error(err)
	chk.Nil(field)
	//
	value = set.V(&a)
	field, err = value.FieldByIndexAsValue(nil)
	chk.Error(err)
	chk.Nil(field)
	field, err = value.FieldByIndexAsValue([]int{})
	chk.Error(err)
	chk.Nil(field)
	//
	{ // Test scalar, aka something not indexable.
		var b bool
		value = set.V(&b)
		{ // When reflect.Value is returned
			field, err := value.FieldByIndex([]int{1, 2})
			chk.Error(err)
			chk.Equal(true, field.IsValid())
		}
		{ // When *Value is returned
			field, err := value.FieldByIndexAsValue([]int{1, 2})
			chk.Error(err)
			chk.Nil(field)
		}
	}
}

func TestValue_fillCodeCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	m := map[string]interface{}{
		"Nested": map[string]interface{}{
			"String": "Hello, World!",
		},
		"Slice": []map[string]interface{}{
			{
				"String": "Hello, World!",
			},
			{
				"String": "Goodbye, World!",
			},
		},
	}
	getter := set.MapGetter(m)
	{
		type T struct {
			Nested struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested []struct {
				String int
			}
		}
		var t T
		err = set.V(t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested []struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested string
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice []struct {
				String int
			}
		}
		var t T
		err = set.V(t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice []struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice int
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
}

func TestValue_appendCodeCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b []bool
		err = set.V(b).Append(42)
		chk.Error(err)
	}
}

func TestValue_newElemCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	//
	{ // Tests NewElem when *Value is not nil but not a map
		var b bool
		v := set.V(&b)
		elem, err := v.NewElem()
		chk.Error(err)
		chk.Nil(elem)
	}
}
