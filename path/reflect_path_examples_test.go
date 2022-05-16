package path_test

import (
	"fmt"
	"reflect"

	"github.com/nofeaturesonlybugs/set/path"
)

func ExampleReflectPath() {
	type Other struct {
		Message string
	}
	type Foo struct {
		Num int
		Str string
		M   Other
	}

	f := Foo{}

	var rp path.ReflectPath
	tree := path.Stat(Foo{})

	v := reflect.Indirect(reflect.ValueOf(&f))

	rp = tree.Leaves["Str"].ReflectPath()
	rp.Value(v).SetString("Blue")

	rp = tree.Leaves["Num"].ReflectPath()
	rp.Value(v).SetInt(42)

	rp = tree.Leaves["M.Message"].ReflectPath()
	rp.Value(v).SetString("hut hut")

	fmt.Println(f.Str, f.Num, f.M.Message)

	// Output: Blue 42 hut hut
}

func ExampleReflectPath_pointer() {
	type Foo struct {
		Num int
		Str string
	}
	type Ptr struct {
		P *Foo
	}

	f := Ptr{}

	var rp path.ReflectPath
	tree := path.Stat(f)

	v := reflect.Indirect(reflect.ValueOf(&f))

	rp = tree.Leaves["P.Str"].ReflectPath()
	rp.Value(v).SetString("Blue")

	rp = tree.Leaves["P.Num"].ReflectPath()
	rp.Value(v).SetInt(42)

	fmt.Println(f.P.Str, f.P.Num)

	// Output: Blue 42
}
