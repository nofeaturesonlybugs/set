package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleInt_Set_bool() {
	i1, i2 := int32(0), int64(0)
	fmt.Println(i1, i2)
	v1, v2 := set.V(&i1), set.V(&i2)

	v1.To(true)
	v2.To(true)
	fmt.Println(i1, i2)

	v1.To(false)
	v2.To(false)
	fmt.Println(i1, i2)

	// Output: 0 0
	// 1 1
	// 0 0
}

func ExampleInt_Set_float() {
	i1, i2 := int32(0), int64(0)
	fmt.Println(i1, i2)
	v1, v2 := set.V(&i1), set.V(&i2)

	v1.To(float32(1.15))
	v2.To(float64(1.59))
	fmt.Println(i1, i2)

	v1.To(float32(0))
	v2.To(float64(0))
	fmt.Println(i1, i2)

	v1.To(float32(-1.15))
	v2.To(float64(-1.59))
	fmt.Println(i1, i2)

	// Output: 0 0
	// 1 1
	// 0 0
	// -1 -1
}

func ExampleInt_Set_uint() {
	i1, i2 := int32(0), int64(0)
	fmt.Println(i1, i2)
	v1, v2 := set.V(&i1), set.V(&i2)

	v1.To(uint32(1))
	v2.To(uint64(1))
	fmt.Println(i1, i2)

	v1.To(uint32(0))
	v2.To(uint64(0))
	fmt.Println(i1, i2)

	v1.To(uint32(999))
	v2.To(uint64(999))
	fmt.Println(i1, i2)

	// Output: 0 0
	// 1 1
	// 0 0
	// 999 999
}

func ExampleInt_Set_string() {
	i1, i2 := int32(0), int64(0)
	fmt.Println(i1, i2)
	v1, v2 := set.V(&i1), set.V(&i2)

	v1.To("1")
	v2.To("1")
	fmt.Println(i1, i2)

	v1.To("-1")
	v2.To("-1")
	fmt.Println(i1, i2)

	v1.To("1.59")
	v2.To("1.59")
	fmt.Println(i1, i2)

	v1.To("-3.14")
	v2.To("-3.14")
	fmt.Println(i1, i2)

	v1.To("0")
	v2.To("0")
	fmt.Println(i1, i2)

	// Output: 0 0
	// 1 1
	// -1 -1
	// 1 1
	// -3 -3
	// 0 0
}
