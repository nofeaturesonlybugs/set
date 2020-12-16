package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleFloat_Set_bool() {
	f1, f2 := float32(0), float64(0)
	fmt.Println(f1, f2)
	v1, v2 := set.V(&f1), set.V(&f2)

	v1.To(true)
	v2.To(true)
	fmt.Println(f1, f2)

	v1.To(false)
	v2.To(false)
	fmt.Println(f1, f2)

	// Output: 0 0
	// 1 1
	// 0 0
}

func ExampleFloat_Set_int() {
	f1, f2 := float32(0), float64(0)
	fmt.Println(f1, f2)
	v1, v2 := set.V(&f1), set.V(&f2)

	v1.To(int32(1))
	v2.To(int64(1))
	fmt.Println(f1, f2)

	v1.To(int32(-3))
	v2.To(int64(-3))
	fmt.Println(f1, f2)

	// Output: 0 0
	// 1 1
	// -3 -3
}

func ExampleFloat_Set_uint() {
	f1, f2 := float32(0), float64(0)
	fmt.Println(f1, f2)
	v1, v2 := set.V(&f1), set.V(&f2)

	v1.To(uint32(1))
	v2.To(uint64(1))
	fmt.Println(f1, f2)

	v1.To(uint32(3))
	v2.To(uint64(3))
	fmt.Println(f1, f2)

	// Output: 0 0
	// 1 1
	// 3 3
}

func ExampleFloat_Set_string() {
	f1, f2 := float32(0), float64(0)
	fmt.Println(f1, f2)
	v1, v2 := set.V(&f1), set.V(&f2)

	v1.To("1")
	v2.To("1")
	fmt.Println(f1, f2)

	v1.To("3.14")
	v2.To("3.14")
	fmt.Println(f1, f2)

	v1.To("-3.59")
	v2.To("-3.59")
	fmt.Println(f1, f2)

	// Output: 0 0
	// 1 1
	// 3.14 3.14
	// -3.59 -3.59
}
