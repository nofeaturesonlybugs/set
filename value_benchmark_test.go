package set_test

import (
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

func BenchmarkValue(b *testing.B) {
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		*Common
		*Timestamps // Not used but present anyways
		First       string
		Last        string
	}
	type Vendor struct {
		*Common
		*Timestamps // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		*Common
		*Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	for k := 0; k < b.N; k++ {
		dest := new(T)
		set.V(dest)
	}
}

func BenchmarkScalarTo(b *testing.B) {
	type B bool
	type F32 float32
	type F64 float64
	type I8 int8
	type I16 int16
	type I32 int32
	type I64 int64
	type I int
	type U8 uint8
	type U16 uint16
	type U32 uint32
	type U64 uint64
	type U uint
	type S string

	Benchit := func(b *testing.B, v set.Value) {
		for n := 0; n < b.N; n++ {
			_ = v.To(true)
			_ = v.To(false)
			_ = v.To(B(true))
			_ = v.To(B(false))
			_ = v.To(float32(33))
			_ = v.To(float64(66))
			_ = v.To(F32(33))
			_ = v.To(F64(66))
			_ = v.To(int8(10))
			_ = v.To(int16(20))
			_ = v.To(int32(30))
			_ = v.To(int64(40))
			_ = v.To(int(50))
			_ = v.To(I8(10))
			_ = v.To(I16(20))
			_ = v.To(I32(30))
			_ = v.To(I64(40))
			_ = v.To(I(50))
			_ = v.To(uint8(10))
			_ = v.To(uint16(20))
			_ = v.To(uint32(30))
			_ = v.To(uint64(40))
			_ = v.To(uint(50))
			_ = v.To(U8(10))
			_ = v.To(U16(20))
			_ = v.To(U32(30))
			_ = v.To(U64(40))
			_ = v.To(U(50))
			if v.Kind == reflect.Bool {
				_ = v.To("1")
				_ = v.To(S("0"))
			} else {
				_ = v.To("42")
				_ = v.To(S("55"))
			}
		}
	}

	b.Run("bool", func(b *testing.B) {
		var dst bool
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("float32", func(b *testing.B) {
		var dst float32
		v := set.V(&dst)
		Benchit(b, v)

	})
	b.Run("float64", func(b *testing.B) {
		var dst float64
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("int", func(b *testing.B) {
		var dst int
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("int8", func(b *testing.B) {
		var dst int8
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("int16", func(b *testing.B) {
		var dst int16
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("int32", func(b *testing.B) {
		var dst int32
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("int64", func(b *testing.B) {
		var dst int64
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("uint", func(b *testing.B) {
		var dst uint
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("uint8", func(b *testing.B) {
		var dst uint8
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("uint16", func(b *testing.B) {
		var dst uint16
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("uint32", func(b *testing.B) {
		var dst uint32
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("uint64", func(b *testing.B) {
		var dst uint64
		v := set.V(&dst)
		Benchit(b, v)
	})
	b.Run("string", func(b *testing.B) {
		var dst string
		v := set.V(&dst)
		Benchit(b, v)
	})
}
