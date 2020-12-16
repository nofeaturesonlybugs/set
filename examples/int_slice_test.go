package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleIntSlice_Set() {
	i := []int32{}
	fmt.Println(i) // []
	v := set.V(&i)

	// A slice of specific types (string)
	v.To([]string{"1.1", "0.5"})
	fmt.Println(i)

	// A slice of int64
	v.To([]float64{0, 1, 0})
	fmt.Println(i)

	// Mixed interface{} slice.
	v.To([]interface{}{"3.14", "-5", false, true})
	fmt.Println(i)

	// Output: []
	// [1 0]
	// [0 1 0]
	// [3 -5 0 1]
}
