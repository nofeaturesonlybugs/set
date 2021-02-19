package set_test

import (
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

func TestTypeList(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct{}
	type B struct{}
	a := set.NewTypeList(A{})
	b := set.NewTypeList(B{})
	//
	chk.Equal(true, a.Has(reflect.TypeOf(A{})))
	chk.Equal(true, b.Has(reflect.TypeOf(B{})))
	//
	a.Merge(b)
	chk.Equal(true, a.Has(reflect.TypeOf(A{})))
	chk.Equal(true, a.Has(reflect.TypeOf(B{})))
}
