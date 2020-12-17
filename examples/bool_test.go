package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleBool_Set_floats() {
	b1, b2 := false, false
	fmt.Println(b1, b2)
	v1, v2 := set.V(&b1), set.V(&b2)

	v1.To(float32(1))
	v2.To(float64(1))
	fmt.Println(b1, b2)

	v1.To(float32(0))
	v2.To(float64(0))
	fmt.Println(b1, b2)

	v1.To(float32(3.14))
	v2.To(float64(3.14))
	fmt.Println(b1, b2)

	// Output: false false
	// true true
	// false false
	// true true
}

func ExampleBool_Set_ints() {
	b1, b2 := false, false
	fmt.Println(b1, b2)
	v1, v2 := set.V(&b1), set.V(&b2)

	v1.To(int32(1))
	v2.To(int64(-1))
	fmt.Println(b1, b2)

	v1.To(int32(0))
	v2.To(int64(0))
	fmt.Println(b1, b2)

	v1.To(int8(3))
	v2.To(int16(-3))
	fmt.Println(b1, b2)

	// Output: false false
	// true true
	// false false
	// true true
}

func ExampleBool_Set_uints() {
	b1, b2 := false, false
	fmt.Println(b1, b2)
	v1, v2 := set.V(&b1), set.V(&b2)

	v1.To(uint32(1))
	v2.To(uint64(1))
	fmt.Println(b1, b2)

	v1.To(uint32(0))
	v2.To(uint64(0))
	fmt.Println(b1, b2)

	v1.To(uint8(3))
	v2.To(uint16(3))
	fmt.Println(b1, b2)

	// Output: false false
	// true true
	// false false
	// true true
}

func ExampleBool_Set_strings() {
	b1, b2 := false, false
	fmt.Println(b1, b2)
	v1, v2 := set.V(&b1), set.V(&b2)

	v1.To("True")
	v2.To("TRUE")
	fmt.Println(b1, b2)

	v1.To("0")
	v2.To("false")
	fmt.Println(b1, b2)

	v1.To("true")
	v2.To("1")
	fmt.Println(b1, b2)

	// Output: false false
	// true true
	// false false
	// true true
}
