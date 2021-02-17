package set_test

import (
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

func TestWritable(t *testing.T) {
	chk := assert.New(t)
	//
	{
		var b bool
		//
		v, info, ok := set.Writable(reflect.ValueOf(b))
		chk.Equal(false, ok)
		chk.Equal(reflect.Bool, info.Kind)
		chk.Equal(true, v.IsZero())
		//
		v, info, ok = set.Writable(reflect.ValueOf(&b))
		chk.Equal(true, ok)
		chk.Equal(reflect.Bool, info.Kind)
		chk.Equal(true, v.IsZero())
		chk.Equal(true, v.CanSet())
		v.SetBool(true)
		chk.Equal(true, b)
	}
	{
		var bp *bool
		//
		v, info, ok := set.Writable(reflect.ValueOf(bp))
		chk.Equal(false, ok)
		chk.Equal(reflect.Bool, info.Kind)
		chk.Equal(false, v.IsValid())
		//
		v, info, ok = set.Writable(reflect.ValueOf(&bp))
		chk.Equal(true, ok)
		chk.Equal(reflect.Bool, info.Kind)
		chk.Equal(true, v.IsZero())
		chk.Equal(true, v.CanSet())
		v.SetBool(true)
		chk.Equal(true, *bp)
	}
	{
		var b bool
		bp := &b
		//
		v, info, ok := set.Writable(reflect.ValueOf(bp))
		chk.Equal(true, ok)
		chk.Equal(reflect.Bool, info.Kind)
		chk.Equal(true, v.IsValid())
		v.SetBool(true)
		chk.Equal(true, b)
		chk.Equal(true, *bp)
	}
}
