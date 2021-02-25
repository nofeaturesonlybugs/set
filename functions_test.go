package set_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
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
	{
		v, ok := set.Writable(reflect.Value{})
		chk.Equal(false, ok)
		chk.Equal(false, v.IsValid())
	}
}

func ExampleWritable() {
	var value, writable reflect.Value
	var ok bool
	var s string
	var sp *string

	value = reflect.ValueOf(s)
	writable, ok = set.Writable(value)
	fmt.Printf("ok= %v\n", ok)

	value = reflect.ValueOf(sp)
	writable, ok = set.Writable(value)
	fmt.Printf("ok= %v\n", ok)

	value = reflect.ValueOf(&sp)
	writable, ok = set.Writable(value)
	writable.SetString("Hello")
	fmt.Printf("ok= %v sp= %v\n", ok, *sp)

	// Output: ok= false
	// ok= false
	// ok= true sp= Hello
}
