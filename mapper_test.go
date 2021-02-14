package set_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
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
		mapping, err := set.DefaultMapper.Map(&data)
		chk.NoError(err)
		chk.Equal("[0]", fmt.Sprintf("%v", mapping.Get("A")))
		chk.Equal("[1]", fmt.Sprintf("%v", mapping.Get("B")))
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
		mapping, err := mapper.Map(&data)
		chk.NoError(err)
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

func TestMapper_Bind(t *testing.T) {
	chk := assert.New(t)
	//
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
		var data Person
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(CommonDb{}),
		}
		bound, err := mapper.Bind(&data)
		chk.NoError(err)
		chk.NotNil(bound)
		//
		field, err := bound.Field("Pk")
		chk.NoError(err)
		chk.NotNil(field)
		err = field.To(10)
		chk.NoError(err)
		chk.Equal(10, data.Pk)
		//
		err = bound.Set("Pk", "20")
		chk.NoError(err)
		chk.Equal(20, data.Pk)
		//
		err = bound.Set("CreatedTime", "created")
		chk.NoError(err)
		chk.Equal("created", data.CreatedTime)
		//
		err = bound.Set("UpdatedTime", "updated")
		chk.NoError(err)
		chk.Equal("updated", data.UpdatedTime)
		//
		field, err = bound.Field("NotFound")
		chk.Error(err)
		chk.Nil(field)
	}
	//
	{ // Test case where CommonDb is embedded but pointer and data is pointer and we pass address of data
		type CommonDb struct {
			Pk          int    `db:"pk" json:"id"`
			CreatedTime string `db:"created_tmz" json:"created_time"`
			UpdatedTime string `db:"modified_tmz" json:"modified_time"`
		}
		type Person struct {
			*CommonDb
			Name string
			Age  int
		}
		var data *Person
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(CommonDb{}),
		}
		bound, err := mapper.Bind(&data)
		chk.NoError(err)
		chk.NotNil(bound)
		//
		field, err := bound.Field("Pk")
		chk.NoError(err)
		chk.NotNil(field)
		err = field.To(10)
		chk.NoError(err)
		chk.Equal(10, data.Pk)
		//
		err = bound.Set("Pk", "20")
		chk.NoError(err)
		chk.Equal(20, data.Pk)
		//
		err = bound.Set("CreatedTime", "created")
		chk.NoError(err)
		chk.Equal("created", data.CreatedTime)
		//
		err = bound.Set("UpdatedTime", "updated")
		chk.NoError(err)
		chk.Equal("updated", data.UpdatedTime)
		//
		field, err = bound.Field("NotFound")
		chk.Error(err)
		chk.Nil(field)
	}
}

func TestMapperCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	{ // Tests case where receiver is nil when calling Mapping.Lookup
		var mapping set.Mapping
		_, ok := mapping.Lookup("Hi")
		chk.Equal(false, ok)
	}
	{ // Tests case where receiver is nil when calling Mapper.Map or Mapper.Bind
		var mapper *set.Mapper
		mapping, err := mapper.Map(struct{}{})
		chk.Error(err)
		chk.Nil(mapping)
		bound, err := mapper.Bind(struct{}{})
		chk.Error(err)
		chk.Nil(bound)
	}
	{ // Tests case when type T is already scanned and uses Mapper.known
		type T struct {
			A string
		}
		mapping, err := set.DefaultMapper.Map(T{})
		chk.NoError(err)
		chk.NotNil(mapping)
		mapping, err = set.DefaultMapper.Map(T{})
		chk.NoError(err)
		chk.NotNil(mapping)
	}
	{ // Tests case when type T is already wrapped in *set.Value when calling Mapper.Map
		type T struct {
			A string
		}
		mapping, err := set.DefaultMapper.Map(set.V(T{}))
		chk.NoError(err)
		chk.NotNil(mapping)
	}
	{ // Tests Mapper.Copy
		type T struct {
			A string
			B string
		}
		m1, err := set.DefaultMapper.Map(T{})
		chk.NoError(err)
		chk.NotNil(m1)
		m2 := m1.Copy()
		chk.NotNil(m2)
		chk.Equal(len(m1), len(m2))
		m2["A"] = nil
		m2["B"] = nil
		chk.NotEqual(len(m1["A"]), len(m2["A"]))
		chk.NotEqual(len(m1["B"]), len(m2["B"]))
	}
	{ // Tests Mapping.Copy when Mapping is nil
		var m1 set.Mapping
		m2 := m1.Copy()
		chk.Nil(m1)
		chk.Nil(m2)
	}
	{ // Tests BoundMapper when value V is not a struct.
		var b bool
		bound, err := set.DefaultMapper.Bind(&b)
		chk.NoError(err)
		chk.NotNil(bound)
		err = bound.Set("Huh", false)
		chk.Error(err)
	}
}

func ExampleMapper() {
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
		mapping, _ := mapper.Map(&data)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("lowercase with dot separators")
		mapper := &set.Mapper{
			Join:      ".",
			Transform: strings.ToLower,
		}
		mapping, _ := mapper.Map(&data)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("specify tags")
		mapper := &set.Mapper{
			Join: "_",
			Tags: []string{"t"},
		}
		mapping, _ := mapper.Map(&data)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
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
