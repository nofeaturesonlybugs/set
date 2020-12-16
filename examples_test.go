package set_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleValue_To_addressability() {
	fmt.Println("When using Value.To the target Value must be addressable.")
	//
	{
		var value bool
		if err := set.V(value).To(true); err != nil {
			// Expected
			fmt.Println("1. Error because address-of value was not passed.")
		} else {
			fmt.Printf("1. Value is %v\n", value)
		}
	}
	{
		var value bool
		if err := set.V(&value).To(true); err != nil {
			fmt.Println("2. Error because address-of value was not passed.")
		} else {
			// Expected
			fmt.Printf("2. Value is %v\n", value)
		}
	}
	{
		var value *int
		if err := set.V(value).To(42); err != nil {
			// Expected
			fmt.Println("3. Even though value is a pointer itself its address still needs to be passed.")
			if value == nil {
				// Expected
				fmt.Println("3. Also worth noting the pointer remains nil.")
			}
		} else {
			fmt.Printf("3. Value is %v\n", value)
		}
	}
	{
		var value *int
		if err := set.V(&value).To(42); err != nil {
			fmt.Println("4. Even though value is a pointer itself its address still needs to be passed.")
			if value == nil {
				fmt.Println("4. Also worth noting the pointer remains nil.")
			}
		} else {
			// Expected
			fmt.Printf("4. Value is %v\n", *value)
			if value != nil {
				// Expected
				fmt.Println("4. A pointer-to-int was created!")
			}
		}
	}
	{
		var value ***int
		if err := set.V(&value).To(42); err != nil {
			fmt.Println("5. Even though value is a pointer itself its address still needs to be passed.")
			if value == nil {
				fmt.Println("5. Also worth noting the pointer remains nil.")
			}
		} else {
			// Expected
			fmt.Printf("5. Value is %v\n", ***value)
			fmt.Println("5. Multiple pointers had to be created for this to work!")
		}
	}
	{
		value := new(int)
		if err := set.V(value).To(42); err != nil {
			fmt.Println("6. Even though value is a pointer itself its address still needs to be passed.")
		} else {
			// Expected
			fmt.Printf("6. Value is %v\n", *value)
		}
	}
	{
		value := new(**int)
		if err := set.V(value).To(42); err != nil {
			fmt.Println("7. Even though value is a pointer itself its address still needs to be passed.")
		} else {
			// Expected
			fmt.Printf("7. Value is %v\n", ***value)
		}
	}
	{
		var value string
		if err := set.V(&value).To("8. It works with strings too."); err == nil {
			fmt.Println(value)
		}
	}
	{
		var value ***string
		if err := set.V(&value).To("9. String pointers are no different."); err == nil {
			fmt.Println(***value)
		}
	}
	// Output: When using Value.To the target Value must be addressable.
	// 1. Error because address-of value was not passed.
	// 2. Value is true
	// 3. Even though value is a pointer itself its address still needs to be passed.
	// 3. Also worth noting the pointer remains nil.
	// 4. Value is 42
	// 4. A pointer-to-int was created!
	// 5. Value is 42
	// 5. Multiple pointers had to be created for this to work!
	// 6. Value is 42
	// 7. Value is 42
	// 8. It works with strings too.
	// 9. String pointers are no different.
}
