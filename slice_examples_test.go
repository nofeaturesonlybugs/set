package set_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleSlice() {
	var v set.Value
	var nums []int
	var numsppp ***[]int

	slice, err := set.Slice(&nums)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, elem := range []interface{}{"42", false, "true", "3.14"} {
		if k == 0 {
			v = set.V(slice.Elem())
		} else {
			v.Rebind(slice.Elem())
		}
		if err = v.To(elem); err != nil {
			fmt.Println(err)
			return
		}
		slice.Append(v.TopValue)
	}

	slice, err = set.Slice(&numsppp)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, elem := range []interface{}{"42", false, "true", "3.14"} {
		if k == 0 {
			v = set.V(slice.Elem())
		} else {
			v.Rebind(slice.Elem())
		}
		if err = v.To(elem); err != nil {
			fmt.Println(err)
			return
		}
		slice.Append(v.TopValue)
	}

	fmt.Println(nums)
	fmt.Println(***numsppp)

	// Output: [42 0 1 3]
	// [42 0 1 3]
}

func ExampleSlice_errors() {
	var n int
	var numspp **[]int

	_, err := set.Slice(n)
	fmt.Println(err)

	_, err = set.Slice(&n)
	fmt.Println(err)

	_, err = set.Slice(reflect.ValueOf(&n))
	fmt.Println(err)

	_, err = set.Slice(numspp)
	fmt.Println(err)

	// Output: set: Slice: invalid slice: expected pointer to slice; got int
	// set: Slice: invalid slice: expected pointer to slice; got *int
	// set: Slice: invalid slice: expected pointer to slice; got *int
	// set: Slice: read only value: can not set **[]int
}
