package examples_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
	_ "github.com/nofeaturesonlybugs/set/examples"
)

func ExampleScalar_Set_slice() {
	// This example demonstrates calling Set() with a slice as its argument when the Value is a scalar:
	//	set.V(&Scalar).To(Slice)
	//
	// The Value is set to its proper zero value.  Then if the slice is non-nil and non-zero length the scalar
	// is set to the last element in the slice.
	{
		b1, b2 := true, true
		v1, v2 := set.V(&b1), set.V(&b2)
		v1.To([]interface{}{"Hello", 1, "False"})
		v2.To(([]interface{})(nil))
		fmt.Println(b1, b2)
	}
	{
		i1, i2 := int(42), int(42)
		v1, v2 := set.V(&i1), set.V(&i2)
		v1.To([]interface{}{"Hello", 1, "False"})
		v2.To(([]interface{})(nil))
		fmt.Println(i1, i2)
		v1.To([]interface{}{"Hello", 1, "3.14"})
		v2.To([]interface{}{})
		fmt.Println(i1, i2)
	}
	// Output: false false
	// 0 0
	// 3 0
}
