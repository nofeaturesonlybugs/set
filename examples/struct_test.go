package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleStruct_Fill() {
	type T struct {
		Bool   bool
		Int    int
		Uint   uint
		String string
	}
	values := map[string]interface{}{
		"Bool":   true,
		"Int":    -42,
		"Uint":   42,
		"String": "Hello, World!",
	}
	fn := func(key string) interface{} {
		if value, ok := values[key]; ok {
			return value
		}
		return nil
	}

	var t T
	if err := set.V(&t).Fill(set.GetterFunc(fn)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t.Bool, t.Int, t.Uint, t.String)
	}
	// Output: true -42 42 Hello, World!
}

func ExampleStruct_Fill_byTag() {
	type T struct {
		Bool   bool   `someTag:"myBool"`
		Int    int    `someTag:"myInt"`
		Uint   uint   `someTag:"myUint"`
		String string `someTag:"myString"`
	}
	values := map[string]interface{}{
		"myBool":   true,
		"myInt":    -42,
		"myUint":   42,
		"myString": "Hello, World!",
	}
	fn := func(key string) interface{} {
		if value, ok := values[key]; ok {
			return value
		}
		return nil
	}

	var t T
	if err := set.V(&t).FillByTag("someTag", set.GetterFunc(fn)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t.Bool, t.Int, t.Uint, t.String)
	}
	// Output: true -42 42 Hello, World!
}

func ExampleStruct_Fill_byTagWithPunctuation() {
	type T struct {
		Bool   bool   `some.tag:"my-bool"`
		Int    int    `some.tag:"my-int"`
		Uint   uint   `some.tag:"my-uint"`
		String string `some.tag:"my-string"`
	}
	values := map[string]interface{}{
		"my-bool":   true,
		"my-int":    -42,
		"my-uint":   42,
		"my-string": "Hello, World!",
	}
	fn := func(key string) interface{} {
		if value, ok := values[key]; ok {
			return value
		}
		return nil
	}

	var t T
	if err := set.V(&t).FillByTag("some.tag", set.GetterFunc(fn)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(t.Bool, t.Int, t.Uint, t.String)
	}
	// Output: true -42 42 Hello, World!
}
