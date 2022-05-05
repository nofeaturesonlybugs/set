package set_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set"
)

func ExamplePreparedMapping_Rebind() {
	// Once a PreparedMapping has been created you can swap out the value it mutates
	// by calling Rebind.  This yields better performance in tight loops where
	// you intend to mutate many instances of the same type.

	m := &set.Mapper{}

	type S struct {
		Str string
		Num int
	}

	values := [][]interface{}{
		{"First", 1},
		{"Second", 2},
		{"Third", 3},
	}
	slice := make([]S, 3)

	p, _ := m.Prepare(&slice[0]) // We can prepare on the first value to get a PreparedMapping

	// We must call Plan with our intended field access order.
	_ = p.Plan("Str", "Num")

	for k := range slice {
		p.Rebind(&slice[k]) // The PreparedMapping now affects slice[k]
		for _, value := range values[k] {
			p.Set(value)
		}
	}

	fmt.Println(slice[0].Num, slice[0].Str)
	fmt.Println(slice[1].Num, slice[1].Str)
	fmt.Println(slice[2].Num, slice[2].Str)

	// Output: 1 First
	// 2 Second
	// 3 Third
}

func ExamplePreparedMapping_Rebind_panic() {
	// PreparedMapping.Rebind panics if the new instance is not the same type.

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

	p, _ := m.Prepare(&s)
	p.Rebind(&d)

	// Output: mismatching types during Rebind; have *set_test.S and got *set_test.Different
}

func ExamplePreparedMapping_Rebind_reflectValue() {
	// As a convenience PreparedMapping.Rebind can be called with an instance of reflect.Value
	// if-and-only-if the reflect.Value is holding a type compatible with the PreparedMapping.

	m := &set.Mapper{}

	type S struct {
		Str string
		Num int
	}
	var s S

	// Errors ignored for brevity.
	p, _ := m.Prepare(&s)
	p.Plan("Num", "Str")

	p.Set(42)
	p.Set("Hello")

	rv := reflect.New(reflect.TypeOf(s)) // reflect.New creates a *S which is the type
	p.Rebind(rv)                         // originally bound.  Therefore b.Rebind(rv) is valid.
	p.Set(100)
	p.Set("reflect.Value!")
	r := rv.Elem().Interface().(S)

	fmt.Println(s.Str, s.Num)
	fmt.Println(r.Str, r.Num)

	// Output: Hello 42
	// reflect.Value! 100
}
