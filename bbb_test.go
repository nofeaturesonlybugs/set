package set_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

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

	// type Value struct {
	// 	Name  string
	// 	Value interface{}
	// }
	// values := []Value{
	// 	{Name: "bool-true", Value: true},
	// 	{Name: "bool-false", Value: false},
	// 	{Name: "B-true", Value: B(true)},
	// 	{Name: "B-false", Value: B(false)},
	// 	{Name: "float32", Value: float32(33)},
	// 	{Name: "float64", Value: float64(66)},
	// 	{Name: "F32", Value: F32(33)},
	// 	{Name: "F64", Value: F64(66)},
	// 	{Name: "int8", Value: int8(10)},
	// 	{Name: "int16", Value: int16(20)},
	// 	{Name: "int32", Value: int32(30)},
	// 	{Name: "int64", Value: int64(40)},
	// 	{Name: "int", Value: int(50)},
	// 	{Name: "I8", Value: I8(10)},
	// 	{Name: "I16", Value: I16(20)},
	// 	{Name: "I32", Value: I32(30)},
	// 	{Name: "I64", Value: I64(40)},
	// 	{Name: "I", Value: I(50)},
	// 	{Name: "uint8", Value: uint8(10)},
	// 	{Name: "uint16", Value: uint16(20)},
	// 	{Name: "uint32", Value: uint32(30)},
	// 	{Name: "uint64", Value: uint64(40)},
	// 	{Name: "uint", Value: uint(50)},
	// 	{Name: "U8", Value: U8(10)},
	// 	{Name: "U16", Value: U16(20)},
	// 	{Name: "U32", Value: U32(30)},
	// 	{Name: "U64", Value: U64(40)},
	// 	{Name: "U", Value: U(50)},
	// 	{Name: "string", Value: "42"},
	// 	{Name: "S", Value: S("55")},
	// }
	// type Dest struct {
	// 	Name  string
	// 	Value interface{}
	// }
	// var bb bool
	// var f32 float32
	// var f64 float64
	// var i int
	// var i8 int8
	// var i16 int16
	// var i32 int32
	// var i64 int64
	// var u uint
	// var u8 uint8
	// var u16 uint16
	// var u32 uint32
	// var u64 uint64
	// var s string
	// dests := []Dest{
	// 	{Name: "bool", Value: &bb},
	// 	{Name: "float32", Value: &f32},
	// 	{Name: "float64", Value: &f64},
	// 	{Name: "int", Value: &i},
	// 	{Name: "int8", Value: &i8},
	// 	{Name: "int16", Value: &i16},
	// 	{Name: "int32", Value: &i32},
	// 	{Name: "int64", Value: &i64},
	// 	{Name: "uint", Value: &u},
	// 	{Name: "uint8", Value: &u8},
	// 	{Name: "uint16", Value: &u16},
	// 	{Name: "uint32", Value: &u32},
	// 	{Name: "uint64", Value: &u64},
	// 	{Name: "string", Value: &s},
	// }
	// for _, dest := range dests {
	// 	b.Run(dest.Name, func(b *testing.B) {
	// 		v := set.V(dest.Value)
	// 		for n := 0; n < b.N; n++ {
	// 			for _, to := range values {
	// 				v.To(to)
	// 			}
	// 		}
	// 	})
	// }
	Benchit := func(b *testing.B, v *set.Value) {
		for n := 0; n < b.N; n++ {
			v.To(true)
			v.To(false)
			v.To(B(true))
			v.To(B(false))
			v.To(float32(33))
			v.To(float64(66))
			v.To(F32(33))
			v.To(F64(66))
			v.To(int8(10))
			v.To(int16(20))
			v.To(int32(30))
			v.To(int64(40))
			v.To(int(50))
			v.To(I8(10))
			v.To(I16(20))
			v.To(I32(30))
			v.To(I64(40))
			v.To(I(50))
			v.To(uint8(10))
			v.To(uint16(20))
			v.To(uint32(30))
			v.To(uint64(40))
			v.To(uint(50))
			v.To(U8(10))
			v.To(U16(20))
			v.To(U32(30))
			v.To(U64(40))
			v.To(U(50))
			v.To("42")
			v.To(S("55"))
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
