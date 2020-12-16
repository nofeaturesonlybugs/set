package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleBoolSlice_Set() {
	b := []bool{}
	fmt.Println(b)
	v := set.V(&b)

	// A slice of specific types (float32)
	v.To([]float32{1, 0})
	fmt.Println(b)

	// A slice of int64
	v.To([]int64{0, 1, 0})
	fmt.Println(b)

	// Mixed interface{} slice.
	v.To([]interface{}{"1", "False", true, 0})
	fmt.Println(b)

	// Output: []
	// [true false]
	// [false true false]
	// [true false true false]
}

func ExampleBoolSlice_Set_createsCopy() {
	slice := []bool{true, false, true}
	var dest []bool
	set.V(&dest).To(slice)

	fmt.Println(slice[1])
	dest[1] = true
	fmt.Println(slice[1])

	// Output: false
	// false
}
