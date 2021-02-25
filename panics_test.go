package set_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestPanics_Append(t *testing.T) {
	chk := assert.New(t)
	//
	type Person struct {
		Name string
		Age  int
	}
	{
		var dest []Person
		v := set.V(&dest)
		times := 10
		for k := 0; k < times; k++ {
			elem, err := v.NewElem()
			chk.NoError(err)
			//
			field, err := elem.FieldByIndexAsValue([]int{0})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			field, err = elem.FieldByIndexAsValue([]int{1})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			set.Panics.Append(v, elem)
		}
		chk.Equal(times, len(dest))
		for z := 0; z < times; z++ {
			chk.Equal(fmt.Sprintf("%v", z), dest[z].Name)
			chk.Equal(z, dest[z].Age)
		}
	}
	{
		var dest []*Person
		v := set.V(&dest)
		times := 10
		for k := 0; k < times; k++ {
			elem, err := v.NewElem()
			chk.NoError(err)
			//
			field, err := elem.FieldByIndexAsValue([]int{0})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			field, err = elem.FieldByIndexAsValue([]int{1})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			set.Panics.Append(v, elem)
		}
		chk.Equal(times, len(dest))
		for z := 0; z < times; z++ {
			chk.Equal(fmt.Sprintf("%v", z), dest[z].Name)
			chk.Equal(z, dest[z].Age)
		}
	}
	{
		var dest []*****Person
		v := set.V(&dest)
		times := 10
		for k := 0; k < times; k++ {
			elem, err := v.NewElem()
			chk.NoError(err)
			//
			field, err := elem.FieldByIndexAsValue([]int{0})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			field, err = elem.FieldByIndexAsValue([]int{1})
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(k)
			chk.NoError(err)
			//
			set.Panics.Append(v, elem)
		}
		chk.Equal(times, len(dest))
		for z := 0; z < times; z++ {
			chk.Equal(fmt.Sprintf("%v", z), (*****dest[z]).Name)
			chk.Equal(z, (*****dest[z]).Age)
		}
	}
}
