package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleFloatSlice_Set() {
	f := []float32{}
	fmt.Println(f) // []
	v := set.V(&f)

	// A slice of specific types (string)
	v.To([]string{"1.1", "0.5"})
	fmt.Println(f)

	// A slice of int64
	v.To([]int64{0, 1, 0})
	fmt.Println(f)

	// Mixed interface{} slice.
	v.To([]interface{}{"3.14", "-5", false, true})
	fmt.Println(f)

	// Output: []
	// [1.1 0.5]
	// [0 1 0]
	// [3.14 -5 0 1]
}
