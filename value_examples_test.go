package set_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleValue_To() {
	// Value wraps around scalars or slices and the To method performs type coercion.
	//
	// When calling set.V(a) `a` must represent an addressable or writable value similarly
	// to invoking deserializers like json.Unmarshal([]byte, a).

	// Fails because not addressable
	var b bool
	err := set.V(b).To("true") // Should have been &b -- address of b
	fmt.Println(err)

	// Succeed because we pass addresses of destinations.
	var s string
	var n int
	var u8 uint8

	err = set.V(&s).To(3.14)
	fmt.Println("s", s, err)

	err = set.V(&n).To("42")
	fmt.Println("n", n, err)

	err = set.V(&u8).To("27")
	fmt.Println("u8", u8, err)

	// When passing the address of a nil ptr it will be created.
	var nptr *int                // nptr == nil
	err = set.V(&nptr).To("100") // nptr != nil now
	fmt.Println("nptr", *nptr, err)

	// If a pointer already points at something you can pass it directly.
	sptr := &s // sptr != nil and points at s
	err = set.V(sptr).To("Something")
	fmt.Println("sptr", *sptr, err, "s", s)

	// new(T) works the same as sptr above.
	f32 := new(float32)
	err = set.V(f32).To("3")
	fmt.Println("f32", *f32, err)

	// Output: set: Value.To: read only value: bool is not writable: hint=[call to set.V(bool) should have been set.V(*bool)]
	// s 3.14 <nil>
	// n 42 <nil>
	// u8 27 <nil>
	// nptr 100 <nil>
	// sptr Something <nil> s Something
	// f32 3 <nil>
}

func ExampleValue_To_scalarToSlice() {
	// This example demonstrates how scalars are coerced to slices.
	// Given
	// 	↪ var t []T    // A target slice
	// 	↪ var s S      // A scalar value
	// Then
	//	↪ _ = set.V(&t).To(s)
	// Yields
	//	↪ t == []T{ s } // A slice with a single element

	var n []int
	var s []string

	err := set.V(&n).To("9")
	fmt.Println("n", n, err)

	err = set.V(&s).To(1234)
	fmt.Println("s", s, err)

	// When the incoming value can not be type coerced the slice will be an empty slice.
	err = set.V(&n).To("Hello")
	fmt.Println("n", n == nil, err == nil)

	// Output: n [9] <nil>
	// s [1234] <nil>
	// n true false
}

func ExampleValue_To_sliceToScalar() {
	// This example demonstrates how slices are coerced to scalars.
	// Given
	// 	↪ var t T    // A target scalar
	// 	↪ var s []S  // A slice value
	// Then
	//	↪ _ = set.V(&t).To(s)
	// Yields
	//	↪ t == s[ len(s) -1 ]    // t is set to last value in slice

	var n int
	var s string

	err := set.V(&n).To([]interface{}{9, "3.14", false, "next value wins", "9999"})
	fmt.Println("n", n, err)

	err = set.V(&s).To([]interface{}{9, "3.14", false, "next value wins", "I win!"})
	fmt.Println("s", s, err)

	// When the incoming slice can not be type coerced the scalar will be zero value.
	err = set.V(&n).To([]interface{}{9, "3.14", false, "next value wins", "9999", "Fred"})
	fmt.Println("n", n == 0, err == nil)

	// When the incoming slice is nil or empty the scalar will be zero value.
	err = set.V(&s).To([]int(nil)) // Sets s to "" -- empty string!
	fmt.Println("s", s, err)

	// Output: n 9999 <nil>
	// s I win! <nil>
	// n true false
	// s  <nil>
}

func ExampleValue_To_sliceToSlice() {
	// This example demonstrates how slices are coerced to slices.
	// Given
	// 	↪ var t []T  // A target slice
	// 	↪ var s []S  // A slice value
	// Then
	//	↪ _ = set.V(&t).To(s)
	// Yields
	//	↪ t == s // t is a slice equal in length to s with elements type coerced to T
	// However
	//	↪ t is a copy of s

	var n []int
	var s []string

	values := []interface{}{9, "3.14", false, "9999"}

	err := set.V(&n).To(values)
	fmt.Println("n", n, err)

	err = set.V(&s).To(values)
	fmt.Println("s", s, err)

	// If any element can not be coerced the target (or dest) will be zero value slice.
	values = append(values, "Hello") // "Hello" can not be coereced to int
	err = set.V(&n).To(values)
	fmt.Println("n", n == nil, err == nil)

	// When dealing with slices the target (or dest) is always a copy.
	m := []int{2, 4, 6, 8}
	err = set.V(&n).To(m) // Even though m and n are same type n will be a copy
	fmt.Println("n", n, err)
	m[1] = -4 // Change element in m
	fmt.Println("m", m, "n", n)

	// Output: n [9 3 0 9999] <nil>
	// s [9 3.14 false 9999] <nil>
	// n true false
	// n [2 4 6 8] <nil>
	// m [2 -4 6 8] n [2 4 6 8]
}

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
