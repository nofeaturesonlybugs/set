package coerce_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set/coerce"
)

func BenchmarkCoerce(b *testing.B) {
	b.Run("bool", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Bool(true)
			coerce.Bool(false)
			coerce.Bool(B(true))
			coerce.Bool(B(false))
			coerce.Bool(float32(33))
			coerce.Bool(float64(66))
			coerce.Bool(F32(33))
			coerce.Bool(F64(66))
			coerce.Bool(int8(10))
			coerce.Bool(int16(20))
			coerce.Bool(int32(30))
			coerce.Bool(int64(40))
			coerce.Bool(int(50))
			coerce.Bool(I8(10))
			coerce.Bool(I16(20))
			coerce.Bool(I32(30))
			coerce.Bool(I64(40))
			coerce.Bool(II(50))
			coerce.Bool(uint8(10))
			coerce.Bool(uint16(20))
			coerce.Bool(uint32(30))
			coerce.Bool(uint64(40))
			coerce.Bool(uint(50))
			coerce.Bool(U8(10))
			coerce.Bool(U16(20))
			coerce.Bool(U32(30))
			coerce.Bool(U64(40))
			coerce.Bool(UU(50))
			coerce.Bool("1")
			coerce.Bool(S("0"))
		}
	})
	b.Run("float32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Float32(true)
			coerce.Float32(false)
			coerce.Float32(B(true))
			coerce.Float32(B(false))
			coerce.Float32(float32(33))
			coerce.Float32(float64(66))
			coerce.Float32(F32(33))
			coerce.Float32(F64(66))
			coerce.Float32(int8(10))
			coerce.Float32(int16(20))
			coerce.Float32(int32(30))
			coerce.Float32(int64(40))
			coerce.Float32(int(50))
			coerce.Float32(I8(10))
			coerce.Float32(I16(20))
			coerce.Float32(I32(30))
			coerce.Float32(I64(40))
			coerce.Float32(II(50))
			coerce.Float32(uint8(10))
			coerce.Float32(uint16(20))
			coerce.Float32(uint32(30))
			coerce.Float32(uint64(40))
			coerce.Float32(uint(50))
			coerce.Float32(U8(10))
			coerce.Float32(U16(20))
			coerce.Float32(U32(30))
			coerce.Float32(U64(40))
			coerce.Float32(UU(50))
			coerce.Float32("42")
			coerce.Float32(S("55"))
		}
	})
	b.Run("float64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Float64(true)
			coerce.Float64(false)
			coerce.Float64(B(true))
			coerce.Float64(B(false))
			coerce.Float64(float32(33))
			coerce.Float64(float64(66))
			coerce.Float64(F32(33))
			coerce.Float64(F64(66))
			coerce.Float64(int8(10))
			coerce.Float64(int16(20))
			coerce.Float64(int32(30))
			coerce.Float64(int64(40))
			coerce.Float64(int(50))
			coerce.Float64(I8(10))
			coerce.Float64(I16(20))
			coerce.Float64(I32(30))
			coerce.Float64(I64(40))
			coerce.Float64(II(50))
			coerce.Float64(uint8(10))
			coerce.Float64(uint16(20))
			coerce.Float64(uint32(30))
			coerce.Float64(uint64(40))
			coerce.Float64(uint(50))
			coerce.Float64(U8(10))
			coerce.Float64(U16(20))
			coerce.Float64(U32(30))
			coerce.Float64(U64(40))
			coerce.Float64(UU(50))
			coerce.Float64("42")
			coerce.Float64(S("55"))
		}
	})
	b.Run("iint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Int(true)
			coerce.Int(false)
			coerce.Int(B(true))
			coerce.Int(B(false))
			coerce.Int(float32(33))
			coerce.Int(float64(66))
			coerce.Int(F32(33))
			coerce.Int(F64(66))
			coerce.Int(int8(10))
			coerce.Int(int16(20))
			coerce.Int(int32(30))
			coerce.Int(int64(40))
			coerce.Int(int(50))
			coerce.Int(I8(10))
			coerce.Int(I16(20))
			coerce.Int(I32(30))
			coerce.Int(I64(40))
			coerce.Int(II(50))
			coerce.Int(uint8(10))
			coerce.Int(uint16(20))
			coerce.Int(uint32(30))
			coerce.Int(uint64(40))
			coerce.Int(uint(50))
			coerce.Int(U8(10))
			coerce.Int(U16(20))
			coerce.Int(U32(30))
			coerce.Int(U64(40))
			coerce.Int(UU(50))
			coerce.Int("42")
			coerce.Int(S("55"))
		}
	})
	b.Run("int8", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Int8(true)
			coerce.Int8(false)
			coerce.Int8(B(true))
			coerce.Int8(B(false))
			coerce.Int8(float32(33))
			coerce.Int8(float64(66))
			coerce.Int8(F32(33))
			coerce.Int8(F64(66))
			coerce.Int8(int8(10))
			coerce.Int8(int16(20))
			coerce.Int8(int32(30))
			coerce.Int8(int64(40))
			coerce.Int8(int(50))
			coerce.Int8(I8(10))
			coerce.Int8(I16(20))
			coerce.Int8(I32(30))
			coerce.Int8(I64(40))
			coerce.Int8(II(50))
			coerce.Int8(uint8(10))
			coerce.Int8(uint16(20))
			coerce.Int8(uint32(30))
			coerce.Int8(uint64(40))
			coerce.Int8(uint(50))
			coerce.Int8(U8(10))
			coerce.Int8(U16(20))
			coerce.Int8(U32(30))
			coerce.Int8(U64(40))
			coerce.Int8(UU(50))
			coerce.Int8("42")
			coerce.Int8(S("55"))
		}
	})
	b.Run("int16", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Int16(true)
			coerce.Int16(false)
			coerce.Int16(B(true))
			coerce.Int16(B(false))
			coerce.Int16(float32(33))
			coerce.Int16(float64(66))
			coerce.Int16(F32(33))
			coerce.Int16(F64(66))
			coerce.Int16(int8(10))
			coerce.Int16(int16(20))
			coerce.Int16(int32(30))
			coerce.Int16(int64(40))
			coerce.Int16(int(50))
			coerce.Int16(I8(10))
			coerce.Int16(I16(20))
			coerce.Int16(I32(30))
			coerce.Int16(I64(40))
			coerce.Int16(II(50))
			coerce.Int16(uint8(10))
			coerce.Int16(uint16(20))
			coerce.Int16(uint32(30))
			coerce.Int16(uint64(40))
			coerce.Int16(uint(50))
			coerce.Int16(U8(10))
			coerce.Int16(U16(20))
			coerce.Int16(U32(30))
			coerce.Int16(U64(40))
			coerce.Int16(UU(50))
			coerce.Int16("42")
			coerce.Int16(S("55"))
		}
	})
	b.Run("int32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Int32(true)
			coerce.Int32(false)
			coerce.Int32(B(true))
			coerce.Int32(B(false))
			coerce.Int32(float32(33))
			coerce.Int32(float64(66))
			coerce.Int32(F32(33))
			coerce.Int32(F64(66))
			coerce.Int32(int8(10))
			coerce.Int32(int16(20))
			coerce.Int32(int32(30))
			coerce.Int32(int64(40))
			coerce.Int32(int(50))
			coerce.Int32(I8(10))
			coerce.Int32(I16(20))
			coerce.Int32(I32(30))
			coerce.Int32(I64(40))
			coerce.Int32(II(50))
			coerce.Int32(uint8(10))
			coerce.Int32(uint16(20))
			coerce.Int32(uint32(30))
			coerce.Int32(uint64(40))
			coerce.Int32(uint(50))
			coerce.Int32(U8(10))
			coerce.Int32(U16(20))
			coerce.Int32(U32(30))
			coerce.Int32(U64(40))
			coerce.Int32(UU(50))
			coerce.Int32("42")
			coerce.Int32(S("55"))
		}
	})
	b.Run("int64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Int64(true)
			coerce.Int64(false)
			coerce.Int64(B(true))
			coerce.Int64(B(false))
			coerce.Int64(float32(33))
			coerce.Int64(float64(66))
			coerce.Int64(F32(33))
			coerce.Int64(F64(66))
			coerce.Int64(int8(10))
			coerce.Int64(int16(20))
			coerce.Int64(int32(30))
			coerce.Int64(int64(40))
			coerce.Int64(int(50))
			coerce.Int64(I8(10))
			coerce.Int64(I16(20))
			coerce.Int64(I32(30))
			coerce.Int64(I64(40))
			coerce.Int64(II(50))
			coerce.Int64(uint8(10))
			coerce.Int64(uint16(20))
			coerce.Int64(uint32(30))
			coerce.Int64(uint64(40))
			coerce.Int64(uint(50))
			coerce.Int64(U8(10))
			coerce.Int64(U16(20))
			coerce.Int64(U32(30))
			coerce.Int64(U64(40))
			coerce.Int64(UU(50))
			coerce.Int64("42")
			coerce.Int64(S("55"))
		}
	})
	b.Run("uuint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Uint(true)
			coerce.Uint(false)
			coerce.Uint(B(true))
			coerce.Uint(B(false))
			coerce.Uint(float32(33))
			coerce.Uint(float64(66))
			coerce.Uint(F32(33))
			coerce.Uint(F64(66))
			coerce.Uint(int8(10))
			coerce.Uint(int16(20))
			coerce.Uint(int32(30))
			coerce.Uint(int64(40))
			coerce.Uint(int(50))
			coerce.Uint(I8(10))
			coerce.Uint(I16(20))
			coerce.Uint(I32(30))
			coerce.Uint(I64(40))
			coerce.Uint(II(50))
			coerce.Uint(uint8(10))
			coerce.Uint(uint16(20))
			coerce.Uint(uint32(30))
			coerce.Uint(uint64(40))
			coerce.Uint(uint(50))
			coerce.Uint(U8(10))
			coerce.Uint(U16(20))
			coerce.Uint(U32(30))
			coerce.Uint(U64(40))
			coerce.Uint(UU(50))
			coerce.Uint("42")
			coerce.Uint(S("55"))
		}
	})
	b.Run("uint8", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Uint8(true)
			coerce.Uint8(false)
			coerce.Uint8(B(true))
			coerce.Uint8(B(false))
			coerce.Uint8(float32(33))
			coerce.Uint8(float64(66))
			coerce.Uint8(F32(33))
			coerce.Uint8(F64(66))
			coerce.Uint8(int8(10))
			coerce.Uint8(int16(20))
			coerce.Uint8(int32(30))
			coerce.Uint8(int64(40))
			coerce.Uint8(int(50))
			coerce.Uint8(I8(10))
			coerce.Uint8(I16(20))
			coerce.Uint8(I32(30))
			coerce.Uint8(I64(40))
			coerce.Uint8(II(50))
			coerce.Uint8(uint8(10))
			coerce.Uint8(uint16(20))
			coerce.Uint8(uint32(30))
			coerce.Uint8(uint64(40))
			coerce.Uint8(uint(50))
			coerce.Uint8(U8(10))
			coerce.Uint8(U16(20))
			coerce.Uint8(U32(30))
			coerce.Uint8(U64(40))
			coerce.Uint8(UU(50))
			coerce.Uint8("42")
			coerce.Uint8(S("55"))
		}
	})
	b.Run("uint16", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Uint16(true)
			coerce.Uint16(false)
			coerce.Uint16(B(true))
			coerce.Uint16(B(false))
			coerce.Uint16(float32(33))
			coerce.Uint16(float64(66))
			coerce.Uint16(F32(33))
			coerce.Uint16(F64(66))
			coerce.Uint16(int8(10))
			coerce.Uint16(int16(20))
			coerce.Uint16(int32(30))
			coerce.Uint16(int64(40))
			coerce.Uint16(int(50))
			coerce.Uint16(I8(10))
			coerce.Uint16(I16(20))
			coerce.Uint16(I32(30))
			coerce.Uint16(I64(40))
			coerce.Uint16(II(50))
			coerce.Uint16(uint8(10))
			coerce.Uint16(uint16(20))
			coerce.Uint16(uint32(30))
			coerce.Uint16(uint64(40))
			coerce.Uint16(uint(50))
			coerce.Uint16(U8(10))
			coerce.Uint16(U16(20))
			coerce.Uint16(U32(30))
			coerce.Uint16(U64(40))
			coerce.Uint16(UU(50))
			coerce.Uint16("42")
			coerce.Uint16(S("55"))
		}
	})
	b.Run("uint32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Uint32(true)
			coerce.Uint32(false)
			coerce.Uint32(B(true))
			coerce.Uint32(B(false))
			coerce.Uint32(float32(33))
			coerce.Uint32(float64(66))
			coerce.Uint32(F32(33))
			coerce.Uint32(F64(66))
			coerce.Uint32(int8(10))
			coerce.Uint32(int16(20))
			coerce.Uint32(int32(30))
			coerce.Uint32(int64(40))
			coerce.Uint32(int(50))
			coerce.Uint32(I8(10))
			coerce.Uint32(I16(20))
			coerce.Uint32(I32(30))
			coerce.Uint32(I64(40))
			coerce.Uint32(II(50))
			coerce.Uint32(uint8(10))
			coerce.Uint32(uint16(20))
			coerce.Uint32(uint32(30))
			coerce.Uint32(uint64(40))
			coerce.Uint32(uint(50))
			coerce.Uint32(U8(10))
			coerce.Uint32(U16(20))
			coerce.Uint32(U32(30))
			coerce.Uint32(U64(40))
			coerce.Uint32(UU(50))
			coerce.Uint32("42")
			coerce.Uint32(S("55"))
		}
	})
	b.Run("uint64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.Uint64(true)
			coerce.Uint64(false)
			coerce.Uint64(B(true))
			coerce.Uint64(B(false))
			coerce.Uint64(float32(33))
			coerce.Uint64(float64(66))
			coerce.Uint64(F32(33))
			coerce.Uint64(F64(66))
			coerce.Uint64(int8(10))
			coerce.Uint64(int16(20))
			coerce.Uint64(int32(30))
			coerce.Uint64(int64(40))
			coerce.Uint64(int(50))
			coerce.Uint64(I8(10))
			coerce.Uint64(I16(20))
			coerce.Uint64(I32(30))
			coerce.Uint64(I64(40))
			coerce.Uint64(II(50))
			coerce.Uint64(uint8(10))
			coerce.Uint64(uint16(20))
			coerce.Uint64(uint32(30))
			coerce.Uint64(uint64(40))
			coerce.Uint64(uint(50))
			coerce.Uint64(U8(10))
			coerce.Uint64(U16(20))
			coerce.Uint64(U32(30))
			coerce.Uint64(U64(40))
			coerce.Uint64(UU(50))
			coerce.Uint64("42")
			coerce.Uint64(S("55"))
		}
	})
	b.Run("string", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			coerce.String(true)
			coerce.String(false)
			coerce.String(B(true))
			coerce.String(B(false))
			coerce.String(float32(33))
			coerce.String(float64(66))
			coerce.String(F32(33))
			coerce.String(F64(66))
			coerce.String(int8(10))
			coerce.String(int16(20))
			coerce.String(int32(30))
			coerce.String(int64(40))
			coerce.String(int(50))
			coerce.String(I8(10))
			coerce.String(I16(20))
			coerce.String(I32(30))
			coerce.String(I64(40))
			coerce.String(II(50))
			coerce.String(uint8(10))
			coerce.String(uint16(20))
			coerce.String(uint32(30))
			coerce.String(uint64(40))
			coerce.String(uint(50))
			coerce.String(U8(10))
			coerce.String(U16(20))
			coerce.String(U32(30))
			coerce.String(U64(40))
			coerce.String(UU(50))
			coerce.String("42")
			coerce.String(S("55"))
		}
	})
}