package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleStringSlice_Set() {
	s := []string{}
	fmt.Println(s) // []
	v := set.V(&s)

	// A slice of specific types (string)
	v.To([]string{"1.1", "0.5"})
	fmt.Println(s)

	// A slice of int64
	v.To([]float64{0, 1, 0})
	fmt.Println(s)

	// Mixed interface{} slice.
	v.To([]interface{}{3.14, 5, false, true})
	fmt.Println(s)

	// Output: []
	// [1.1 0.5]
	// [0 1 0]
	// [3.14 5 false true]
}
