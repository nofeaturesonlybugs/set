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
			b.Set(key, value)
		}
	}

	fmt.Println(slice[0].Num, slice[0].Str)
	fmt.Println(slice[1].Num, slice[1].Str)
	fmt.Println(slice[2].Num, slice[2].Str)

	// Output: 1 First
	// 2 Second
	// 3 Third
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
	b.Set("Num", 42)
	b.Set("Str", "Hello")

	rv := reflect.New(reflect.TypeOf(s)) // reflect.New creates a *S which is the type
	b.Rebind(rv)                         // originally bound.  Therefore b.Rebind(rv) is valid.
	b.Set("Num", 100)
	b.Set("Str", "reflect.Value!")
	r := rv.Elem().Interface().(S)

	fmt.Println(s.Str, s.Num)
	fmt.Println(r.Str, r.Num)

	// Output: Hello 42
	// reflect.Value! 100
}
