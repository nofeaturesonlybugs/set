package coerce_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// S is a new type derived from string.
type S string

// StringTest is the struct used to build up table driven tests for strings.
type StringTest struct {
	To     interface{}
	Error  error
	Expect string
}

// StringTests is a table of StringTests.
type StringTests map[string]*StringTest

// Run iterates the BoolTests and runs each one.
func (tests StringTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("string "+name, func(t *testing.T) {
			chk := assert.New(t)
			s, err := coerce.String(test.To)
			chk.True(errors.Is(err, test.Error))
			chk.Equal(test.Expect, s)
		})
	}
}

func TestStringFromBool(t *testing.T) {
	tests := StringTests{
		"false": {
			To: false, Expect: "false",
		},
		"true": {
			To: true, Expect: "true",
		},
		"B(false)": {
			To: B(false), Expect: "false",
		},
		"B(true)": {
			To: B(true), Expect: "true",
		},
	}
	tests.Run(t)
}

func TestStringFromKind(t *testing.T) {
	tests := StringTests{
		"empty": {
			To: S(""), Expect: "",
		},
		"hello": {
			To: S("Hello!"), Expect: "Hello!",
		},
	}
	tests.Run(t)
}

func TestStringFromFloat(t *testing.T) {
	tests := StringTests{
		"max8": {
			To: float32(math.MaxInt8), Expect: fmt.Sprintf("%v", float32(math.MaxInt8)),
		},
		"max16": {
			To: float32(math.MaxInt16), Expect: fmt.Sprintf("%v", float32(math.MaxInt16)),
		},
		"max32": {
			To: float32(math.MaxInt32) / 4, Expect: fmt.Sprintf("%v", float32(math.MaxInt32)/4),
		},
		"max64": {
			To: float32(math.MaxInt64) / 4, Expect: fmt.Sprintf("%v", float32(math.MaxInt64)/4),
		},
		"float64 max": {
			To: float64(math.MaxFloat64), Expect: fmt.Sprintf("%v", float64(math.MaxFloat64)),
		},
		"float32 min": {
			To: float32(-1 * math.MaxFloat32), Expect: fmt.Sprintf("%v", float32(-1*math.MaxFloat32)),
		},

		"F32(max16)": {
			To: F32(math.MaxInt16), Expect: fmt.Sprintf("%v", F32(math.MaxInt16)),
		},
		"F64(float64 max)": {
			To: F64(math.MaxFloat64), Expect: fmt.Sprintf("%v", float64(math.MaxFloat64)),
		},
	}
	tests.Run(t)
}

func TestStringFromInt(t *testing.T) {
	tests := StringTests{
		"max8": {
			To: int8(math.MaxInt8), Expect: fmt.Sprintf("%v", int8(math.MaxInt8)),
		},
		"max16": {
			To: int16(math.MaxInt16), Expect: fmt.Sprintf("%v", int16(math.MaxInt16)),
		},
		"max32": {
			To: int32(math.MaxInt32), Expect: fmt.Sprintf("%v", int32(math.MaxInt32)),
		},
		"max64": {
			To: int64(math.MaxInt64), Expect: fmt.Sprintf("%v", int64(math.MaxInt64)),
		},

		"I8(max8)": {
			To: I8(math.MaxInt8), Expect: fmt.Sprintf("%v", int8(math.MaxInt8)),
		},
		"I16(max16)": {
			To: I16(math.MaxInt16), Expect: fmt.Sprintf("%v", int16(math.MaxInt16)),
		},
		"I32(max32)": {
			To: I32(math.MaxInt32), Expect: fmt.Sprintf("%v", int32(math.MaxInt32)),
		},
		"I64(max64)": {
			To: I64(math.MaxInt64), Expect: fmt.Sprintf("%v", int64(math.MaxInt64)),
		},
	}
	tests.Run(t)
}

func TestStringUnsupported(t *testing.T) {
	tests := StringTests{
		"unsupported": {
			To:    map[string]string{},
			Error: coerce.ErrUnsupported,
		},
	}
	tests.Run(t)
}

func TestStringFromNil(t *testing.T) {
	chk := assert.New(t)
	var dst string
	var err error
	dst, err = coerce.String(nil)
	chk.NoError(err)
	chk.Equal("", dst)
}

func TestStringFromPtr(t *testing.T) {
	n := (*string)(nil)
	nn := &n
	s := "Hello!"
	ps := &s
	pps := &ps
	ppps := &pps
	//
	tests := StringTests{
		"nil": {
			To: n, Expect: "",
		},
		"*nil": {
			To: nn, Expect: "",
		},
		"ppps": {
			To: ppps, Expect: "Hello!",
		},
	}
	tests.Run(t)
}

func TestStringFromSlice(t *testing.T) {
	tests := StringTests{
		"nil": {
			To: []string(nil), Expect: "",
		},
		"slice": {
			To: []interface{}{"42", "78"}, Expect: "78",
		},
		"slice kind": {
			To: []interface{}{"42", "78", II(78)}, Expect: "78",
		},
	}
	tests.Run(t)
}

func TestStringFromUint(t *testing.T) {
	tests := StringTests{
		"max8": {
			To: uint8(math.MaxInt8), Expect: fmt.Sprintf("%v", uint8(math.MaxInt8)),
		},
		"max16": {
			To: uint16(math.MaxInt16), Expect: fmt.Sprintf("%v", uint16(math.MaxInt16)),
		},
		"max32": {
			To: uint32(math.MaxInt32), Expect: fmt.Sprintf("%v", uint32(math.MaxInt32)),
		},
		"max64": {
			To: math.MaxInt64, Expect: fmt.Sprintf("%v", math.MaxInt64),
		},
		"uint64 max": {
			To: uint64(math.MaxUint64), Expect: fmt.Sprintf("%v", uint64(math.MaxUint64)),
		},
		"max32+1": {
			To: uint(math.MaxInt32 + 1), Expect: fmt.Sprintf("%v", uint(math.MaxInt32+1)),
		},

		"U8(max8)": {
			To: U8(math.MaxInt8), Expect: fmt.Sprintf("%v", uint8(math.MaxInt8)),
		},
		"U16(max16)": {
			To: U16(math.MaxInt16), Expect: fmt.Sprintf("%v", uint16(math.MaxInt16)),
		},
		"U32(max32)": {
			To: U32(math.MaxInt32), Expect: fmt.Sprintf("%v", uint32(math.MaxInt32)),
		},
		"U64(uint64 max)": {
			To: U64(math.MaxUint64), Expect: fmt.Sprintf("%v", uint64(math.MaxUint64)),
		},
		"UU(max32+1)": {
			To: UU(math.MaxInt32 + 1), Expect: fmt.Sprintf("%v", uint(math.MaxInt32+1)),
		},
	}
	tests.Run(t)
}

func TestStringFromString(t *testing.T) {
	String := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	tests := StringTests{
		"3.14": {
			To: "3.14", Expect: "3.14",
		},
		"uint max": {
			To: String(uint64(math.MaxUint64)), Expect: fmt.Sprintf("%v", uint64(math.MaxUint64)),
		},
		"hello": {
			To: "Hello", Expect: "Hello",
		},
	}
	tests.Run(t)
}
