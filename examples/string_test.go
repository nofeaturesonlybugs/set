package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleString_Set_bool() {
	s1, s2 := "", ""
	fmt.Println(s1, s2, "(empty)")
	v1, v2 := set.V(&s1), set.V(&s2)

	v1.To(true)
	v2.To(true)
	fmt.Println(s1, s2)

	v1.To(false)
	v2.To(false)
	fmt.Println(s1, s2)

	// Output: (empty)
	// true true
	// false false
}

func ExampleString_Set_float() {
	s1, s2 := "", ""
	fmt.Println(s1, s2, "(empty)")
	v1, v2 := set.V(&s1), set.V(&s2)

	v1.To(float32(1.15))
	v2.To(float64(1.59))
	fmt.Println(s1, s2)

	v1.To(float32(0))
	v2.To(float64(0))
	fmt.Println(s1, s2)

	v1.To(float32(-1.15))
	v2.To(float64(-1.59))
	fmt.Println(s1, s2)

	// Output: (empty)
	// 1.15 1.59
	// 0 0
	// -1.15 -1.59
}

func ExampleString_Set_int() {
	s1, s2 := "", ""
	fmt.Println(s1, s2, "(empty)")
	v1, v2 := set.V(&s1), set.V(&s2)

	v1.To(int32(1))
	v2.To(int64(1))
	fmt.Println(s1, s2)

	v1.To(int32(0))
	v2.To(int64(0))
	fmt.Println(s1, s2)

	v1.To(int32(999))
	v2.To(int64(999))
	fmt.Println(s1, s2)

	v1.To(int32(-999))
	v2.To(int64(-999))
	fmt.Println(s1, s2)

	// Output: (empty)
	// 1 1
	// 0 0
	// 999 999
	// -999 -999
}

func ExampleString_Set_uint() {
	s1, s2 := "", ""
	fmt.Println(s1, s2, "(empty)")
	v1, v2 := set.V(&s1), set.V(&s2)

	v1.To(uint32(1))
	v2.To(uint64(1))
	fmt.Println(s1, s2)

	v1.To(uint32(0))
	v2.To(uint64(0))
	fmt.Println(s1, s2)

	v1.To(uint32(999))
	v2.To(uint64(999))
	fmt.Println(s1, s2)

	// Output: (empty)
	// 1 1
	// 0 0
	// 999 999
}
