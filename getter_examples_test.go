package set_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleGetterFunc() {
	// A GetterFunc allows a function to be used as the producer for
	// populating a struct.

	// Note in this example the key names match the field names of the struct
	// and Value.Fill is used to populate the struct.

	// F has the correct signature to become a GetterFunc.
	F := func(name string) interface{} {
		switch name {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		default:
			return nil
		}
	}
	myGetter := set.GetterFunc(F)

	type T struct {
		Name string
		Age  uint
	}
	var t T

	set.V(&t).Fill(myGetter)
	fmt.Println(t.Name, t.Age)

	// Output: Bob 42
}

func ExampleGetterFunc_fillByTag() {
	// In this example the key names are lower case and match values
	// in the struct tags.  To populate a struct by tag use Value.FillByTag.

	// Also note that `address` returns a new GetterFunc used to populate
	// the nested Address struct in Person.

	F := func(key string) interface{} {
		switch key {
		case "name":
			return "Bob"
		case "age":
			return "42"
		case "address":
			return set.GetterFunc(func(key string) interface{} {
				switch key {
				case "street1":
					return "97531 Some Street"
				case "street2":
					return ""
				case "city":
					return "Big City"
				case "state":
					return "ST"
				case "zip":
					return "12345"
				default:
					return nil
				}
			})
		default:
			return nil
		}
	}
	myGetter := set.GetterFunc(F)

	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	var t Person

	set.V(&t).FillByTag("key", myGetter)
	fmt.Println(t.Name, t.Age)
	fmt.Printf("%v, %v, %v  %v\n", t.Address.Street1, t.Address.City, t.Address.State, t.Address.Zip)

	// Output: Bob 42
	// 97531 Some Street, Big City, ST  12345
}

func ExampleMapGetter() {
	// A map can also be used as a Getter by casting to MapGetter.
	//
	// Note that the map key must be either string or interface{}.
	//
	// If the map contains nested maps that can be used as MapGetter
	// then those maps can populate nested structs in the hierarchy.

	m := map[string]interface{}{
		"name": "Bob",
		"age":  42,
		"address": map[interface{}]string{
			"street1": "97531 Some Street",
			"street2": "",
			"city":    "Big City",
			"state":   "ST",
			"zip":     "12345",
		},
	}
	myGetter := set.MapGetter(m)

	type Address struct {
		Street1 string `key:"street1"`
		Street2 string `key:"street2"`
		City    string `key:"city"`
		State   string `key:"state"`
		Zip     string `key:"zip"`
	}
	type Person struct {
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	var t Person

	set.V(&t).FillByTag("key", myGetter)
	fmt.Println(t.Name, t.Age)
	fmt.Printf("%v, %v, %v  %v\n", t.Address.Street1, t.Address.City, t.Address.State, t.Address.Zip)

	// Output: Bob 42
	// 97531 Some Street, Big City, ST  12345
}
