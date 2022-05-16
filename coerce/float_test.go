package coerce_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// F32 is a type derived from float32.
type F32 float32

// F64 is a type derived from float64.
type F64 float64

// FloatTest is the struct used to build up table driven tests for floats.
type FloatTest struct {
	To       interface{}
	Error32  error
	Expect32 float32
	Error64  error
	Expect64 float64
}

// FloatTests is a table of FloatTest.
type FloatTests map[string]*FloatTest

// Run iterates the FloatTests and runs each one.
func (tests FloatTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("float32 "+name, func(t *testing.T) {
			chk := assert.New(t)
			f32, err := coerce.Float32(test.To)
			chk.True(errors.Is(err, test.Error32))
			if test.Expect32 != float32(0) {
				chk.InEpsilon(test.Expect32, f32, 0.1)
			} else {
				chk.Equal(test.Expect32, f32)
			}
		})
		t.Run("float64 "+name, func(t *testing.T) {
			chk := assert.New(t)
			f64, err := coerce.Float64(test.To)
			chk.True(errors.Is(err, test.Error64))
			if test.Expect64 != float64(0) {
				chk.InEpsilon(test.Expect64, f64, 0.1)
			} else {
				chk.Equal(test.Expect64, f64)
			}
		})
	}
}

func TestFloatFromBool(t *testing.T) {
	tests := FloatTests{
		"false": {To: false, Expect32: 0, Expect64: 0},
		"true":  {To: true, Expect32: 1, Expect64: 1},

		"B(false)": {To: B(false), Expect32: 0, Expect64: 0},
		"B(true)":  {To: B(true), Expect32: 1, Expect64: 1},
	}
	tests.Run(t)
}

func TestFloatFromKind(t *testing.T) {
	tests := FloatTests{
		"max32": {
			To:       F32(math.MaxFloat32),
			Expect32: math.MaxFloat32, Expect64: math.MaxFloat32,
		},
		"max64": {
			To:      F64(math.MaxFloat64),
			Error32: coerce.ErrOverflow, Expect64: math.MaxFloat64,
		},
		"F64 max32": {
			To:       F64(math.MaxFloat32),
			Expect32: math.MaxFloat32, Expect64: math.MaxFloat32,
		},
	}
	tests.Run(t)
}

func TestFloatFromFloat(t *testing.T) {
	tests := FloatTests{
		"max32": {
			To:       float32(math.MaxFloat32),
			Expect32: math.MaxFloat32, Expect64: math.MaxFloat32,
		},
		"max64": {
			To:      float64(math.MaxFloat64),
			Error32: coerce.ErrOverflow, Expect64: math.MaxFloat64,
		},

		"F32(max32)": {
			To:       F32(math.MaxFloat32),
			Expect32: math.MaxFloat32, Expect64: math.MaxFloat32,
		},
		"F64(max64)": {
			To:      F64(math.MaxFloat64),
			Error32: coerce.ErrOverflow, Expect64: math.MaxFloat64,
		},
	}
	tests.Run(t)
}

func TestFloatFromInt(t *testing.T) {
	tests := FloatTests{
		"max8": {
			To:       int8(math.MaxInt8),
			Expect32: math.MaxInt8, Expect64: math.MaxInt8,
		},
		"max16": {
			To:       int16(math.MaxInt16),
			Expect32: math.MaxInt16, Expect64: math.MaxInt16,
		},
		"max32": {
			To:       int32(math.MaxInt32),
			Expect32: math.MaxInt32, Expect64: math.MaxInt32,
		},
		"max64": {
			To:       int64(math.MaxInt64),
			Expect32: math.MaxInt64, Expect64: math.MaxInt64,
		},
		"int": {
			To:       int(math.MaxInt32),
			Expect32: math.MaxInt32, Expect64: math.MaxInt32,
		},

		"I8(max8)": {
			To:       I8(math.MaxInt8),
			Expect32: math.MaxInt8, Expect64: math.MaxInt8,
		},
		"I16(max16)": {
			To:       I16(math.MaxInt16),
			Expect32: math.MaxInt16, Expect64: math.MaxInt16,
		},
		"I32(max32)": {
			To:       I32(math.MaxInt32),
			Expect32: math.MaxInt32, Expect64: math.MaxInt32,
		},
		"I64(max64)": {
			To:       I64(math.MaxInt64),
			Expect32: math.MaxInt64, Expect64: math.MaxInt64,
		},
		"II(int)": {
			To:       II(math.MaxInt32),
			Expect32: math.MaxInt32, Expect64: math.MaxInt32,
		},
	}
	tests.Run(t)
}

func TestFloatUnsupported(t *testing.T) {
	tests := FloatTests{
		"unsupported": {
			To: map[string]string{}, Error32: coerce.ErrUnsupported, Error64: coerce.ErrUnsupported,
		},
	}
	tests.Run(t)
}

func TestFloatFromNil(t *testing.T) {
	tests := FloatTests{
		"nil": {
			To: nil, Expect32: 0, Expect64: 0,
		},
	}
	tests.Run(t)
}

func TestFloatFromPtr(t *testing.T) {
	n32 := (*float32)(nil)
	nn32 := &n32
	f32, f64 := float32(math.MaxFloat32), float64(math.MaxFloat64)
	p32, p64 := &f32, &f64
	pp32, pp64 := &p32, &p64
	ppp32, ppp64 := &pp32, &pp64
	//
	tests := FloatTests{
		"nil": {
			To: n32, Expect32: 0, Expect64: 0,
		},
		"*nil": {
			To: nn32, Expect32: 0, Expect64: 0,
		},
		"ppp32": {
			To: ppp32, Expect32: *p32, Expect64: float64(*p32),
		},
		"ppp64": {
			To: ppp64, Error32: coerce.ErrOverflow, Expect64: *p64,
		},
	}
	tests.Run(t)
}

func TestFloatFromSlice(t *testing.T) {
	tests := FloatTests{
		"nil":        {To: []string(nil), Expect32: 0, Expect64: 0},
		"slice":      {To: []interface{}{"42", "78"}, Expect32: 78, Expect64: 78},
		"slice kind": {To: []interface{}{"42", "78", S("27")}, Expect32: 27, Expect64: 27},
	}
	tests.Run(t)
}

func TestFloatFromUint(t *testing.T) {
	tests := FloatTests{
		"uint8": {
			To: uint8(math.MaxUint8), Expect32: math.MaxUint8, Expect64: math.MaxUint8,
		},
		"uint16": {
			To: uint16(math.MaxUint16), Expect32: math.MaxUint16, Expect64: math.MaxUint16,
		},
		"uint32": {
			To: uint32(math.MaxUint32), Expect32: math.MaxUint32, Expect64: math.MaxUint32,
		},
		"uint64": {
			To: uint64(math.MaxUint64), Expect32: math.MaxUint64, Expect64: math.MaxUint64,
		},
		"uint": {
			To: uint(math.MaxUint32), Expect32: math.MaxUint32, Expect64: math.MaxUint32,
		},

		"U8(uint8)": {
			To: U8(math.MaxUint8), Expect32: math.MaxUint8, Expect64: math.MaxUint8,
		},
		"U16(uint16)": {
			To: U16(math.MaxUint16), Expect32: math.MaxUint16, Expect64: math.MaxUint16,
		},
		"U32(uint32)": {
			To: U32(math.MaxUint32), Expect32: math.MaxUint32, Expect64: math.MaxUint32,
		},
		"U64(uint64)": {
			To: U64(math.MaxUint64), Expect32: math.MaxUint64, Expect64: math.MaxUint64,
		},
		"UU(uint)": {
			To: UU(math.MaxUint32), Expect32: math.MaxUint32, Expect64: math.MaxUint32,
		},
	}
	tests.Run(t)
}

func TestFloatFromString(t *testing.T) {
	String := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	tests := FloatTests{
		"max8": {
			To:       String(math.MaxInt8),
			Expect32: math.MaxInt8, Expect64: math.MaxInt8,
		},
		"max16": {
			To:       String(math.MaxInt16),
			Expect32: math.MaxInt16, Expect64: math.MaxInt16,
		},
		"max32": {
			To:       String(math.MaxInt32),
			Expect32: math.MaxInt32, Expect64: math.MaxInt32,
		},
		"max64": {
			To:       String(int64(math.MaxInt64)),
			Expect32: math.MaxInt64, Expect64: math.MaxInt64,
		},
		"3.14": {
			To:       "3.14",
			Expect32: 3.14, Expect64: 3.14,
		},
		"max float32": {
			To:       String(float32(math.MaxFloat32)),
			Expect32: math.MaxFloat32, Expect64: float64(math.MaxFloat32),
		},
		"uint": {
			To:       String(uint64(math.MaxInt32 + 1)),
			Expect32: math.MaxInt32 + 1, Expect64: math.MaxInt32 + 1,
		},
		"uint max": {
			To:       String(uint64(math.MaxUint64)),
			Expect32: math.MaxUint64, Expect64: math.MaxUint64,
		},
		"max float64": {
			To:      String(float64(math.MaxFloat64)),
			Error32: coerce.ErrOverflow, Expect64: math.MaxFloat64,
		},
		"overflow float64": {
			To:      "1.797693134862315708145274237317043567981e+1000",
			Error32: coerce.ErrOverflow, Error64: coerce.ErrOverflow,
		},
		"hello": {
			To:      "Hello",
			Error32: coerce.ErrInvalid, Error64: coerce.ErrInvalid,
		},

		"S(3.14)": {
			To:       S("3.14"),
			Expect32: 3.14, Expect64: 3.14,
		},
		"S(max float32)": {
			To:       S(String(float32(math.MaxFloat32))),
			Expect32: math.MaxFloat32, Expect64: float64(math.MaxFloat32),
		},

		`"false"`: {
			To:       "false",
			Expect32: 0, Expect64: 0,
		},
		`"true"`: {
			To:       "true",
			Expect32: 1, Expect64: 1,
		},
		`S("false")`: {
			To:       S("false"),
			Expect32: 0, Expect64: 0,
		},
		`S("true")`: {
			To:       S("true"),
			Expect32: 1, Expect64: 1,
		},
	}
	tests.Run(t)
}
