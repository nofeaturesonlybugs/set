package set

import (
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/set/assert"
)

func TestCoerceToBool(t *testing.T) {
	chk := assert.New(t)
	//
	var target, value reflect.Value
	var err error
	b := false
	target = reflect.ValueOf(&b)
	{
		// float-to-bool
		//
		for _, v := range []float32{1, -1, 3.14, -3.14, 0.5, -0.5} {
			value = reflect.ValueOf(float32(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(float64(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
		}
		//
		value = reflect.ValueOf(float32(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
		//
		value = reflect.ValueOf(float64(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		// int-to-bool
		//
		for _, v := range []int8{1, -1, 32, -32, 127, -127} {
			value = reflect.ValueOf(int(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(int8(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(int16(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(int32(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(int64(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
		}
		//
		value = reflect.ValueOf(int(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
		//
		value = reflect.ValueOf(int8(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		// uint-to-bool
		//
		for _, v := range []uint8{1, 32, 127, 255} {
			value = reflect.ValueOf(uint(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(uint8(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(uint16(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(uint32(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
			value = reflect.ValueOf(uint64(v))
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
		}
		//
		value = reflect.ValueOf(uint(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
		//
		value = reflect.ValueOf(uint8(0))
		err = coerce(reflect.Indirect(target), value)
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		// string-to-bool
		//
		for _, v := range []string{"1", "True", "TRUE", "true"} {
			value = reflect.ValueOf(v)
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(true, b)
			//
		}
		for _, v := range []string{"0", "False", "FALSE", "false"} {
			value = reflect.ValueOf(v)
			err = coerce(reflect.Indirect(target), value)
			chk.NoError(err)
			chk.Equal(false, b)
			//
		}
		for _, v := range []string{"tRuE", "asdf"} {
			b = true
			chk.Equal(true, b)
			value = reflect.ValueOf(v)
			err = coerce(reflect.Indirect(target), value)
			chk.Error(err)
			chk.Equal(false, b)
			//
		}
	}
}

func TestCoerceToFloat(t *testing.T) {
	chk := assert.New(t)
	//
	var target, value reflect.Value
	var err error
	f32, f64 := float32(0), float64(0)
	for _, ptr := range []interface{}{&f32, &f64} {
		// We test against both targets, float32 and float64
		target = reflect.Indirect(reflect.ValueOf(ptr))
		//
		// bool-to-float
		err = coerce(target, reflect.ValueOf(true))
		chk.NoError(err)
		if ptr == &f32 {
			chk.Equal(float32(1), *(ptr.(*float32)))
		} else {
			chk.Equal(float64(1), *(ptr.(*float64)))
		}
		//
		err = coerce(target, reflect.ValueOf(false))
		chk.NoError(err)
		if ptr == &f32 {
			chk.Equal(float32(0), *(ptr.(*float32)))
		} else {
			chk.Equal(float64(0), *(ptr.(*float64)))
		}
		//
		// int-to-float
		for _, v := range []int{0, 1, 2, 16, 32, -1, -2, -16, -32} {
			// int
			value = reflect.ValueOf(v)
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), int(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), int(*(ptr.(*float64))))
			}
			// int8
			value = reflect.ValueOf(int8(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), int8(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), int8(*(ptr.(*float64))))
			}
			// int16
			value = reflect.ValueOf(int16(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), int16(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), int16(*(ptr.(*float64))))
			}
			// int32
			value = reflect.ValueOf(int32(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), int32(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), int32(*(ptr.(*float64))))
			}
			// int64
			value = reflect.ValueOf(int64(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), int64(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), int64(*(ptr.(*float64))))
			}
		} //
		// uint-to-float
		for _, v := range []uint{0, 1, 2, 16, 32} {
			// int
			value = reflect.ValueOf(v)
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), uint(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), uint(*(ptr.(*float64))))
			}
			// int8
			value = reflect.ValueOf(uint8(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), uint8(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), uint8(*(ptr.(*float64))))
			}
			// int16
			value = reflect.ValueOf(uint16(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), uint16(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), uint16(*(ptr.(*float64))))
			}
			// int32
			value = reflect.ValueOf(uint32(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), uint32(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), uint32(*(ptr.(*float64))))
			}
			// int64
			value = reflect.ValueOf(uint64(v))
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.Equal(value.Interface(), uint64(*(ptr.(*float32))))
			} else {
				chk.Equal(value.Interface(), uint64(*(ptr.(*float64))))
			}
		}
		//
		// string-to-float
		for _, v := range []struct {
			S string
			E float32
		}{
			{"0", 0},
			{"3.14", 3.14},
			{"-3.14", -3.14},
		} {
			value = reflect.ValueOf(v.S)
			err = coerce(target, value)
			chk.NoError(err)
			if ptr == &f32 {
				chk.InDelta(float32(v.E), *(ptr.(*float32)), 0.01)
			} else {
				chk.InDelta(float64(v.E), *(ptr.(*float64)), 0.01)
			}
		}
		for _, v := range []string{"asdf", "failure"} {
			value = reflect.ValueOf(v)
			err = coerce(target, value)
			chk.Error(err)
		}
	}
}

func TestCoerceToInt(t *testing.T) {
	chk := assert.New(t)
	//
	var target, value reflect.Value
	var err error
	i, i8, i16, i32, i64 := int(0), int8(0), int16(0), int32(0), int64(0)
	checkEqual := func(expect int8, ptr interface{}) {
		if ptr == &i {
			chk.Equal(int(expect), *(ptr.(*int)))
		} else if ptr == &i8 {
			chk.Equal(int8(expect), *(ptr.(*int8)))
		} else if ptr == &i16 {
			chk.Equal(int16(expect), *(ptr.(*int16)))
		} else if ptr == &i32 {
			chk.Equal(int32(expect), *(ptr.(*int32)))
		} else if ptr == &i64 {
			chk.Equal(int64(expect), *(ptr.(*int64)))
		}
	}
	for _, ptr := range []interface{}{&i, &i8, &i16, &i32, &i64} {
		target = reflect.Indirect(reflect.ValueOf(ptr))
		//
		// bool-to-int
		value = reflect.ValueOf(true)
		err = coerce(target, value)
		chk.NoError(err)
		checkEqual(1, ptr)
		//
		value = reflect.ValueOf(false)
		err = coerce(target, value)
		chk.NoError(err)
		checkEqual(0, ptr)
		//
		// float-to-int
		for _, v := range []struct {
			F float32
			E int8
		}{
			{0, 0},
			{3.14, 3}, {-3.14, -3},
			{4.75, 4}, {-4.75, -4},
		} {
			value = reflect.ValueOf(v.F)
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(v.E, ptr)
		}
		//
		// uint-to-int
		for _, v := range []uint8{0, 1, 2, 16, 32} {
			value = reflect.ValueOf(uint(v))
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(int8(v), ptr)
			//
			value = reflect.ValueOf(uint8(v))
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(int8(v), ptr)
			//
			value = reflect.ValueOf(uint16(v))
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(int8(v), ptr)
			//
			value = reflect.ValueOf(uint32(v))
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(int8(v), ptr)
			//
			value = reflect.ValueOf(uint64(v))
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(int8(v), ptr)
		}
		//
		// string-to-int
		for _, v := range []struct {
			S string
			E int8
		}{
			{"0", 0},
			{"1", 1}, {"-1", -1},
			{"2", 2}, {"-2", -2},
			{"32", 32}, {"-32", -32},
			{"3.14", 3}, {"-3.14", -3},
		} {
			value = reflect.ValueOf(v.S)
			err = coerce(target, value)
			chk.NoError(err)
			checkEqual(v.E, ptr)
		}
	}
}

func TestCoerceToUint(t *testing.T) {
	chk := assert.New(t)
	//
	var target, value reflect.Value
	var err error
	ui, ui8, ui16, ui32, ui64 := uint(0), uint8(0), uint16(0), uint32(0), uint64(0)
	checkEqual := func(expect uint8, ptr interface{}) {
		if ptr == &ui {
			chk.Equal(uint(expect), *(ptr.(*uint)))
		} else if ptr == &ui8 {
			chk.Equal(uint8(expect), *(ptr.(*uint8)))
		} else if ptr == &ui16 {
			chk.Equal(uint16(expect), *(ptr.(*uint16)))
		} else if ptr == &ui32 {
			chk.Equal(uint32(expect), *(ptr.(*uint32)))
		} else if ptr == &ui64 {
			chk.Equal(uint64(expect), *(ptr.(*uint64)))
		}
	}
	for _, ptr := range []interface{}{&ui, &ui8, &ui16, &ui32, &ui64} {
		target = reflect.Indirect(reflect.ValueOf(ptr))
		//
		// bool-to-uint
		value = reflect.ValueOf(true)
		err = coerce(target, value)
		chk.NoError(err)
		checkEqual(1, ptr)
		//
		value = reflect.ValueOf(false)
		err = coerce(target, value)
		chk.NoError(err)
		checkEqual(0, ptr)
		//
		// float-to-uint
		for _, v := range []struct {
			F     float32
			E     uint8
			Error bool
		}{
			{0, 0, false},
			{3.14, 3, false}, {-3.14, 0, true},
			{4.75, 4, false}, {-4.75, 0, true},
		} {
			value = reflect.ValueOf(v.F)
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
		}
		//
		// int-to-uint
		for _, v := range []struct {
			I     int
			E     uint8
			Error bool
		}{
			{0, 0, false},
		} {

			value = reflect.ValueOf(int(v.I))
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
			//
			value = reflect.ValueOf(int8(v.I))
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
			//
			value = reflect.ValueOf(int16(v.I))
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
			//
			value = reflect.ValueOf(int32(v.I))
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
			//
			value = reflect.ValueOf(int64(v.I))
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
		}
		//
		// string-to-uint
		for _, v := range []struct {
			S     string
			E     uint8
			Error bool
		}{
			{"0", 0, false},
			{"1", 1, false}, {"-1", 0, true},
			{"2", 2, false}, {"-2", 0, true},
			{"32", 32, false}, {"-32", 0, true},
			{"3.14", 3, false}, {"-3.14", 0, true},
		} {
			value = reflect.ValueOf(v.S)
			err = coerce(target, value)
			if v.Error == false {
				chk.NoError(err)
			} else {
				chk.Error(err)
			}
			checkEqual(v.E, ptr)
		}
	}
}

func TestCoerceToString(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	var s string
	target := reflect.Indirect(reflect.ValueOf(&s))
	for _, v := range []struct {
		V     interface{}
		E     string
		Error bool
	}{
		{0, "0", false},
		{uint64(2342342), "2342342", false},
		{int32(2342342), "2342342", false},
		{float64(3.14), "3.14", false},
		{true, "true", false},
		{false, "false", false},
		{"Hello", "", true},
	} {
		s = ""
		err = coerce(target, reflect.ValueOf(v.V))
		if v.Error {
			chk.Error(err)
			chk.Equal("", s)
		} else {
			chk.NoError(err)
			chk.Equal(v.E, s)
		}
	}
}
