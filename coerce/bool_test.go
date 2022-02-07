package coerce_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set/coerce"
)

// B is a new type derived from bool.
type B bool

// BoolTest is the struct used to build up table driven tests for bools.
type BoolTest struct {
	To     interface{}
	Error  error
	Expect bool
}

// BoolTests is a table of BoolTest.
type BoolTests map[string]*BoolTest

// Run iterates the BoolTests and runs each one.
func (tests BoolTests) Run(t *testing.T) {
	for name, test := range tests {
		t.Run("bool "+name, func(t *testing.T) {
			chk := assert.New(t)
			got, err := coerce.Bool(test.To)
			chk.True(errors.Is(err, test.Error))
			chk.Equal(test.Expect, got)
		})
	}
}

func TestBoolFromBool(t *testing.T) {
	tests := BoolTests{
		"false": {
			To: false, Expect: false,
		},
		"true": {
			To: true, Expect: true,
		},

		"B(false)": {
			To: B(false), Expect: false,
		},
		"B(true)": {
			To: B(true), Expect: true,
		},
	}
	tests.Run(t)
}

func TestBoolFromFloat(t *testing.T) {
	tests := BoolTests{
		"max8": {
			To: float32(math.MaxInt8), Expect: true,
		},
		"max16": {
			To: float32(math.MaxInt16), Expect: true,
		},
		"float32(0)": {
			To: float32(0), Expect: false,
		},
		"float64(100)": {
			To: float64(100), Expect: true,
		},

		"F32(max8)": {
			To: F32(math.MaxInt8), Expect: true,
		},
		"F32(max16)": {
			To: F32(math.MaxInt16), Expect: true,
		},
		"F32(0)": {
			To: F32(0), Expect: false,
		},
		"F64(100)": {
			To: F64(100), Expect: true,
		},
	}
	tests.Run(t)
}

func TestBoolFromInt(t *testing.T) {
	tests := BoolTests{
		"max8": {
			To: int8(math.MaxInt8), Expect: true,
		},
		"max16": {
			To: int16(math.MaxInt16), Expect: true,
		},
		"min32": {
			To: int32(math.MinInt32), Expect: true,
		},
		"max64": {
			To: int64(0), Expect: false,
		},
		"0": {
			To: int(0), Expect: false,
		},

		"I8(max8)": {
			To: I8(math.MaxInt8), Expect: true,
		},
		"I16(max16)": {
			To: I16(math.MaxInt16), Expect: true,
		},
		"I32(min32)": {
			To: I32(math.MinInt32), Expect: true,
		},
		"I64(max64)": {
			To: I64(0), Expect: false,
		},
		"II(0)": {
			To: II(0), Expect: false,
		},
	}
	tests.Run(t)
}

func TestBoolUnsupported(t *testing.T) {
	chk := assert.New(t)
	var dst bool
	var err error
	{
		dst, err = coerce.Bool(map[string]string{})
		chk.True(errors.Is(err, coerce.ErrUnsupported))
		chk.Equal(false, dst)
	}
}

func TestBoolFromNil(t *testing.T) {
	chk := assert.New(t)
	var dst bool
	var err error
	dst, err = coerce.Bool(nil)
	chk.NoError(err)
	chk.Equal(false, dst)
}

func TestBoolFromPtr(t *testing.T) {
	n := (*bool)(nil)
	nn := &n
	bf, bt := false, true
	pf, pt := &bf, &bt
	ppf, ppt := &pf, &pt
	pppf, pppt := &ppf, &ppt
	tests := BoolTests{
		"nil": {
			To: n, Expect: false,
		},
		"*nil": {
			To: nn, Expect: false,
		},
		"pppf": {
			To: pppf, Expect: bf,
		},
		"pppt": {
			To: pppt, Expect: bt,
		},
	}
	tests.Run(t)
}

func TestBoolFromSlice(t *testing.T) {
	tests := BoolTests{
		"nil slice": {
			To: []string(nil), Expect: false,
		},
		"slice": {
			To: []interface{}{42, "true"}, Expect: true,
		},
		"slice kind": {
			To: []interface{}{42, "true", F32(1)}, Expect: true,
		},
	}
	tests.Run(t)
}

func TestBoolFromUint(t *testing.T) {
	tests := BoolTests{
		"max8": {
			To: uint8(math.MaxInt8), Expect: true,
		},
		"max16": {
			To: uint16(math.MaxInt16), Expect: true,
		},
		"max32": {
			To: uint32(math.MaxInt32), Expect: true,
		},
		"max64": {
			To: uint64(math.MaxInt64), Expect: true,
		},
		"uint64 max": {
			To: uint64(math.MaxUint64), Expect: true,
		},
		"0": {
			To: uint(0), Expect: false,
		},

		"U8(max8)": {
			To: U8(math.MaxInt8), Expect: true,
		},
		"U16(max16)": {
			To: U16(math.MaxInt16), Expect: true,
		},
		"U32(max32)": {
			To: U32(math.MaxInt32), Expect: true,
		},
		"U64(max64)": {
			To: U64(math.MaxInt64), Expect: true,
		},
		"U64(uint64 max)": {
			To: U64(math.MaxUint64), Expect: true,
		},
		"U(0)": {
			To: UU(0), Expect: false,
		},
	}
	tests.Run(t)
}

func TestBoolFromString(t *testing.T) {
	String := func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}
	tests := BoolTests{
		"max8": {
			To: String(math.MaxInt8), Error: coerce.ErrInvalid,
		},
		"max16": {
			To: String(math.MaxInt16), Error: coerce.ErrInvalid,
		},
		"1": {
			To: "1", Expect: true,
		},
		"0": {
			To: "0", Expect: false,
		},
		"true": {
			To: "true", Expect: true,
		},
		"false": {
			To: "false", Expect: false,
		},
		"True": {
			To: "True", Expect: true,
		},
		"False": {
			To: "False", Expect: false,
		},

		"S(max8)": {
			To: S(String(math.MaxInt8)), Error: coerce.ErrInvalid,
		},
		"S(max16)": {
			To: S(String(math.MaxInt16)), Error: coerce.ErrInvalid,
		},
		"S(1)": {
			To: S("1"), Expect: true,
		},
		"S(0)": {
			To: S("0"), Expect: false,
		},
		"S(true)": {
			To: S("true"), Expect: true,
		},
		"S(false)": {
			To: S("false"), Expect: false,
		},
		"S(True)": {
			To: S("True"), Expect: true,
		},
		"S(False)": {
			To: S("False"), Expect: false,
		},
	}
	tests.Run(t)
}
