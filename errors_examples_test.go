package set_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/set"
)

func Example_valueErrors() {
	// This example demonstrates some of the errors returned from the Value type.

	type S struct {
		A struct {
			AA int
			BB int
		}
		B int
	}
	var s S
	var slice []int
	var n int
	var err error

	// set: Value.Append: unsupported: nil value: hint=[set.V(nil) was called]
	err = set.V(nil).Append(42, 24)
	fmt.Println(err)

	// set: Value.Append: read only value: []int is not writable: hint=[call to set.V([]int) should have been set.V(*[]int)]
	err = set.V(slice).Append(42, 24)
	fmt.Println(err)

	// set: Value.Append: unsupported: can not append to int
	err = set.V(&n).Append(42)
	fmt.Println(err)

	// set: Value.FieldByIndex: unsupported: nil value: hint=[set.V(nil) was called]
	_, err = set.V(nil).FieldByIndex([]int{0})
	fmt.Println(err)

	// set: Value.FieldByIndex: read only value: set_test.S is not writable: hint=[call to set.V(set_test.S) should have been set.V(*set_test.S)]
	_, err = set.V(s).FieldByIndex([]int{0})
	fmt.Println(err)

	// set: Value.FieldByIndex: unsupported: empty index
	_, err = set.V(&s).FieldByIndex(nil)
	fmt.Println(err)

	// set: Value.FieldByIndex: set: index out of bounds: index 2 exceeds max 1
	_, err = set.V(&s).FieldByIndex([]int{2})
	fmt.Println(err)

	// set: Value.FieldByIndex: unsupported: want struct but got int
	_, err = set.V(&s).FieldByIndex([]int{1, 0})
	fmt.Println(err)

	// TODO Fill+FillByTag

	// set: Value.Zero: unsupported: nil value: hint=[set.V(nil) was called]
	err = set.V(nil).Zero()
	fmt.Println(err)

	// set: Value.Zero: read only value: int is not writable: hint=[call to set.V(int) should have been set.V(*int)]
	err = set.V(n).Zero()
	fmt.Println(err)

	// set: Value.To: unsupported: nil value: hint=[set.V(nil) was called]
	err = set.V(nil).To("Hello!")
	fmt.Println(err)

	err = set.V(n).To(42)
	// set: Value.To: read only value: int is not writable: hint=[call to set.V(int) should have been set.V(*int)]
	fmt.Println(err)

	// Output: set: Value.Append: unsupported: nil value: hint=[set.V(nil) was called]
	// set: Value.Append: read only value: []int is not writable: hint=[call to set.V([]int) should have been set.V(*[]int)]
	// set: Value.Append: unsupported: can not append to int
	// set: Value.FieldByIndex: unsupported: nil value: hint=[set.V(nil) was called]
	// set: Value.FieldByIndex: read only value: set_test.S is not writable: hint=[call to set.V(set_test.S) should have been set.V(*set_test.S)]
	// set: Value.FieldByIndex: unsupported: empty index
	// set: Value.FieldByIndex: index out of bounds: index 2 exceeds max 1
	// set: Value.FieldByIndex: unsupported: want struct but got int
	// set: Value.Zero: unsupported: nil value: hint=[set.V(nil) was called]
	// set: Value.Zero: read only value: int is not writable: hint=[call to set.V(int) should have been set.V(*int)]
	// set: Value.To: unsupported: nil value: hint=[set.V(nil) was called]
	// set: Value.To: read only value: int is not writable: hint=[call to set.V(int) should have been set.V(*int)]
}

func Example_mapperErrors() {
	// This example demonstrates some of the errors returned from the Mapper type.

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Company struct {
		Name  string `json:"name"`
		Owner Person `json:"owner"`
	}

	mapper := set.Mapper{
		Tags: []string{"json"},
		Join: "_",
	}

	var company Company
	var err error

	// set: Mapper.Bind: read only value: set_test.Company is not writable: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	_, err = mapper.Bind(company)
	fmt.Println(err)

	// set: Mapper.Prepare: read only value: set_test.Company is not writable: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	_, err = mapper.Prepare(company)
	fmt.Println(err)

	// Output: set: Mapper.Bind: read only value: set_test.Company is not writable: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	// set: Mapper.Prepare: read only value: set_test.Company is not writable: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
}

func Example_boundMappingErrors() {
	// This example demonstrates some of the errors returned from the BoundMapping type.

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Company struct {
		Name  string `json:"name"`
		Owner Person `json:"owner"`
	}

	mapper := set.Mapper{
		Tags: []string{"json"},
		Join: "_",
	}

	var company Company
	var err error

	readonly, _ := mapper.Bind(company)
	b, _ := mapper.Bind(&company)

	// set: BoundMapping.Assignables: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	_, err = readonly.Assignables([]string{"name"}, nil)
	fmt.Println(err)

	// set: BoundMapping.Field: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	_, err = readonly.Field("name")
	fmt.Println(err)

	// set: BoundMapping.Fields: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	_, err = readonly.Fields([]string{"name"}, nil)
	fmt.Println(err)

	// set: BoundMapping.Set: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	err = readonly.Set("foobar", "ABC Widgets")
	fmt.Println(err)

	// set: BoundMapping.Assignables: unknown field: field [foobar] not found in type *set_test.Company
	_, err = b.Assignables([]string{"foobar"}, nil)
	fmt.Println(err)

	// set: BoundMapping.Field: unknown field: field [foobar] not found in type *set_test.Company
	_, err = b.Field("foobar")
	fmt.Println(err)

	// set: BoundMapping.Fields: unknown field: field [foobar] not found in type *set_test.Company
	_, err = b.Fields([]string{"foobar"}, nil)
	fmt.Println(err)

	// set: BoundMapping.Set: unknown field: field [foobar] not found in type *set_test.Company
	err = b.Set("foobar", "ABC Widgets")
	fmt.Println(err)

	// Output: set: BoundMapping.Assignables: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	// set: BoundMapping.Field: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	// set: BoundMapping.Fields: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	// set: BoundMapping.Set: read only value: hint=[call to Mapper.Bind(set_test.Company) should have been Mapper.Bind(*set_test.Company)]
	// set: BoundMapping.Assignables: unknown field: field [foobar] not found in type *set_test.Company
	// set: BoundMapping.Field: unknown field: field [foobar] not found in type *set_test.Company
	// set: BoundMapping.Fields: unknown field: field [foobar] not found in type *set_test.Company
	// set: BoundMapping.Set: unknown field: field [foobar] not found in type *set_test.Company
}

func Example_preparedMappingErrors() {
	// This example demonstrates some of the errors returned from the PreparedMapping type.

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Company struct {
		Name  string `json:"name"`
		Owner Person `json:"owner"`
	}

	mapper := set.Mapper{
		Tags: []string{"json"},
		Join: "_",
	}

	var company Company
	var err error

	readonly, _ := mapper.Prepare(company)
	p, _ := mapper.Prepare(&company)

	// set: PreparedMapping.Assignables: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	_, err = readonly.Assignables(nil)
	fmt.Println(err)

	// set: PreparedMapping.Field: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	_, err = readonly.Field()
	fmt.Println(err)

	// set: PreparedMapping.Fields: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	_, err = readonly.Fields(nil)
	fmt.Println(err)

	// set: PreparedMapping.Plan: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	err = readonly.Plan("name")
	fmt.Println(err)

	// set: PreparedMapping.Field: no plan: hint=[call PreparedMapping.Plan to prepare access plan for *set_test.Company]
	_, err = p.Field()
	fmt.Println(err)

	// set: PreparedMapping.Set: no plan: hint=[call PreparedMapping.Plan to prepare access plan for *set_test.Company]
	err = p.Set("ABC Widgets")
	fmt.Println(err)

	// set: PreparedMapping.Plan: unknown field: field [foobar] not found in type *set_test.Company
	err = p.Plan("foobar")
	fmt.Println(err)

	_ = p.Plan("name", "owner_name") // No error
	_, _ = p.Field()                 // No error
	_, _ = p.Field()                 // No error
	// set: PreparedMapping.Field: attempted access extends plan: value of *set_test.Company
	_, err = p.Field()
	fmt.Println(err)

	_ = p.Plan("name", "owner_name") // No error
	_ = p.Set("ABC Widgets")         // No error
	_ = p.Set("Larry")               // No error
	// set: PreparedMapping.Set: attempted access extends plan: value of *set_test.Company
	err = p.Set("extended the plan")
	fmt.Println(err)

	// Output: set: PreparedMapping.Assignables: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	// set: PreparedMapping.Field: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	// set: PreparedMapping.Fields: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	// set: PreparedMapping.Plan: read only value: hint=[call to Mapper.Prepare(set_test.Company) should have been Mapper.Prepare(*set_test.Company)]
	// set: PreparedMapping.Field: no plan: hint=[call PreparedMapping.Plan to prepare access plan for *set_test.Company]
	// set: PreparedMapping.Set: no plan: hint=[call PreparedMapping.Plan to prepare access plan for *set_test.Company]
	// set: PreparedMapping.Plan: unknown field: field [foobar] not found in type *set_test.Company
	// set: PreparedMapping.Field: attempted access extends plan: value of *set_test.Company
	// set: PreparedMapping.Set: attempted access extends plan: value of *set_test.Company
}
