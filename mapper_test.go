package set_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestMapper(t *testing.T) {
	chk := assert.New(t)
	{
		// Default options test to increase coverage.
		type A struct {
			A string
			B string
		}
		var data A
		mapping := set.DefaultMapper.Map(&data)
		chk.Equal("[0]", fmt.Sprintf("%v", mapping.Get("A")))
		chk.Equal("[1]", fmt.Sprintf("%v", mapping.Get("B")))
	}
	{
		// Skips unexported fields.
		type A struct {
			A     string
			B     string
			shhhh int
		}
		var data A
		mapping := set.DefaultMapper.Map(&data)
		chk.Equal("[0]", fmt.Sprintf("%v", mapping.Get("A")))
		chk.Equal("[1]", fmt.Sprintf("%v", mapping.Get("B")))
		chk.Empty(mapping.Get("shhhh"))
	}
	{
		type CommonDb struct {
			Pk          int    `db:"pk" json:"id"`
			CreatedTime string `db:"created_tmz" json:"created_time"`
			UpdatedTime string `db:"modified_tmz" json:"modified_time"`
		}
		type Person struct {
			CommonDb
			Name string
			Age  int
		}
		type Logger struct {
			Info  func()
			Warn  func()
			Error func()
		}
		type Combined struct {
			Child     Person
			Parent    Person
			Emergency Person `db:"Emergency"`
			Logger
		}
		var data Combined
		mapper := &set.Mapper{
			Ignored:   set.NewTypeList(Logger{}),
			Elevated:  set.NewTypeList(CommonDb{}),
			Tags:      []string{"db", "json"},
			Join:      "_",
			Transform: strings.ToLower,
		}
		mapping := mapper.Map(&data)
		//
		chk.Equal("[0 0 0]", fmt.Sprintf("%v", mapping.Get("child_pk")))
		chk.Equal("[0 0 1]", fmt.Sprintf("%v", mapping.Get("child_created_tmz")))
		chk.Equal("[0 0 2]", fmt.Sprintf("%v", mapping.Get("child_modified_tmz")))
		chk.Equal("[0 1]", fmt.Sprintf("%v", mapping.Get("child_name")))
		chk.Equal("[0 2]", fmt.Sprintf("%v", mapping.Get("child_age")))
		//
		chk.Equal("[1 0 0]", fmt.Sprintf("%v", mapping.Get("parent_pk")))
		chk.Equal("[1 0 1]", fmt.Sprintf("%v", mapping.Get("parent_created_tmz")))
		chk.Equal("[1 0 2]", fmt.Sprintf("%v", mapping.Get("parent_modified_tmz")))
		chk.Equal("[1 1]", fmt.Sprintf("%v", mapping.Get("parent_name")))
		chk.Equal("[1 2]", fmt.Sprintf("%v", mapping.Get("parent_age")))
		//
		chk.Equal("[2 0 0]", fmt.Sprintf("%v", mapping.Get("Emergency_pk")))
		chk.Equal("[2 0 1]", fmt.Sprintf("%v", mapping.Get("Emergency_created_tmz")))
		chk.Equal("[2 0 2]", fmt.Sprintf("%v", mapping.Get("Emergency_modified_tmz")))
		chk.Equal("[2 1]", fmt.Sprintf("%v", mapping.Get("Emergency_name")))
		chk.Equal("[2 2]", fmt.Sprintf("%v", mapping.Get("Emergency_age")))
		//
		// This just increases code coverage.  We're not particularly concerned with the string
		// representation at this point.  As long as it doesn't crash and isn't empty we're happy.
		chk.NotEqual("", mapping.String())
	}
}

func TestMapper_Map_TaggedFieldsOnly(t *testing.T) {
	chk := assert.New(t)
	type Address struct {
		Street string `t:"street"`
		City   string `t:"city"`
		State  string `t:"state"`
		Zip    string `t:"zip"`
		// Should be ignored.
		Zoning string
	}
	type Person struct {
		Name string `t:"name"`
		Age  int    `t:"age"`
		// Should be ignored
		Occupation string
		Address    `t:"address"`
	}
	values := map[string]string{
		"name":       "Bob",
		"age":        "42",
		"Occupation": "Basket Weaving",
		//
		"address_street": "12345 Street",
		"address_city":   "Small City",
		"address_state":  "ST",
		"address_zip":    "98765",
		"address_Zoning": "residential",
	}
	{
		m := &set.Mapper{
			Join: "_",
			Tags: []string{"t"},
		}
		p := &Person{}
		b, err := m.Bind(p)
		chk.NoError(err)
		for k, v := range values {
			b.Set(k, v)
		}
		chk.NoError(b.Err())
		chk.Equal("Bob", p.Name)
		chk.Equal(42, p.Age)
		chk.Equal("Basket Weaving", p.Occupation)
		//
		chk.Equal("12345 Street", p.Address.Street)
		chk.Equal("Small City", p.Address.City)
		chk.Equal("ST", p.Address.State)
		chk.Equal("98765", p.Address.Zip)
		chk.Equal("residential", p.Address.Zoning)
	}
	{
		m := &set.Mapper{
			Join:             "_",
			Tags:             []string{"t"},
			TaggedFieldsOnly: true,
		}
		p := &Person{}
		b, err := m.Bind(p)
		chk.NoError(err)
		//
		occupation, zoning := values["Occupation"], values["address_Zoning"]
		chk.Equal("Basket Weaving", occupation)
		chk.Equal("residential", zoning)
		delete(values, "Occupation")
		delete(values, "address_Zoning")
		//
		for k, v := range values {
			b.Set(k, v)
		}
		chk.NoError(b.Err())
		chk.Equal("Bob", p.Name)
		chk.Equal(42, p.Age)
		chk.Equal("", p.Occupation)
		//
		chk.Equal("12345 Street", p.Address.Street)
		chk.Equal("Small City", p.Address.City)
		chk.Equal("ST", p.Address.State)
		chk.Equal("98765", p.Address.Zip)
		chk.Equal("", p.Address.Zoning)
		//
		err = b.Set("Occupation", occupation)
		chk.ErrorIs(err, set.ErrUnknownField)
		chk.Equal("", p.Occupation)
		err = b.Set("address_Zoning", zoning)
		chk.ErrorIs(err, set.ErrUnknownField)
		chk.Equal("", p.Address.Zoning)
	}
}

func TestMapperBindCollision(t *testing.T) {
	chk := assert.New(t)
	type Db struct {
		Id int
	}
	type A struct {
		A string
		Db
	}
	type B struct {
		B string
		Db
	}
	type T struct {
		A
		B
	}
	{
		var t T
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(Db{}),
			Join:     "_",
		}
		bound, err := mapper.Bind(&t)
		chk.NoError(err)
		bound.Set("A_Id", "15")
		bound.Set("B_Id", 25)
		err = bound.Err()
		chk.NoError(err)
		chk.Equal(15, t.A.Id)
		chk.Equal(25, t.B.Id)
	}
}

func TestMapperCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	{ // Tests case where mapper is empty when calling Mapping.Lookup ~AND~ Mapping.String
		var mapping set.Mapping
		_, ok := mapping.Lookup("Hi")
		chk.Equal(false, ok)
		s := mapping.String()
		chk.Equal("", s)
	}
	{ // Tests case when type T is already scanned and uses Mapper.known
		// ~AND~ mapping by reflect.Type
		type T struct {
			A string
		}
		mapping := set.DefaultMapper.Map(T{})
		chk.NotNil(mapping)
		mapping = set.DefaultMapper.Map(T{})
		chk.NotNil(mapping)
		// by reflect.Type
		m2 := set.DefaultMapper.Map(reflect.TypeOf(T{}))
		chk.NotNil(m2)
		chk.Equal(mapping, m2)
	}
	{ // Tests case when type T is already wrapped in *set.Value when calling Mapper.Map
		type T struct {
			A string
		}
		mapping := set.DefaultMapper.Map(set.V(T{}))
		chk.NotNil(mapping)
	}
	{ // Tests Mapper.Copy
		type T struct {
			A string
			B string
		}
		m1 := set.DefaultMapper.Map(T{})
		chk.NotNil(m1)
		m2 := m1.Copy()
		chk.NotNil(m2)
		chk.Equal(len(m1.Indeces), len(m2.Indeces))
		m2.Indeces["A"] = nil
		m2.Indeces["B"] = nil
		chk.NotEqual(len(m1.Indeces["A"]), len(m2.Indeces["A"]))
		chk.NotEqual(len(m1.Indeces["B"]), len(m2.Indeces["B"]))
	}
	{ // Tests BoundMapping when value V is not a struct.
		var b bool
		bound, err := set.DefaultMapper.Bind(&b)
		chk.NoError(err)
		err = bound.Set("Huh", false)
		chk.ErrorIs(err, set.ErrUnknownField)
	}
	{ // Tests Mapper.Bind when bound value is already a *set.Value and BoundMapping.Rebind when value is already a *set.Value
		type T struct {
			A string
		}
		var t, u T
		vt, vu := set.V(&t), set.V(&u)
		bound, err := set.DefaultMapper.Bind(vt)
		chk.NoError(err)
		bound.Rebind(vu)
	}
	{ // Tests BoundMapping.Set when the underlying set can not be performed.
		type A struct {
			I int
		}
		bound, err := set.DefaultMapper.Bind(&A{})
		chk.NoError(err)
		err = bound.Set("I", "Hello, World!")
		chk.Error(err) // TODO Want to check for a specific error type; requires Value.To() to be updated.
	}
	{ // Tests Mapper.Map for structs inherently treated as scalars, such as time.Time and *time.Time
		type T struct {
			Time  time.Time
			PTime *time.Time
		}
		mapping := set.DefaultMapper.Map(&T{})
		chk.NotNil(mapping)
		//
		field, ok := mapping.Lookup("Time")
		chk.Equal(true, ok)
		chk.NotNil(field)
		field, ok = mapping.Lookup("PTime")
		chk.Equal(true, ok)
		chk.NotNil(field)
	}
}
