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
		v, ok := set.Writable(reflect.ValueOf(b))
		chk.Equal(false, ok)
		chk.Equal(true, v.IsZero())
		//
		v, ok = set.Writable(reflect.ValueOf(&b))
		chk.Equal(true, ok)
		chk.Equal(true, v.IsZero())
		chk.Equal(true, v.CanSet())
		v.SetBool(true)
		chk.Equal(true, b)
	}
	{
		var bp *bool
		//
		v, ok := set.Writable(reflect.ValueOf(bp))
		chk.Equal(false, ok)
		chk.Equal(false, v.IsValid())
		//
		v, ok = set.Writable(reflect.ValueOf(&bp))
		chk.Equal(true, ok)
		chk.Equal(true, v.IsZero())
		chk.Equal(true, v.CanSet())
		v.SetBool(true)
		chk.Equal(true, *bp)
	}
	{
		var b bool
		bp := &b
		//
		v, ok := set.Writable(reflect.ValueOf(bp))
		chk.Equal(true, ok)
		chk.Equal(true, v.IsValid())
		v.SetBool(true)
		chk.Equal(true, b)
		chk.Equal(true, *bp)
	}
}
