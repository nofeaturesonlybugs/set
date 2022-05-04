package set_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleValue_Rebind() {
	// Once a Value has been created you can swap out the value it mutates
	// by calling Rebind.  This yields better performance in tight loops where
	// you intend to mutate many instances of the same type.

	values := []string{"3.14", "false", "true", "5"}
	slice := make([]int, 4)

	v := set.V(&slice[0])
	for k, str := range values {
		v.Rebind(&slice[k])
		if err := v.To(str); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println(slice)

	// Output: [3 0 1 5]
}
func ExampleValue_Rebind_panic() {
	// Value.Rebind panics if the new instance is not the same type.

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	var a int
	var b string

	v := set.V(&a)
	v.Rebind(&b)

	// Output: mismatching types during Rebind; have *int and got *string
}

func ExampleValue_Rebind_reflectValue() {
	// As a convenience Value.Rebind will accept reflect.Value
	// if-and-only-if the reflect.Value is holding a type compatible
	// with the Value.

	var a, b int

	v, rv := set.V(&a), reflect.ValueOf(&b)

	if err := v.To("42"); err != nil {
		fmt.Println(err)
		return
	}

	v.Rebind(rv)
	if err := v.To("24"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("a=%v b=%v", a, b)

	// Output: a=42 b=24
}

func ExampleValue_Rebind_value() {
	// As a convenience Value.Rebind will accept another Value
	// as long as the internal types are compatible.

	var a, b int

	av, bv := set.V(&a), set.V(&b)

	if err := av.To("42"); err != nil {
		fmt.Println(err)
		return
	}

	av.Rebind(bv)
	if err := av.To("24"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("a=%v b=%v", a, b)

	// Output: a=42 b=24
}
