package coerce_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// I8 is a new type derived from int8.
type I8 int8

// I16 is a new type derived from int16.
type I16 int16

// I32 is a new type derived from int32.
type I32 int32

// I64 is a new type derived from int64.
type I64 int64

// II is a new type derived from int.
type II int

// IntTest is the struct used to build up table driven tests for ints.
type IntTest struct {
	To       interface{}
	Error8   error
	Expect8  int8
	Error16  error
	Expect16 int16
	Error32  error
	Expect32 int32
	Error64  error
	Expect64 int64
	Error    error
	Expect   int
}

// IntTests is a table of IntTest.
type IntTests map[string]*IntTest

// Run iterates the IntTests and runs each one.
func (tests IntTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("int8 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i8, err := coerce.Int8(test.To)
			chk.True(errors.Is(err, test.Error8))
			chk.Equal(test.Expect8, i8)
		})
		t.Run("int16 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i16, err := coerce.Int16(test.To)
			chk.True(errors.Is(err, test.Error16))
			chk.Equal(test.Expect16, i16)
		})
		t.Run("int32 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i32, err := coerce.Int32(test.To)
			chk.True(errors.Is(err, test.Error32))
			chk.Equal(test.Expect32, i32)
		})
		t.Run("int64 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i64, err := coerce.Int64(test.To)
			chk.True(errors.Is(err, test.Error64))
			chk.Equal(test.Expect64, i64)
		})
		t.Run("int "+name, func(t *testing.T) {
			chk := assert.New(t)
			i, err := coerce.Int(test.To)
			chk.True(errors.Is(err, test.Error))
			chk.Equal(test.Expect, i)
		})
	}
}

// RunDelta iterates the IntTests and runs each one but uses InDelta checks to handle
// cases where To is some type of float and may be imprecise.
func (tests IntTests) RunDelta(t *testing.T) {
	for name, test := range tests {
		t.Run("int8 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i8, err := coerce.Int8(test.To)
			chk.True(errors.Is(err, test.Error8))
			if test.Expect8 != int8(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect8, i8, 1.0)
			} else {
				chk.Equal(test.Expect8, i8)
			}
		})
		t.Run("int16 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i16, err := coerce.Int16(test.To)
			chk.True(errors.Is(err, test.Error16))
			if test.Expect16 != int16(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect16, i16, 1.0)
			} else {
				chk.Equal(test.Expect16, i16)
			}
		})
		t.Run("int32 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i32, err := coerce.Int32(test.To)
			chk.True(errors.Is(err, test.Error32))
			if test.Expect32 != int32(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect32, i32, 1.0)
			} else {
				chk.Equal(test.Expect32, i32)
			}
		})
		t.Run("int64 "+name, func(t *testing.T) {
			chk := assert.New(t)
			i64, err := coerce.Int64(test.To)
			chk.True(errors.Is(err, test.Error64))
			if test.Expect64 != int64(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect64, i64, 1.0)
			} else {
				chk.Equal(test.Expect64, i64)
			}
		})
		t.Run("int "+name, func(t *testing.T) {
			chk := assert.New(t)
			i, err := coerce.Int(test.To)
			chk.True(errors.Is(err, test.Error))
			if test.Expect != int(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect, i, 1.0)
			} else {
				chk.Equal(test.Expect, i)
			}
		})
	}
}

func TestIntFromBool(t *testing.T) {
	tests := IntTests{
		"false": {
			To:      false,
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"true": {
			To:      true,
			Expect8: 1, Expect16: 1, Expect32: 1, Expect64: 1, Expect: 1,
		},

		"B(false)": {
			To:      B(false),
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"B(true)": {
			To:      B(true),
			Expect8: 1, Expect16: 1, Expect32: 1, Expect64: 1, Expect: 1,
		},
	}
	tests.Run(t)
}

func TestIntFromKind(t *testing.T) {
	tests := IntTests{
		"I8 max8": {
			To:      I8(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"I16 max16": {
			To:     I16(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"I32 max32": {
			To:     I32(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"I64 max64": {
			To:     I64(math.MaxInt64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64,
		},
		"II max32": {
			To:     II(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
	}
	if strconv.IntSize == 64 {
		tests["I64 max64"].Expect = math.MaxInt64
	} else {
		tests["I64 max64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestIntFromFloat(t *testing.T) {
	tests := IntTests{
		"max8": {
			To:      float32(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"max16": {
			To:     float32(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"max32": {
			To:     float32(math.MaxInt32) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32 / 4, Expect64: math.MaxInt32 / 4, Expect: math.MaxInt32 / 4,
		},
		"max64": {
			To:     float32(math.MaxInt64) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64 / 4,
		},
		"float64 max": {
			To:     float64(math.MaxFloat64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"float32 min": {
			To:     float32(-1 * math.MaxFloat32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},

		"F32(max8)": {
			To:      F32(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"F32(max16)": {
			To:     F32(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"F32(max32)": {
			To:     F32(math.MaxInt32) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32 / 4, Expect64: math.MaxInt32 / 4, Expect: math.MaxInt32 / 4,
		},
		"F32(max64)": {
			To:     F32(math.MaxInt64) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64 / 4,
		},
		"F64(float64 max)": {
			To:     F64(math.MaxFloat64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"F32(float32 min)": {
			To:     F32(-1 * math.MaxFloat32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxInt64 / 4
		tests["F32(max64)"].Expect = math.MaxInt64 / 4
	} else {
		tests["max64"].Error = coerce.ErrOverflow
		tests["F32(max64)"].Error = coerce.ErrOverflow
	}
	tests.RunDelta(t)
}

func TestIntFromInt(t *testing.T) {
	tests := IntTests{
		"max8": {
			To:      int8(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"max16": {
			To:     int16(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"max32": {
			To:     int32(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"max64": {
			To:     math.MaxInt64,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64,
		},

		"low int8": {
			To:      int8(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low int16": {
			To:      int16(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low int32": {
			To:      int32(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low int64": {
			To:      int64(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low int": {
			To:      int(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxInt64
	} else {
		tests["max64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestIntUnsupported(t *testing.T) {
	tests := IntTests{
		"unsupported": {
			To:      map[string]string{},
			Error8:  coerce.ErrUnsupported,
			Error16: coerce.ErrUnsupported,
			Error32: coerce.ErrUnsupported,
			Error64: coerce.ErrUnsupported,
			Error:   coerce.ErrUnsupported,
		},
	}
	tests.Run(t)
}

func TestIntFromNil(t *testing.T) {
	tests := IntTests{
		"nil": {
			To: nil,
		},
	}
	tests.Run(t)
}

func TestIntFromPtr(t *testing.T) {
	n32 := (*int32)(nil)
	nn32 := &n32
	i32, i64 := int32(math.MaxInt32), int64(math.MaxInt64)
	p32, p64 := &i32, &i64
	pp32, pp64 := &p32, &p64
	ppp32, ppp64 := &pp32, &pp64
	//
	tests := IntTests{
		"nil": {
			To:      n32,
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"*nil": {
			To:      nn32,
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"ppp32": {
			To:     ppp32,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: *p32, Expect64: int64(*p32), Expect: int(*p32),
		},
		"ppp64": {
			To:     ppp64,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: *p64,
		},
	}
	if strconv.IntSize == 64 {
		tests["ppp64"].Expect = int(*p64)
	} else {
		tests["ppp64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestIntFromSlice(t *testing.T) {
	tests := IntTests{
		"nil": {
			To:      []string(nil),
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"slice": {
			To:      []interface{}{"42", "78"},
			Expect8: 78, Expect16: 78, Expect32: 78, Expect64: 78, Expect: 78,
		},
		"slice kind": {
			To:      []interface{}{"42", "78", F32(78)},
			Expect8: 78, Expect16: 78, Expect32: 78, Expect64: 78, Expect: 78,
		},
	}
	tests.Run(t)
}

func TestIntFromUint(t *testing.T) {
	tests := IntTests{
		"max8": {
			To:      uint8(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"max u8": {
			To:     uint8(math.MaxUint8),
			Error8: coerce.ErrOverflow, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"max16": {
			To:     uint16(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"max u16": {
			To:     uint16(math.MaxUint16),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"max32": {
			To:     uint32(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"max u32": {
			To:     uint32(math.MaxUint32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint32, Expect: math.MaxUint32,
		},
		"max64": {
			To:     uint64(math.MaxInt64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64,
		},
		"uint64 max": {
			To:     uint64(math.MaxUint64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"max32+1": {
			To:     uint(math.MaxInt32 + 1),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt32 + 1, Expect: math.MaxInt32 + 1,
		},

		"low uint8": {
			To:      uint8(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low uint16": {
			To:      uint16(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low uint32": {
			To:      uint32(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low uint64": {
			To:      uint64(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},
		"low uint": {
			To:      uint(4),
			Expect8: 4, Expect16: 4, Expect32: 4, Expect64: 4, Expect: 4,
		},

		"U8(max8)": {
			To:      U8(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"U16(max16)": {
			To:     U16(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"U32(max32)": {
			To:     U32(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"U64(max64)": {
			To:     U64(math.MaxInt64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64,
		},
		"U64(uint64 max)": {
			To:     U64(math.MaxUint64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"UU(max32+1)": {
			To:     UU(math.MaxInt32 + 1),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt32 + 1, Expect: math.MaxInt32 + 1,
		},
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxInt64
		tests["U64(max64)"].Expect = math.MaxInt64
	} else {
		tests["max64"].Error = coerce.ErrOverflow
		tests["U64(max64)"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestIntFromString(t *testing.T) {
	String := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	tests := IntTests{
		"max8": {
			To:      String(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"max16": {
			To:     String(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"max32": {
			To:     String(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"max64": {
			To:     String(int64(math.MaxInt64)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64,
		},
		"3.14": {
			To:      "3.14",
			Expect8: 3, Expect16: 3, Expect32: 3, Expect64: 3, Expect: 3,
		},
		"max float32": {
			To:     String(float32(math.MaxFloat32)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"uint": {
			To:     String(uint64(math.MaxInt32 + 1)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt32 + 1,
		},
		"uint max": {
			To:     String(uint64(math.MaxUint64)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"hello": {
			To:     "Hello",
			Error8: coerce.ErrInvalid, Error16: coerce.ErrInvalid, Error32: coerce.ErrInvalid, Error64: coerce.ErrInvalid, Error: coerce.ErrInvalid,
		},

		"S(3.14)": {
			To:      S("3.14"),
			Expect8: 3, Expect16: 3, Expect32: 3, Expect64: 3, Expect: 3,
		},
	}
	//
	// Expectations for some tests depend on int bit-size
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxInt64
		tests["uint"].Expect = math.MaxInt32 + 1
	} else {
		tests["max64"].Error = coerce.ErrOverflow
		tests["uint"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}
