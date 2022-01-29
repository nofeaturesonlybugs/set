package set_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

func BenchmarkValue(b *testing.B) {
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		*Common
		*Timestamps // Not used but present anyways
		First       string
		Last        string
	}
	type Vendor struct {
		*Common
		*Timestamps // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type T struct {
		*Common
		*Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	//
	for k := 0; k < b.N; k++ {
		dest := new(T)
		set.V(dest)
	}
}
