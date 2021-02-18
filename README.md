[![Documentation](https://godoc.org/github.com/nofeaturesonlybugs/set?status.svg)](http://godoc.org/github.com/nofeaturesonlybugs/set)
[![Go Report Card](https://goreportcard.com/badge/github.com/nofeaturesonlybugs/set)](https://goreportcard.com/report/github.com/nofeaturesonlybugs/set)
[![Build Status](https://travis-ci.com/nofeaturesonlybugs/set.svg?branch=master)](https://travis-ci.com/nofeaturesonlybugs/set)
[![codecov](https://codecov.io/gh/nofeaturesonlybugs/set/branch/master/graph/badge.svg)](https://codecov.io/gh/nofeaturesonlybugs/set)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Package `set` is a small wrapper around the official reflect package that facilitates loose type conversion, assignment into native Go types, and utilities to populate deeply nested Go structs.

Read the `godoc` for more detailed explanations and examples but here are some enticing snippets.

## Scalars and Type-Coercion
```go
{
    // type coercion
    b, i := true, 42
    set.V(&b).To("False")    // Sets b to false
    set.V(&i).To("3.14")     // Sets i to 3
}

{
    // type coercion
    a := int(0)
    b := uint(42)
    set.V(&a).To(b)             // This coerces b into a if possible.
    set.V(&a).To("-57")         // Also works.
    set.V(&a).To("Hello")       // Returns an error.    
}
```

## Pointer Allocation Plus Type-Coercion
```go
{
    // pointer allocation and type coercion
    var bppp ***bool
    set.V(&bppp).To("True")
    fmt.Println(***bppp) // Prints true
}
```

## Scalars-to-Slices and Slices-to-Scalars
```go
{
    // assign scalars to slices
    var b []bool
    set.V(&b).To("True") // b is []bool{ true }

    // or slices to scalars (last element wins)
    var b bool
    set.V(&b).To([]bool{ false, false, true } ) // b is true, coercion not needed.
    set.V(&b).To([]interface{}{ float32(1), uint(0) }) // b is false, coercion needed.
}
```

## Slices to Slices Including Type-Coercion
```go
{
    // slices to slices with or without type coercion; new slice is always created!
    var t []bool
    var s []interface{}
    s = []interface{}{ "true", 0, float64(1) }
    set.V(&t).To(s) // t is []bool{ true, false, true }
}

{
    var t []bool
    var s []bool
    s = []bool{ true, false, true }
    set.V(&t).To(s) // t is []bool{ true, false, true } and t != s
}
```

## Filling Structs by Field Name
```go
m := map[string]interface{}{
    "Name": "Bob",
    "Age":  42,
    "Address": map[interface{}]string{
        "Street1": "97531 Some Street",
        "Street2": "",
        "City":    "Big City",
        "State":   "ST",
        "Zip":     "12345",
    },
}
myGetter := set.MapGetter(m)

type Address struct {
    Street1 string 
    Street2 string 
    City    string 
    State   string 
    Zip     string 
}
type Person struct {
    Name    string  
    Age     uint    
    Address Address 
}
var t Person
set.V(&t).Fill(myGetter)
```

## Filling Structs by Struct Tag
```go
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
```

## Allocating Struct Pointers and Pointer Fields
```go
type Address struct {
    Street1 string `key:"street1"`
    Street2 string `key:"street2"`
    City    string `key:"city"`
    State   string `key:"state"`
    Zip     string `key:"zip"`
}
type Person struct {
    Name    string   `key:"name"`
    Age     uint     `key:"age"`
    Address *Address `key:"address"`
}
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
getter := set.MapGetter(m)

var t *Person
set.V(&t).FillByTag("key", getter)
fmt.Println(t.Name)                 // Bob
fmt.Println(t.Address.Street1)      // 97531 Some Street
```

## Field Pointers Always Allocated
```go
type Address struct {
    Street1 string `key:"street1"`
    Street2 string `key:"street2"`
    City    string `key:"city"`
    State   string `key:"state"`
    Zip     string `key:"zip"`
}
type Person struct {
    Name    string   `key:"name"`
    Age     uint     `key:"age"`
    Address *Address `key:"address"`
}
m := map[string]interface{}{
    "name": "Bob",
    "age":  42,
    // address is missing!
}
getter := set.MapGetter(m)
var t *Person
set.V(&t).FillByTag("key", getter)
fmt.Printf("%p\n", t.Address) // Prints a memory address; the field was allocated.
```

## Pointer-to-Slices-of-Struct-Pointers -- OH MY!

Also noteworthy in this example is the same fuzzy logic for assigning `scalar-to-slice` or `slice-to-scalar` also works for `struct-to-[]struct` and `[]struct-to-struct`.
```go
func TestValue_fillNestedStructPointerToSlicesAsPointers(t *testing.T) {
    chk := assert.New(t)
    //
    var err error
    type Address struct {
        Street1 string `key:"street1"`
        Street2 string `key:"street2"`
        City    string `key:"city"`
        State   string `key:"state"`
        Zip     string `key:"zip"`
    }
    type Person struct {
        Name    string   `key:"name"`
        Age     uint     `key:"age"`
        Address *Address `key:"address"`
    }
    type Company struct {
        Name         string     `key:"name"`
        // Within the Getter "employees" is a []map so it is intuitive that Company.Employees becomes
        // a slice with as many entries as Getter( "employees" ).
        Employees    *[]*Person `key:"employees"`
        // The "employees" key is reused here but now it is going into a single struct pointer; the
        // last set of data in Getter( "employees" ) wins, aka Sally.
        LastEmployee *Person    `key:"employees"`
        // Note the "slice" key in the getter is not a []map itself; but a single *Person is
        // created and inserted into Company.Slice.
        Slice        *[]*Person `key:"slice"`
    }
    //
    // Also noteworthy is that Company.Employees and Company.Slice are pointers to slices and
    // they are instantiated automatically by this package.  The syntax to use them becomes unwieldly
    // so I don't know why you'd want to do this but hey -- it works and that's neat.
    //
    m := map[string]interface{}{
        "name": "Some Company",
        "slice": map[string]interface{}{
            "name": "Slice",
            "age":  2,
            "address": map[interface{}]string{
                "street1": "Slice Street",
                "street2": "",
                "city":    "Slice City",
                "state":   "SL",
                "zip":     "99999",
            },
        },
        "employees": []map[string]interface{}{
            {
                "name": "Bob",
                "age":  42,
                "address": map[interface{}]string{
                    "street1": "97531 Some Street",
                    "street2": "",
                    "city":    "Big City",
                    "state":   "ST",
                    "zip":     "12345",
                },
            },
            {
                "name": "Sally",
                "age":  48,
                "address": map[interface{}]string{
                    "street1": "555 Small Lane",
                    "street2": "",
                    "city":    "Other City",
                    "state":   "OO",
                    "zip":     "54321",
                },
            },
        },
    }
    getter := set.MapGetter(m)

    var t *Company
    err = set.V(&t).FillByTag("key", getter)
    chk.NoError(err)
    chk.Equal("Some Company", t.Name)
    //
    chk.Equal(2, len(*t.Employees))
    //
    chk.Equal("Bob", (*t.Employees)[0].Name)
    chk.Equal(uint(42), (*t.Employees)[0].Age)
    chk.Equal("97531 Some Street", (*t.Employees)[0].Address.Street1)
    chk.Equal("", (*t.Employees)[0].Address.Street2)
    chk.Equal("Big City", (*t.Employees)[0].Address.City)
    chk.Equal("ST", (*t.Employees)[0].Address.State)
    chk.Equal("12345", (*t.Employees)[0].Address.Zip)
    //
    chk.Equal("Sally", (*t.Employees)[1].Name)
    chk.Equal(uint(48), (*t.Employees)[1].Age)
    chk.Equal("555 Small Lane", (*t.Employees)[1].Address.Street1)
    chk.Equal("", (*t.Employees)[1].Address.Street2)
    chk.Equal("Other City", (*t.Employees)[1].Address.City)
    chk.Equal("OO", (*t.Employees)[1].Address.State)
    chk.Equal("54321", (*t.Employees)[1].Address.Zip)
    //
    chk.Equal("Sally", t.LastEmployee.Name)
    chk.Equal(uint(48), t.LastEmployee.Age)
    chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
    chk.Equal("", t.LastEmployee.Address.Street2)
    chk.Equal("Other City", t.LastEmployee.Address.City)
    chk.Equal("OO", t.LastEmployee.Address.State)
    chk.Equal("54321", t.LastEmployee.Address.Zip)
    //
    chk.Equal(1, len(*t.Slice))
    chk.Equal("Slice", (*t.Slice)[0].Name)
    chk.Equal("Slice Street", (*t.Slice)[0].Address.Street1)
    chk.Equal("", (*t.Slice)[0].Address.Street2)
    chk.Equal("Slice City", (*t.Slice)[0].Address.City)
    chk.Equal("SL", (*t.Slice)[0].Address.State)
    chk.Equal("99999", (*t.Slice)[0].Address.Zip)
}
```

## Mapper and Mapping for Nested Struct Access
`Mapper` traverses a nested struct hierarchy to generate a `Mapping`.  From a `Mapping` you can use `string` keys to
access the struct members by `struct field index`.
```go
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
```
Depending on the `Mapper` options the following mappings can be created:
```go
// Type CommonDb doesn't affect names; join nestings with "_"
mapper := &set.Mapper{
    Elevated: set.NewTypeList(CommonDb{}),
    Join:     "_",
}
generates = `
[0 0] Pk
[0 1] CreatedTime
[0 2] UpdatedTime
[1] Name
[2] Age
`

// lowercase with dot separators
mapper := &set.Mapper{
    Join:      ".",
    Transform: strings.ToLower,
}
generates = `
[0 0] commondb.pk
[0 1] commondb.createdtime
[0 2] commondb.updatedtime
[1] name
[2] age
`

// specify tags
mapper := &set.Mapper{
    Join: "_",
    Tags: []string{"t"},
}
generates = `
[0 0] common_pk
[0 1] common_created_time
[0 2] common_updated_time
[1] name
[2] age
`
```
