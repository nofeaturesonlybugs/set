package set_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleBoundMapping_Rebind() {
	// Once a BoundMapping has been created you can swap out the value it mutates
	// by calling Rebind.  This yields better performance in tight loops where
	// you intend to mutate many instances of the same type.

	m := &set.Mapper{}

	type S struct {
		Str string
		Num int
	}

	values := []map[string]interface{}{
		{"Str": "First", "Num": 1},
		{"Str": "Second", "Num": 2},
		{"Str": "Third", "Num": 3},
	}
	slice := make([]S, 3)

	b, _ := m.Bind(&slice[0]) // We can bind on the first value to get a BoundMapping
	for k := range slice {
		b.Rebind(&slice[k]) // The BoundMapping now affects slice[k]
		for key, value := range values[k] {
			_ = b.Set(key, value) // error ignored for brevity
		}
	}

	fmt.Println(slice[0].Num, slice[0].Str)
	fmt.Println(slice[1].Num, slice[1].Str)
	fmt.Println(slice[2].Num, slice[2].Str)

	// Output: 1 First
	// 2 Second
	// 3 Third
}

func ExampleBoundMapping_Rebind_panic() {
	// BoundMapping.Rebind panics if the new instance is not the same type.

	m := &set.Mapper{}

	type S struct {
		Str string
	}
	type Different struct {
		Str string
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	var s S
	var d Different

	b, _ := m.Bind(&s)
	b.Rebind(&d)

	// Output: mismatching types during Rebind; have *set_test.S and got *set_test.Different
}

func ExampleBoundMapping_Rebind_reflectValue() {
	// As a convenience BoundMapping.Rebind can be called with an instance of reflect.Value
	// if-and-only-if the reflect.Value is holding a type compatible with the BoundMapping.

	m := &set.Mapper{}

	type S struct {
		Str string
		Num int
	}
	var s S

	// Errors ignored for brevity.
	b, _ := m.Bind(&s)
	_ = b.Set("Num", 42) // b.Set errors ignored for brevity
	_ = b.Set("Str", "Hello")

	rv := reflect.New(reflect.TypeOf(s)) // reflect.New creates a *S which is the type
	b.Rebind(rv)                         // originally bound.  Therefore b.Rebind(rv) is valid.
	_ = b.Set("Num", 100)
	_ = b.Set("Str", "reflect.Value!")
	r := rv.Elem().Interface().(S)

	fmt.Println(s.Str, s.Num)
	fmt.Println(r.Str, r.Num)

	// Output: Hello 42
	// reflect.Value! 100
}
