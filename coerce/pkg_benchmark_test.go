package coerce_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set/coerce"
)

func BenchmarkCoerce(b *testing.B) {
	b.Run("bool", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Bool(true)
			_, _ = coerce.Bool(false)
			_, _ = coerce.Bool(B(true))
			_, _ = coerce.Bool(B(false))
			_, _ = coerce.Bool(float32(33))
			_, _ = coerce.Bool(float64(66))
			_, _ = coerce.Bool(F32(33))
			_, _ = coerce.Bool(F64(66))
			_, _ = coerce.Bool(int8(10))
			_, _ = coerce.Bool(int16(20))
			_, _ = coerce.Bool(int32(30))
			_, _ = coerce.Bool(int64(40))
			_, _ = coerce.Bool(int(50))
			_, _ = coerce.Bool(I8(10))
			_, _ = coerce.Bool(I16(20))
			_, _ = coerce.Bool(I32(30))
			_, _ = coerce.Bool(I64(40))
			_, _ = coerce.Bool(II(50))
			_, _ = coerce.Bool(uint8(10))
			_, _ = coerce.Bool(uint16(20))
			_, _ = coerce.Bool(uint32(30))
			_, _ = coerce.Bool(uint64(40))
			_, _ = coerce.Bool(uint(50))
			_, _ = coerce.Bool(U8(10))
			_, _ = coerce.Bool(U16(20))
			_, _ = coerce.Bool(U32(30))
			_, _ = coerce.Bool(U64(40))
			_, _ = coerce.Bool(UU(50))
			_, _ = coerce.Bool("1")
			_, _ = coerce.Bool(S("0"))
		}
	})
	b.Run("float32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Float32(true)
			_, _ = coerce.Float32(false)
			_, _ = coerce.Float32(B(true))
			_, _ = coerce.Float32(B(false))
			_, _ = coerce.Float32(float32(33))
			_, _ = coerce.Float32(float64(66))
			_, _ = coerce.Float32(F32(33))
			_, _ = coerce.Float32(F64(66))
			_, _ = coerce.Float32(int8(10))
			_, _ = coerce.Float32(int16(20))
			_, _ = coerce.Float32(int32(30))
			_, _ = coerce.Float32(int64(40))
			_, _ = coerce.Float32(int(50))
			_, _ = coerce.Float32(I8(10))
			_, _ = coerce.Float32(I16(20))
			_, _ = coerce.Float32(I32(30))
			_, _ = coerce.Float32(I64(40))
			_, _ = coerce.Float32(II(50))
			_, _ = coerce.Float32(uint8(10))
			_, _ = coerce.Float32(uint16(20))
			_, _ = coerce.Float32(uint32(30))
			_, _ = coerce.Float32(uint64(40))
			_, _ = coerce.Float32(uint(50))
			_, _ = coerce.Float32(U8(10))
			_, _ = coerce.Float32(U16(20))
			_, _ = coerce.Float32(U32(30))
			_, _ = coerce.Float32(U64(40))
			_, _ = coerce.Float32(UU(50))
			_, _ = coerce.Float32("42")
			_, _ = coerce.Float32(S("55"))
		}
	})
	b.Run("float64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Float64(true)
			_, _ = coerce.Float64(false)
			_, _ = coerce.Float64(B(true))
			_, _ = coerce.Float64(B(false))
			_, _ = coerce.Float64(float32(33))
			_, _ = coerce.Float64(float64(66))
			_, _ = coerce.Float64(F32(33))
			_, _ = coerce.Float64(F64(66))
			_, _ = coerce.Float64(int8(10))
			_, _ = coerce.Float64(int16(20))
			_, _ = coerce.Float64(int32(30))
			_, _ = coerce.Float64(int64(40))
			_, _ = coerce.Float64(int(50))
			_, _ = coerce.Float64(I8(10))
			_, _ = coerce.Float64(I16(20))
			_, _ = coerce.Float64(I32(30))
			_, _ = coerce.Float64(I64(40))
			_, _ = coerce.Float64(II(50))
			_, _ = coerce.Float64(uint8(10))
			_, _ = coerce.Float64(uint16(20))
			_, _ = coerce.Float64(uint32(30))
			_, _ = coerce.Float64(uint64(40))
			_, _ = coerce.Float64(uint(50))
			_, _ = coerce.Float64(U8(10))
			_, _ = coerce.Float64(U16(20))
			_, _ = coerce.Float64(U32(30))
			_, _ = coerce.Float64(U64(40))
			_, _ = coerce.Float64(UU(50))
			_, _ = coerce.Float64("42")
			_, _ = coerce.Float64(S("55"))
		}
	})
	b.Run("iint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Int(true)
			_, _ = coerce.Int(false)
			_, _ = coerce.Int(B(true))
			_, _ = coerce.Int(B(false))
			_, _ = coerce.Int(float32(33))
			_, _ = coerce.Int(float64(66))
			_, _ = coerce.Int(F32(33))
			_, _ = coerce.Int(F64(66))
			_, _ = coerce.Int(int8(10))
			_, _ = coerce.Int(int16(20))
			_, _ = coerce.Int(int32(30))
			_, _ = coerce.Int(int64(40))
			_, _ = coerce.Int(int(50))
			_, _ = coerce.Int(I8(10))
			_, _ = coerce.Int(I16(20))
			_, _ = coerce.Int(I32(30))
			_, _ = coerce.Int(I64(40))
			_, _ = coerce.Int(II(50))
			_, _ = coerce.Int(uint8(10))
			_, _ = coerce.Int(uint16(20))
			_, _ = coerce.Int(uint32(30))
			_, _ = coerce.Int(uint64(40))
			_, _ = coerce.Int(uint(50))
			_, _ = coerce.Int(U8(10))
			_, _ = coerce.Int(U16(20))
			_, _ = coerce.Int(U32(30))
			_, _ = coerce.Int(U64(40))
			_, _ = coerce.Int(UU(50))
			_, _ = coerce.Int("42")
			_, _ = coerce.Int(S("55"))
		}
	})
	b.Run("int8", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Int8(true)
			_, _ = coerce.Int8(false)
			_, _ = coerce.Int8(B(true))
			_, _ = coerce.Int8(B(false))
			_, _ = coerce.Int8(float32(33))
			_, _ = coerce.Int8(float64(66))
			_, _ = coerce.Int8(F32(33))
			_, _ = coerce.Int8(F64(66))
			_, _ = coerce.Int8(int8(10))
			_, _ = coerce.Int8(int16(20))
			_, _ = coerce.Int8(int32(30))
			_, _ = coerce.Int8(int64(40))
			_, _ = coerce.Int8(int(50))
			_, _ = coerce.Int8(I8(10))
			_, _ = coerce.Int8(I16(20))
			_, _ = coerce.Int8(I32(30))
			_, _ = coerce.Int8(I64(40))
			_, _ = coerce.Int8(II(50))
			_, _ = coerce.Int8(uint8(10))
			_, _ = coerce.Int8(uint16(20))
			_, _ = coerce.Int8(uint32(30))
			_, _ = coerce.Int8(uint64(40))
			_, _ = coerce.Int8(uint(50))
			_, _ = coerce.Int8(U8(10))
			_, _ = coerce.Int8(U16(20))
			_, _ = coerce.Int8(U32(30))
			_, _ = coerce.Int8(U64(40))
			_, _ = coerce.Int8(UU(50))
			_, _ = coerce.Int8("42")
			_, _ = coerce.Int8(S("55"))
		}
	})
	b.Run("int16", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Int16(true)
			_, _ = coerce.Int16(false)
			_, _ = coerce.Int16(B(true))
			_, _ = coerce.Int16(B(false))
			_, _ = coerce.Int16(float32(33))
			_, _ = coerce.Int16(float64(66))
			_, _ = coerce.Int16(F32(33))
			_, _ = coerce.Int16(F64(66))
			_, _ = coerce.Int16(int8(10))
			_, _ = coerce.Int16(int16(20))
			_, _ = coerce.Int16(int32(30))
			_, _ = coerce.Int16(int64(40))
			_, _ = coerce.Int16(int(50))
			_, _ = coerce.Int16(I8(10))
			_, _ = coerce.Int16(I16(20))
			_, _ = coerce.Int16(I32(30))
			_, _ = coerce.Int16(I64(40))
			_, _ = coerce.Int16(II(50))
			_, _ = coerce.Int16(uint8(10))
			_, _ = coerce.Int16(uint16(20))
			_, _ = coerce.Int16(uint32(30))
			_, _ = coerce.Int16(uint64(40))
			_, _ = coerce.Int16(uint(50))
			_, _ = coerce.Int16(U8(10))
			_, _ = coerce.Int16(U16(20))
			_, _ = coerce.Int16(U32(30))
			_, _ = coerce.Int16(U64(40))
			_, _ = coerce.Int16(UU(50))
			_, _ = coerce.Int16("42")
			_, _ = coerce.Int16(S("55"))
		}
	})
	b.Run("int32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Int32(true)
			_, _ = coerce.Int32(false)
			_, _ = coerce.Int32(B(true))
			_, _ = coerce.Int32(B(false))
			_, _ = coerce.Int32(float32(33))
			_, _ = coerce.Int32(float64(66))
			_, _ = coerce.Int32(F32(33))
			_, _ = coerce.Int32(F64(66))
			_, _ = coerce.Int32(int8(10))
			_, _ = coerce.Int32(int16(20))
			_, _ = coerce.Int32(int32(30))
			_, _ = coerce.Int32(int64(40))
			_, _ = coerce.Int32(int(50))
			_, _ = coerce.Int32(I8(10))
			_, _ = coerce.Int32(I16(20))
			_, _ = coerce.Int32(I32(30))
			_, _ = coerce.Int32(I64(40))
			_, _ = coerce.Int32(II(50))
			_, _ = coerce.Int32(uint8(10))
			_, _ = coerce.Int32(uint16(20))
			_, _ = coerce.Int32(uint32(30))
			_, _ = coerce.Int32(uint64(40))
			_, _ = coerce.Int32(uint(50))
			_, _ = coerce.Int32(U8(10))
			_, _ = coerce.Int32(U16(20))
			_, _ = coerce.Int32(U32(30))
			_, _ = coerce.Int32(U64(40))
			_, _ = coerce.Int32(UU(50))
			_, _ = coerce.Int32("42")
			_, _ = coerce.Int32(S("55"))
		}
	})
	b.Run("int64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Int64(true)
			_, _ = coerce.Int64(false)
			_, _ = coerce.Int64(B(true))
			_, _ = coerce.Int64(B(false))
			_, _ = coerce.Int64(float32(33))
			_, _ = coerce.Int64(float64(66))
			_, _ = coerce.Int64(F32(33))
			_, _ = coerce.Int64(F64(66))
			_, _ = coerce.Int64(int8(10))
			_, _ = coerce.Int64(int16(20))
			_, _ = coerce.Int64(int32(30))
			_, _ = coerce.Int64(int64(40))
			_, _ = coerce.Int64(int(50))
			_, _ = coerce.Int64(I8(10))
			_, _ = coerce.Int64(I16(20))
			_, _ = coerce.Int64(I32(30))
			_, _ = coerce.Int64(I64(40))
			_, _ = coerce.Int64(II(50))
			_, _ = coerce.Int64(uint8(10))
			_, _ = coerce.Int64(uint16(20))
			_, _ = coerce.Int64(uint32(30))
			_, _ = coerce.Int64(uint64(40))
			_, _ = coerce.Int64(uint(50))
			_, _ = coerce.Int64(U8(10))
			_, _ = coerce.Int64(U16(20))
			_, _ = coerce.Int64(U32(30))
			_, _ = coerce.Int64(U64(40))
			_, _ = coerce.Int64(UU(50))
			_, _ = coerce.Int64("42")
			_, _ = coerce.Int64(S("55"))
		}
	})
	b.Run("uuint", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Uint(true)
			_, _ = coerce.Uint(false)
			_, _ = coerce.Uint(B(true))
			_, _ = coerce.Uint(B(false))
			_, _ = coerce.Uint(float32(33))
			_, _ = coerce.Uint(float64(66))
			_, _ = coerce.Uint(F32(33))
			_, _ = coerce.Uint(F64(66))
			_, _ = coerce.Uint(int8(10))
			_, _ = coerce.Uint(int16(20))
			_, _ = coerce.Uint(int32(30))
			_, _ = coerce.Uint(int64(40))
			_, _ = coerce.Uint(int(50))
			_, _ = coerce.Uint(I8(10))
			_, _ = coerce.Uint(I16(20))
			_, _ = coerce.Uint(I32(30))
			_, _ = coerce.Uint(I64(40))
			_, _ = coerce.Uint(II(50))
			_, _ = coerce.Uint(uint8(10))
			_, _ = coerce.Uint(uint16(20))
			_, _ = coerce.Uint(uint32(30))
			_, _ = coerce.Uint(uint64(40))
			_, _ = coerce.Uint(uint(50))
			_, _ = coerce.Uint(U8(10))
			_, _ = coerce.Uint(U16(20))
			_, _ = coerce.Uint(U32(30))
			_, _ = coerce.Uint(U64(40))
			_, _ = coerce.Uint(UU(50))
			_, _ = coerce.Uint("42")
			_, _ = coerce.Uint(S("55"))
		}
	})
	b.Run("uint8", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Uint8(true)
			_, _ = coerce.Uint8(false)
			_, _ = coerce.Uint8(B(true))
			_, _ = coerce.Uint8(B(false))
			_, _ = coerce.Uint8(float32(33))
			_, _ = coerce.Uint8(float64(66))
			_, _ = coerce.Uint8(F32(33))
			_, _ = coerce.Uint8(F64(66))
			_, _ = coerce.Uint8(int8(10))
			_, _ = coerce.Uint8(int16(20))
			_, _ = coerce.Uint8(int32(30))
			_, _ = coerce.Uint8(int64(40))
			_, _ = coerce.Uint8(int(50))
			_, _ = coerce.Uint8(I8(10))
			_, _ = coerce.Uint8(I16(20))
			_, _ = coerce.Uint8(I32(30))
			_, _ = coerce.Uint8(I64(40))
			_, _ = coerce.Uint8(II(50))
			_, _ = coerce.Uint8(uint8(10))
			_, _ = coerce.Uint8(uint16(20))
			_, _ = coerce.Uint8(uint32(30))
			_, _ = coerce.Uint8(uint64(40))
			_, _ = coerce.Uint8(uint(50))
			_, _ = coerce.Uint8(U8(10))
			_, _ = coerce.Uint8(U16(20))
			_, _ = coerce.Uint8(U32(30))
			_, _ = coerce.Uint8(U64(40))
			_, _ = coerce.Uint8(UU(50))
			_, _ = coerce.Uint8("42")
			_, _ = coerce.Uint8(S("55"))
		}
	})
	b.Run("uint16", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Uint16(true)
			_, _ = coerce.Uint16(false)
			_, _ = coerce.Uint16(B(true))
			_, _ = coerce.Uint16(B(false))
			_, _ = coerce.Uint16(float32(33))
			_, _ = coerce.Uint16(float64(66))
			_, _ = coerce.Uint16(F32(33))
			_, _ = coerce.Uint16(F64(66))
			_, _ = coerce.Uint16(int8(10))
			_, _ = coerce.Uint16(int16(20))
			_, _ = coerce.Uint16(int32(30))
			_, _ = coerce.Uint16(int64(40))
			_, _ = coerce.Uint16(int(50))
			_, _ = coerce.Uint16(I8(10))
			_, _ = coerce.Uint16(I16(20))
			_, _ = coerce.Uint16(I32(30))
			_, _ = coerce.Uint16(I64(40))
			_, _ = coerce.Uint16(II(50))
			_, _ = coerce.Uint16(uint8(10))
			_, _ = coerce.Uint16(uint16(20))
			_, _ = coerce.Uint16(uint32(30))
			_, _ = coerce.Uint16(uint64(40))
			_, _ = coerce.Uint16(uint(50))
			_, _ = coerce.Uint16(U8(10))
			_, _ = coerce.Uint16(U16(20))
			_, _ = coerce.Uint16(U32(30))
			_, _ = coerce.Uint16(U64(40))
			_, _ = coerce.Uint16(UU(50))
			_, _ = coerce.Uint16("42")
			_, _ = coerce.Uint16(S("55"))
		}
	})
	b.Run("uint32", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Uint32(true)
			_, _ = coerce.Uint32(false)
			_, _ = coerce.Uint32(B(true))
			_, _ = coerce.Uint32(B(false))
			_, _ = coerce.Uint32(float32(33))
			_, _ = coerce.Uint32(float64(66))
			_, _ = coerce.Uint32(F32(33))
			_, _ = coerce.Uint32(F64(66))
			_, _ = coerce.Uint32(int8(10))
			_, _ = coerce.Uint32(int16(20))
			_, _ = coerce.Uint32(int32(30))
			_, _ = coerce.Uint32(int64(40))
			_, _ = coerce.Uint32(int(50))
			_, _ = coerce.Uint32(I8(10))
			_, _ = coerce.Uint32(I16(20))
			_, _ = coerce.Uint32(I32(30))
			_, _ = coerce.Uint32(I64(40))
			_, _ = coerce.Uint32(II(50))
			_, _ = coerce.Uint32(uint8(10))
			_, _ = coerce.Uint32(uint16(20))
			_, _ = coerce.Uint32(uint32(30))
			_, _ = coerce.Uint32(uint64(40))
			_, _ = coerce.Uint32(uint(50))
			_, _ = coerce.Uint32(U8(10))
			_, _ = coerce.Uint32(U16(20))
			_, _ = coerce.Uint32(U32(30))
			_, _ = coerce.Uint32(U64(40))
			_, _ = coerce.Uint32(UU(50))
			_, _ = coerce.Uint32("42")
			_, _ = coerce.Uint32(S("55"))
		}
	})
	b.Run("uint64", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.Uint64(true)
			_, _ = coerce.Uint64(false)
			_, _ = coerce.Uint64(B(true))
			_, _ = coerce.Uint64(B(false))
			_, _ = coerce.Uint64(float32(33))
			_, _ = coerce.Uint64(float64(66))
			_, _ = coerce.Uint64(F32(33))
			_, _ = coerce.Uint64(F64(66))
			_, _ = coerce.Uint64(int8(10))
			_, _ = coerce.Uint64(int16(20))
			_, _ = coerce.Uint64(int32(30))
			_, _ = coerce.Uint64(int64(40))
			_, _ = coerce.Uint64(int(50))
			_, _ = coerce.Uint64(I8(10))
			_, _ = coerce.Uint64(I16(20))
			_, _ = coerce.Uint64(I32(30))
			_, _ = coerce.Uint64(I64(40))
			_, _ = coerce.Uint64(II(50))
			_, _ = coerce.Uint64(uint8(10))
			_, _ = coerce.Uint64(uint16(20))
			_, _ = coerce.Uint64(uint32(30))
			_, _ = coerce.Uint64(uint64(40))
			_, _ = coerce.Uint64(uint(50))
			_, _ = coerce.Uint64(U8(10))
			_, _ = coerce.Uint64(U16(20))
			_, _ = coerce.Uint64(U32(30))
			_, _ = coerce.Uint64(U64(40))
			_, _ = coerce.Uint64(UU(50))
			_, _ = coerce.Uint64("42")
			_, _ = coerce.Uint64(S("55"))
		}
	})
	b.Run("string", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, _ = coerce.String(true)
			_, _ = coerce.String(false)
			_, _ = coerce.String(B(true))
			_, _ = coerce.String(B(false))
			_, _ = coerce.String(float32(33))
			_, _ = coerce.String(float64(66))
			_, _ = coerce.String(F32(33))
			_, _ = coerce.String(F64(66))
			_, _ = coerce.String(int8(10))
			_, _ = coerce.String(int16(20))
			_, _ = coerce.String(int32(30))
			_, _ = coerce.String(int64(40))
			_, _ = coerce.String(int(50))
			_, _ = coerce.String(I8(10))
			_, _ = coerce.String(I16(20))
			_, _ = coerce.String(I32(30))
			_, _ = coerce.String(I64(40))
			_, _ = coerce.String(II(50))
			_, _ = coerce.String(uint8(10))
			_, _ = coerce.String(uint16(20))
			_, _ = coerce.String(uint32(30))
			_, _ = coerce.String(uint64(40))
			_, _ = coerce.String(uint(50))
			_, _ = coerce.String(U8(10))
			_, _ = coerce.String(U16(20))
			_, _ = coerce.String(U32(30))
			_, _ = coerce.String(U64(40))
			_, _ = coerce.String(UU(50))
			_, _ = coerce.String("42")
			_, _ = coerce.String(S("55"))
		}
	})
}
