package set_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nofeaturesonlybugs/set"
)

func ExampleMapper() {
	// This example maps the same type (Person) with three different Mapper instances.
	// Each Mapper instance has slightly different configuration to control the mapped
	// names that are generated.

	type CommonDb struct {
		Pk          int    `t:"pk"`
		CreatedTime string `t:"created_time"`
		UpdatedTime string `t:"updated_time"`
	}
	type Person struct {
		CommonDb `t:"common"`
		Name     string `t:"name"`
		Age      int    `t:"age"`
	}
	var data Person
	{
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(CommonDb{}),
			Join:     "_",
		}
		mapping := mapper.Map(&data)
		fmt.Println(strings.ReplaceAll(mapping.String(), "\t\t", " "))
	}
	{
		fmt.Println("")
		fmt.Println("lowercase with dot separators")
		mapper := &set.Mapper{
			Join:      ".",
			Transform: strings.ToLower,
		}
		mapping := mapper.Map(&data)
		fmt.Println(strings.ReplaceAll(mapping.String(), "\t\t", " "))
	}
	{
		fmt.Println("")
		fmt.Println("specify tags")
		mapper := &set.Mapper{
			Join: "_",
			Tags: []string{"t"},
		}
		mapping := mapper.Map(&data)
		fmt.Println(strings.ReplaceAll(mapping.String(), "\t\t", " "))
	}

	// Output: [0 0] Pk
	// [0 1] CreatedTime
	// [0 2] UpdatedTime
	// [1] Name
	// [2] Age
	//
	// lowercase with dot separators
	// [0 0] commondb.pk
	// [0 1] commondb.createdtime
	// [0 2] commondb.updatedtime
	// [1] name
	// [2] age
	//
	// specify tags
	// [0 0] common_pk
	// [0 1] common_created_time
	// [0 2] common_updated_time
	// [1] name
	// [2] age
}
func ExampleMapper_treatAsScalar() {
	// By default Mapper maps **all** exported fields in nested or embedded structs.

	type S struct {
		N sql.NullString // sql.NullString has public fields that can be mapped.
	}
	var s S

	// Notice that **this** Mapper maps the public fields in the sql.NullString type.
	fmt.Println("Without TreatAsScalar")
	all := &set.Mapper{
		Join: ".",
	}
	m := all.Map(&s)
	fmt.Println(strings.ReplaceAll(m.String(), "\t\t", " "))

	// While **this** Mapper treats N as a scalar field and does not map its public fields.
	fmt.Println("TreatAsScalar")
	scalar := &set.Mapper{
		TreatAsScalar: set.NewTypeList(sql.NullString{}), // Fields of sql.NullString are now treated as scalars.
		Join:          ".",
	}
	m = scalar.Map(&s)
	fmt.Println(strings.ReplaceAll(m.String(), "\t\t", " "))

	// Output: Without TreatAsScalar
	// [0 0] N.String
	// [0 1] N.Valid
	// TreatAsScalar
	// [0] N
}

func ExampleMapper_treatAsScalarTime() {
	// As a convenience time.Time and *time.Time are always treated as scalars.

	type S struct {
		T    time.Time
		TPtr *time.Time
	}
	var s S

	// Even though this mapper does not configure TreatAsScalar the time.Time fields **are**
	// treated as scalars.
	all := &set.Mapper{
		Join: ".",
	}
	m := all.Map(&s)
	fmt.Println(strings.ReplaceAll(m.String(), "\t\t", " "))

	// Output: [0] T
	// [1] TPtr
}

func ExampleMapper_Bind() {
	// This example demonstrates a simple flat struct of primitives.
	type S struct {
		Num int
		Str string
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{}
	b, _ := m.Bind(&s) // err ignored for brevity

	// b.Set errors ignored for brevity
	_ = b.Set("Str", 3.14) // 3.14 coerced to "3.14"
	_ = b.Set("Num", "42") // "42" coerced to 42

	fmt.Println(s.Num, s.Str)
	// Output: 42 3.14
}

func ExampleMapper_Bind_nesting() {
	// This example demonstrates how nested structs are accessible via their mapped names.
	type Nest struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Foo Nest
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Join: ".", // Nested and embedded fields join with DOT
	}
	b, _ := m.Bind(&s) // err ignored for brevity

	// b.Set errors ignored for brevity
	_ = b.Set("Str", 3.14)   // 3.14 coerced to "3.14"
	_ = b.Set("Num", "42")   // "42" coerced to 42
	_ = b.Set("Foo.V", 3.14) // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.Foo.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Bind_embedded() {
	// This example demonstrates how to access fields in an embedded struct.
	// Notice in this example the field is named Embed.V -- the struct type name
	// becomes part of the mapping name.
	type Embed struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Embed
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Join: ".", // Nested and embedded fields join with DOT
	}
	b, _ := m.Bind(&s) // err ignored for brevity

	// b.Set errors ignored for brevity
	_ = b.Set("Str", 3.14)     // 3.14 coerced to "3.14"
	_ = b.Set("Num", "42")     // "42" coerced to 42
	_ = b.Set("Embed.V", 3.14) // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Bind_elevatedEmbed() {
	// This example demonstrates how to access fields in an elevated embedded struct.
	// Elevated types can be embedded without their type name becoming part of the field name.
	type Embed struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Embed
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Elevated: set.NewTypeList(Embed{}), // Elevate embedded fields of type Embed
		Join:     ".",                      // Nested and embedded fields join with DOT
	}
	b, _ := m.Bind(&s) // err ignored for brevity

	// b.Set errors ignored for brevity
	_ = b.Set("Str", 3.14) // 3.14 coerced to "3.14"
	_ = b.Set("Num", "42") // "42" coerced to 42
	_ = b.Set("V", 3.14)   // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Bind_reflectValue() {
	// As a convenience Mapper.Bind will accept a reflect.Value to perform the binding.
	type S struct {
		Num int
		Str string
	}
	var s, t, u S

	// Create a mapper.
	m := &set.Mapper{}
	b, _ := m.Bind(reflect.ValueOf(&s)) // err ignored for brevity

	// b.Set errors ignored for brevity
	_ = b.Set("Str", 3.14)
	_ = b.Set("Num", "42")

	b.Rebind(reflect.ValueOf(&t))
	_ = b.Set("Str", -3.14)
	_ = b.Set("Num", "24")

	// Even though the BoundMapping was created with reflect.Value it will still accept *S directly.
	b.Rebind(&u)
	_ = b.Set("Str", "Works!")
	_ = b.Set("Num", uint(100))

	fmt.Println("s", s.Num, s.Str)
	fmt.Println("t", t.Num, t.Str)
	fmt.Println("u", u.Num, u.Str)
	// Output: s 42 3.14
	// t 24 -3.14
	// u 100 Works!
}

func ExampleMapper_Prepare() {
	// This example demonstrates a simple flat struct of primitives.
	type S struct {
		Num int
		Str string
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{}
	p, _ := m.Prepare(&s) // err ignored for brevity

	// PreparedBindings require a call to Plan indicating field order.
	_ = p.Plan("Str", "Num") // err ignored for brevity

	_ = p.Set(3.14) // 3.14 coerced to "3.14"
	_ = p.Set("42") // "42" coerced to 42

	fmt.Println(s.Num, s.Str)
	// Output: 42 3.14
}

func ExampleMapper_Prepare_nesting() {
	// This example demonstrates how nested structs are accessible via their mapped names.
	type Nest struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Foo Nest
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Join: ".", // Nested and embedded fields join with DOT
	}
	p, _ := m.Prepare(&s) // err ignored for brevity

	// PreparedBindings require a call to Plan indicating field order.
	_ = p.Plan("Str", "Num", "Foo.V") // err ignored for brevity

	_ = p.Set(3.14) // 3.14 coerced to "3.14"
	_ = p.Set("42") // "42" coerced to 42
	_ = p.Set(3.14) // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.Foo.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Prepare_embedded() {
	// This example demonstrates how to access fields in an embedded struct.
	// Notice in this example the field is named Embed.V -- the struct type name
	// becomes part of the mapping name.
	type Embed struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Embed
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Join: ".", // Nested and embedded fields join with DOT
	}
	p, _ := m.Prepare(&s) // err ignored for brevity

	// PreparedBindings require a call to Plan indicating field order.
	_ = p.Plan("Str", "Num", "Embed.V") // err ignored for brevity

	_ = p.Set(3.14) // 3.14 coerced to "3.14"
	_ = p.Set("42") // "42" coerced to 42
	_ = p.Set(3.14) // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Prepare_elevatedEmbed() {
	// This example demonstrates how to access fields in an elevated embedded struct.
	// Elevated types can be embedded without their type name becoming part of the field name.
	type Embed struct {
		V int
	}
	type S struct {
		Num int
		Str string
		Embed
	}
	var s S

	// Create a mapper.
	m := &set.Mapper{
		Elevated: set.NewTypeList(Embed{}), // Elevate embedded fields of type Embed
		Join:     ".",                      // Nested and embedded fields join with DOT
	}
	p, _ := m.Prepare(&s) // err ignored for brevity

	// PreparedBindings require a call to Plan indicating field order.
	_ = p.Plan("Str", "Num", "V") // err ignored for brevity

	_ = p.Set(3.14) // 3.14 coerced to "3.14"
	_ = p.Set("42") // "42" coerced to 42
	_ = p.Set(3.14) // 3.14 coerced to 3

	fmt.Println(s.Num, s.Str, s.V)
	// Output: 42 3.14 3
}

func ExampleMapper_Prepare_reflectValue() {
	// As a convenience Mapper.Prepare will accept a reflect.Value to perform the binding.
	type S struct {
		Num int
		Str string
	}
	var s, t, u S

	// Create a mapper.
	m := &set.Mapper{}
	p, _ := m.Prepare(reflect.ValueOf(&s)) // err ignored for brevity
	_ = p.Plan("Str", "Num")               // err ignored for brevity

	_ = p.Set(3.14)
	_ = p.Set("42")

	p.Rebind(reflect.ValueOf(&t))
	_ = p.Set(-3.14)
	_ = p.Set("24")

	// Even though the PreparedMapping was created with reflect.Value it will still accept *S directly.
	p.Rebind(&u)
	_ = p.Set("Works!")
	_ = p.Set(uint(100))

	fmt.Println("s", s.Num, s.Str)
	fmt.Println("t", t.Num, t.Str)
	fmt.Println("u", u.Num, u.Str)
	// Output: s 42 3.14
	// t 24 -3.14
	// u 100 Works!
}
