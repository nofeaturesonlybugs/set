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

// U8 is a new type derived from uint8.
type U8 uint8

// U16 is a new type derived from uint16.
type U16 uint16

// U32 is a new type derived from uint32.
type U32 uint32

// U64 is a new type derived from uint64.
type U64 uint64

// UU is a new type derived from uint.
type UU uint

// UintTest is the struct used to build up table driven tests for uints.
type UintTest struct {
	To       interface{}
	Error8   error
	Expect8  uint8
	Error16  error
	Expect16 uint16
	Error32  error
	Expect32 uint32
	Error64  error
	Expect64 uint64
	Error    error
	Expect   uint
}

// UintTests is a table of UintTest.
type UintTests map[string]*UintTest

// Run iterates the UintTests and runs each one.
func (tests UintTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("uint8 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u8, err := coerce.Uint8(test.To)
			chk.True(errors.Is(err, test.Error8))
			chk.Equal(test.Expect8, u8)
		})
		t.Run("uint16 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u16, err := coerce.Uint16(test.To)
			chk.True(errors.Is(err, test.Error16))
			chk.Equal(test.Expect16, u16)
		})
		t.Run("uint32 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u32, err := coerce.Uint32(test.To)
			chk.True(errors.Is(err, test.Error32))
			chk.Equal(test.Expect32, u32)
		})
		t.Run("uint64 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u64, err := coerce.Uint64(test.To)
			chk.True(errors.Is(err, test.Error64))
			chk.Equal(test.Expect64, u64)
		})
		t.Run("uint "+name, func(t *testing.T) {
			chk := assert.New(t)
			u, err := coerce.Uint(test.To)
			chk.True(errors.Is(err, test.Error))
			chk.Equal(test.Expect, u)
		})
	}
}

// RunDelta iterates the UintTests and runs each one but uses InDelta checks to handle
// cases where To is some type of float and may be imprecise.
func (tests UintTests) RunDelta(t *testing.T) {
	for name, test := range tests {
		t.Run("uint8 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u8, err := coerce.Uint8(test.To)
			chk.True(errors.Is(err, test.Error8))
			if test.Expect8 != uint8(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect8, u8, 1.0)
			} else {
				chk.Equal(test.Expect8, u8)
			}
		})
		t.Run("uint16 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u16, err := coerce.Uint16(test.To)
			chk.True(errors.Is(err, test.Error16))
			if test.Expect16 != uint16(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect16, u16, 1.0)
			} else {
				chk.Equal(test.Expect16, u16)
			}
		})
		t.Run("uint32 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u32, err := coerce.Uint32(test.To)
			chk.True(errors.Is(err, test.Error32))
			if test.Expect32 != uint32(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect32, u32, 1.0)
			} else {
				chk.Equal(test.Expect32, u32)
			}
		})
		t.Run("uint64 "+name, func(t *testing.T) {
			chk := assert.New(t)
			u64, err := coerce.Uint64(test.To)
			chk.True(errors.Is(err, test.Error64))
			if test.Expect64 != uint64(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect64, u64, 1.0)
			} else {
				chk.Equal(test.Expect64, u64)
			}
		})
		t.Run("uint "+name, func(t *testing.T) {
			chk := assert.New(t)
			u, err := coerce.Uint(test.To)
			chk.True(errors.Is(err, test.Error))
			if test.Expect != uint(0) { // type-case 0 so compiler warns if testing against wrong Expect*
				chk.InDelta(test.Expect, u, 1.0)
			} else {
				chk.Equal(test.Expect, u)
			}
		})
	}
}

func TestUintFromBool(t *testing.T) {
	tests := UintTests{
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

func TestUintFromKind(t *testing.T) {
	tests := UintTests{
		"U8 max8": {
			To: U8(math.MaxUint8), Expect8: math.MaxUint8, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"U16 max16": {
			To: U16(math.MaxUint16), Error8: coerce.ErrOverflow, Expect16: math.MaxUint16, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"U32 max32": {
			To: U32(math.MaxUint32), Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32, Expect64: math.MaxUint32, Expect: math.MaxUint32,
		},
		"U64 max64": {
			To: U64(math.MaxUint64), Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64,
		},
		"UU max32": {
			To: UU(math.MaxUint32), Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32, Expect64: math.MaxUint32, Expect: math.MaxUint32,
		},
	}
	if strconv.IntSize == 64 {
		tests["U64 max64"].Expect = math.MaxUint64
	} else {
		tests["U64 max64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestUintFromFloat(t *testing.T) {
	tests := UintTests{
		"max8": {
			To:      float32(math.MaxUint8),
			Expect8: math.MaxUint8, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"max16": {
			To:     float32(math.MaxUint16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxUint16, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"max32": {
			To:     float32(math.MaxUint32) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32 / 4, Expect64: math.MaxUint32 / 4, Expect: math.MaxUint32 / 4,
		},
		"max64": {
			To:     float32(math.MaxUint64) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64 / 4,
		},
		"float64 max": {
			To:     float64(math.MaxFloat64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"float32 min int32": {
			To:     float32(math.MinInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},

		"F32(max8)": {
			To:      F32(math.MaxUint8),
			Expect8: math.MaxUint8, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"F32(max16)": {
			To:     F32(math.MaxUint16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxUint16, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"F32(max32)": {
			To:     F32(math.MaxUint32) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32 / 4, Expect64: math.MaxUint32 / 4, Expect: math.MaxUint32 / 4,
		},
		"F32(max64)": {
			To:     F32(math.MaxUint64) / 4,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64 / 4,
		},
		"F64(float64 max)": {
			To:     F64(math.MaxFloat64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"F32(float32 min int32)": {
			To:     F32(math.MinInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxUint64 / 4
		tests["F32(max64)"].Expect = math.MaxUint64 / 4
	} else {
		tests["max64"].Error = coerce.ErrOverflow
		tests["F32(max64)"].Error = coerce.ErrOverflow
	}
	tests.RunDelta(t)
}

func TestUintFromInt(t *testing.T) {
	tests := UintTests{
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
			To:     int64(math.MaxInt64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64, Expect: math.MaxInt64,
		},
		"int": {
			To:     int(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
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

		"neg int8": {
			To:     int8(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"neg int16": {
			To:     int16(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"neg int32": {
			To:     int32(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"neg int64": {
			To:     int64(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
		"neg int": {
			To:     int(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},

		"I8(max8)": {
			To:      I8(math.MaxInt8),
			Expect8: math.MaxInt8, Expect16: math.MaxInt8, Expect32: math.MaxInt8, Expect64: math.MaxInt8, Expect: math.MaxInt8,
		},
		"I16(max16)": {
			To:     I16(math.MaxInt16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxInt16, Expect32: math.MaxInt16, Expect64: math.MaxInt16, Expect: math.MaxInt16,
		},
		"I32(max32)": {
			To:     I32(math.MaxInt32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxInt32, Expect64: math.MaxInt32, Expect: math.MaxInt32,
		},
		"I64(max64)": {
			To:     I64(math.MaxInt64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxInt64, Expect: math.MaxInt64,
		},
		"I8(min8)": {
			To:     I8(math.MinInt8),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow, Error: coerce.ErrOverflow,
		},
	}
	tests.Run(t)
}

func TestUintUnsupported(t *testing.T) {
	tests := UintTests{
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

func TestUintFromNil(t *testing.T) {
	tests := UintTests{
		"nil": {
			To:       nil,
			Expect8:  0,
			Expect16: 0,
			Expect32: 0,
			Expect64: 0,
			Expect:   0,
		},
	}
	tests.Run(t)
}

func TestUintFromPtr(t *testing.T) {
	n32 := (*uint32)(nil)
	nn32 := &n32
	i32, i64 := uint32(math.MaxUint32), uint64(math.MaxUint64)
	p32, p64 := &i32, &i64
	pp32, pp64 := &p32, &p64
	ppp32, ppp64 := &pp32, &pp64
	//
	tests := UintTests{
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
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: *p32, Expect64: uint64(*p32), Expect: uint(*p32),
		},
		"ppp64": {
			To:     ppp64,
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: *p64,
		},
	}
	if strconv.IntSize == 64 {
		tests["ppp64"].Expect = uint(*p64)
	} else {
		tests["ppp64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestUintFromSlice(t *testing.T) {
	tests := UintTests{
		"nil": {
			To:      []string(nil),
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		"slice": {
			To:      []interface{}{10, "42", "78"},
			Expect8: 78, Expect16: 78, Expect32: 78, Expect64: 78, Expect: 78,
		},
		"slice kind": {
			To:      []interface{}{10, "42", "78", S("78")},
			Expect8: 78, Expect16: 78, Expect32: 78, Expect64: 78, Expect: 78,
		},
	}
	tests.Run(t)
}

func TestUintFromUint(t *testing.T) {
	tests := UintTests{
		"max8": {
			To:      uint8(math.MaxUint8),
			Expect8: math.MaxUint8, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"max16": {
			To:     uint16(math.MaxUint16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxUint16, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"max32": {
			To:     uint32(math.MaxUint32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32, Expect64: math.MaxUint32, Expect: math.MaxUint32,
		},
		"max64": {
			To:     uint64(math.MaxUint64),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64,
		},
		"uint": {
			To:     uint(math.MaxUint32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32, Expect64: math.MaxUint32, Expect: math.MaxUint32,
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
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxUint64
	} else {
		tests["max64"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}

func TestUintFromString(t *testing.T) {
	String := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	tests := UintTests{
		"max8": {
			To:      String(math.MaxUint8),
			Expect8: math.MaxUint8, Expect16: math.MaxUint8, Expect32: math.MaxUint8, Expect64: math.MaxUint8, Expect: math.MaxUint8,
		},
		"max16": {
			To:     String(math.MaxUint16),
			Error8: coerce.ErrOverflow, Expect16: math.MaxUint16, Expect32: math.MaxUint16, Expect64: math.MaxUint16, Expect: math.MaxUint16,
		},
		"max32": {
			To:     String(math.MaxUint32),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Expect32: math.MaxUint32, Expect64: math.MaxUint32, Expect: math.MaxUint32,
		},
		"max64": {
			To:     String(uint64(math.MaxUint64)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64,
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
			To:     String(uint64(math.MaxUint32 + 1)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint32 + 1,
		},
		"uint max": {
			To:     String(uint64(math.MaxUint64)),
			Error8: coerce.ErrOverflow, Error16: coerce.ErrOverflow, Error32: coerce.ErrOverflow, Expect64: math.MaxUint64,
		},
		"negative int": {
			To:     String(-1),
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

		`"false"`: {
			To:      "false",
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		`"true"`: {
			To:      "true",
			Expect8: 1, Expect16: 1, Expect32: 1, Expect64: 1, Expect: 1,
		},
		`S("false")`: {
			To:      S("false"),
			Expect8: 0, Expect16: 0, Expect32: 0, Expect64: 0, Expect: 0,
		},
		`S("true")`: {
			To:      S("true"),
			Expect8: 1, Expect16: 1, Expect32: 1, Expect64: 1, Expect: 1,
		},
	}
	if strconv.IntSize == 64 {
		tests["max64"].Expect = math.MaxUint64
		tests["uint"].Expect = math.MaxUint32 + 1
		tests["uint max"].Expect = math.MaxUint64
	} else {
		tests["max64"].Error = coerce.ErrOverflow
		tests["uint"].Error = coerce.ErrOverflow
		tests["uint max"].Error = coerce.ErrOverflow
	}
	tests.Run(t)
}
