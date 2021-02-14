package set_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

func TestStructMapping(t *testing.T) {
	chk := assert.New(t)
	{
		// Default options test to increase coverage.
		type A struct {
			A string
			B string
		}
		var data A
		mapping := set.NewStructMapping(set.V(&data), nil)
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
		opts := &set.StructMappingOptions{
			Ignored:   []interface{}{Logger{}},
			Elevated:  []interface{}{CommonDb{}},
			Tags:      []string{"db", "json"},
			Join:      "_",
			Transform: strings.ToLower,
		}
		mapping := set.NewStructMapping(set.V(&data), opts)
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

func TestStructMappingCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	var mapping *set.StructMapping
	_, ok := mapping.Lookup("Hi")
	chk.Equal(false, ok)
}

func ExampleStructMapping() {
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
		opts := &set.StructMappingOptions{
			Elevated: []interface{}{CommonDb{}},
			Join:     "_",
		}
		mapping := set.NewStructMapping(set.V(&data), opts)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("lowercase with dot separators")
		opts := &set.StructMappingOptions{
			Join:      ".",
			Transform: strings.ToLower,
		}
		mapping := set.NewStructMapping(set.V(&data), opts)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("specify tags")
		opts := &set.StructMappingOptions{
			Join: "_",
			Tags: []string{"t"},
		}
		mapping := set.NewStructMapping(set.V(&data), opts)
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
