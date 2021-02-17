package set_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

// TODO RM
// func TestValue_Broken(t *testing.T) {
// 	chk := assert.New(t)
// 	var err error
// 	//
// 	{
// 		var b []*bool
// 		chk.Equal(0, len(b))
// 		err = set.V(&b).Append(true, false)
// 		chk.NoError(err)
// 		chk.Equal(2, len(b))
// 		chk.Equal(true, *b[0])
// 		chk.Equal(false, *b[1])
// 	}
// }

func TestValue_fields(t *testing.T) {
	chk := assert.New(t)
	//
	{
		var b bool
		value := set.V(&b)
		chk.NotNil(value)
		fields := value.Fields()
		chk.Nil(fields)
	}
	{
		var b **bool
		value := set.V(&b)
		chk.NotNil(value)
		fields := value.Fields()
		chk.Nil(fields)
	}
	{
		// Can't set unexported fields.
		type T struct {
			a string
			b string
		}
		var t T
		value := set.V(&t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.Error(err)
		}
		chk.Equal("", t.a)
		chk.Equal("", t.b)
	}
	{
		// Address-of t not passed; can not set fields.
		type T struct {
			a string
			b string
		}
		var t T
		value := set.V(t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.Error(err)
		}
		chk.Equal("", t.a)
		chk.Equal("", t.b)
	}
	{
		// Settable.
		type T struct {
			A string
			B string
		}
		var t T
		value := set.V(&t)
		chk.NotNil(value)
		fields := value.Fields()
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		for _, field := range fields {
			err := field.Value.To(".")
			chk.NoError(err)
		}
		chk.Equal(".", t.A)
		chk.Equal(".", t.B)
	}
}

func TestValue_set(t *testing.T) {
	chk := assert.New(t)
	//
	{ // Only addressable values can be set; passing local variable fails.
		var v bool
		err := set.V(v).To(true)
		chk.Error(err)
	}
	{ // Only addressable values can be set; passing address-of local variable works.
		var v bool
		chk.Equal(false, v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.Equal(true, v)
	}
	{ // The local variable is a pointer but its address is not passed; this fails.
		var v *bool
		chk.Equal((*bool)(nil), v)
		err := set.V(v).To(true)
		chk.Error(err)
		chk.Nil(v)
	}
	{ // The local variable is a pointer and its address is passed; a new pointer is created and assignable.
		var o, v *bool
		chk.Equal((*bool)(nil), v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		o = v // Save o so we can compare if a new pointer is created or existing is reused.
		//
		err = set.V(&v).To(false)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(false, *v)
		chk.Equal(o, v)
	}
	{ // The local variable is an instantiated pointer and we pass its address.
		v := new(bool)
		o := v
		chk.NotNil(v)
		err := set.V(&v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		chk.Equal(o, v)
	}
	{ // The local variable is an instantiated pointer and we do not pass its address.
		v := new(bool)
		o := v
		chk.NotNil(v)
		err := set.V(v).To(true)
		chk.NoError(err)
		chk.NotNil(v)
		chk.Equal(true, *v)
		chk.Equal(o, v)
	}
}

func TestValue_setPointers(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		bp := &b
		err = set.V(bp).To("1")
		chk.NoError(err)
		chk.Equal(true, *bp)
		chk.Equal(true, b)
	}
	{
		var bp *bool
		err = set.V(&bp).To("True")
		chk.NoError(err)
		chk.Equal(true, *bp)
	}
	{
		var bppp ***bool
		err = set.V(&bppp).To("True")
		chk.NoError(err)
		chk.Equal(true, ***bppp)
	}
	{
		var ippp ***int
		s := "42"
		sp := &s
		spp := &sp
		err = set.V(&ippp).To(spp)
		chk.NoError(err)
		chk.Equal(42, ***ippp)
	}
	{
		s := "True"
		var bpp **bool
		err = set.V(&bpp).To(&s)
		chk.NoError(err)
		chk.Equal(true, **bpp)
	}
	{
		s := "True"
		var bp *bool
		err = set.V(&bp).To(&s)
		chk.NoError(err)
		chk.Equal(true, *bp)
	}
	{
		b, s := false, "True"
		bp, sp := &b, &s
		bpp, spp := &bp, &sp
		err = set.V(bpp).To(spp)
		chk.NoError(err)
		chk.Equal(true, b)
	}
}

func TestValue_setSlice(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(b).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(3, len(b))
	}
	{
		i := []int{2, 4, 6}
		chk.Equal(3, len(i))
		err = set.V(&i).To([]interface{}{"Hi"})
		chk.Error(err)
		chk.Equal(0, len(i))
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]interface{}{false, true, false, true})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To([]interface{}{"false", 1, 0, "True"})
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(false, b[0])
		chk.Equal(true, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		b := []bool{true, false, true}
		chk.Equal(3, len(b))
		err = set.V(&b).To("True")
		chk.NoError(err)
		chk.Equal(1, len(b))
		chk.Equal(true, b[0])
		err = set.V(&b).To(0)
		chk.NoError(err)
		chk.Equal(1, len(b))
		chk.Equal(false, b[0])
	}
}

func TestValue_setSliceToBool(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		err = set.V(b).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(false, b)
	}
	{
		var b bool
		err = set.V(&b).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(true, b)
	}
	{
		b := true
		err = set.V(&b).To([]bool{})
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := true
		err = set.V(&b).To(([]bool)(nil))
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := true
		err = set.V(&b).To([]interface{}{true, 1, 0})
		chk.NoError(err)
		chk.Equal(false, b)
	}
	{
		b := false
		err = set.V(&b).To([]interface{}{true, float32(0), "True"})
		chk.NoError(err)
		chk.Equal(true, b)
	}
}

func TestValue_setSliceToInt(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var i int
		err = set.V(i).To([]bool{false, true, false})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal(0, i)
	}
	{
		var i int
		err = set.V(&i).To([]bool{false, true, false, true})
		chk.NoError(err)
		chk.Equal(1, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]int{})
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To(([]int)(nil))
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]interface{}{true, float32(32), float64(0)})
		chk.NoError(err)
		chk.Equal(0, i)
	}
	{
		i := int(42)
		err = set.V(&i).To([]interface{}{true, float32(0), "3.14"})
		chk.NoError(err)
		chk.Equal(3, i)
	}
}

func TestValue_setSliceToString(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var s string
		err = set.V(s).To([]string{"a", "b", "c"})
		// Expect an error because was not &b and therefore original length is not changed either.
		chk.Error(err)
		chk.Equal("", s)
	}
	{
		var s string
		err = set.V(&s).To([]string{"a", "b", "c"})
		chk.NoError(err)
		chk.Equal("c", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]string{})
		chk.NoError(err)
		chk.Equal("", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To(([]string)(nil))
		chk.NoError(err)
		chk.Equal("", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]interface{}{true, float32(32), float64(0), false})
		chk.NoError(err)
		chk.Equal("false", s)
	}
	{
		s := "Hello"
		err = set.V(&s).To([]interface{}{true, float32(0), 64})
		chk.NoError(err)
		chk.Equal("64", s)
	}
}
func TestValue_setSliceCreatesCopies(t *testing.T) {
	chk := assert.New(t)
	//
	{
		slice := []bool{true, true, true}
		dest := []bool{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = false
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []float32{2, 4, 6}
		dest := []float32{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = 42
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []int{2, 4, 6}
		dest := []int{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = -4
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []uint{2, 4, 6}
		dest := []uint{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = 42
		chk.NotEqual(dest[1], slice[1])
	}
	{
		slice := []string{"Hello", "World", "foo"}
		dest := []string{}
		set.V(&dest).To(slice)
		chk.Equal(3, len(slice))
		chk.Equal(3, len(dest))
		dest[1] = "bar"
		chk.NotEqual(dest[1], slice[1])
	}
}

func TestValue_zero(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		v := 42
		value := set.V(v)
		err = value.Zero()
		chk.Error(err)
		chk.Equal(42, v)
	}
	{
		v := 42
		value := set.V(&v)
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, v)
	}
	{
		v := []int{42, -42}
		value := set.V(v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.Error(err)
		chk.Equal(2, len(v))
	}
	{
		v := []int{42, -42}
		value := set.V(&v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(v))
	}
	{
		// test zero followed by append.
		v := []int{42, -42}
		value := set.V(&v)
		chk.Equal(2, len(v))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(v))
		v = append(v, 1000)
		v = append(v, 10000, 100000)
		chk.Equal(3, len(v))
	}
	{
		type Test struct {
			S []string
		}
		instance := &Test{[]string{"Hello", "World"}}
		value := set.V(instance.S)
		chk.Equal(2, len(instance.S))
		err = value.Zero()
		chk.Error(err)
		chk.Equal(2, len(instance.S))
	}
	{
		type Test struct {
			S []string
		}
		instance := &Test{[]string{"Hello", "World"}}
		value := set.V(&instance.S)
		chk.Equal(2, len(instance.S))
		err = value.Zero()
		chk.NoError(err)
		chk.Equal(0, len(instance.S))
	}
}

func TestValue_append(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b bool
		err = set.V(&b).Append(true, false)
		chk.Error(err)
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
	}
	{
		var b []*bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, *b[0])
		chk.Equal(false, *b[1])
	}
	{
		var b []****bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false)
		chk.NoError(err)
		chk.Equal(2, len(b))
		chk.Equal(true, ****b[0])
		chk.Equal(false, ****b[1])
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false, "false", "1")
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
	{
		var b []bool
		chk.Equal(0, len(b))
		err = set.V(&b).Append(true, false, "false", "1")
		chk.NoError(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
		// None of the following are appended.
		err = set.V(&b).Append(true, "asdf", false)
		chk.Error(err)
		chk.Equal(4, len(b))
		chk.Equal(true, b[0])
		chk.Equal(false, b[1])
		chk.Equal(false, b[2])
		chk.Equal(true, b[3])
	}
}

func TestValue_fill(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	type T struct {
		Name string
		Age  uint
	}
	type Tags struct {
		String string `key:"Name"`
		Number uint   `key:"Age"`
	}
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		default:
			return nil
		}
	})
	{
		var t T
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
	}
	{
		var t Tags
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.String)
		chk.Equal(uint(42), t.Number)
	}
}

func TestValue_fillNonStruct(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		default:
			return nil
		}
	})
	{
		var t bool
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
	}
	{
		var t int
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
	}
}

func TestValue_fillNested(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
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
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "Name":
			return "Bob"
		case "Age":
			return "42"
		case "Address":
			return set.GetterFunc(func(key string) interface{} {
				switch key {
				case "Street1":
					return "97531 Some Street"
				case "Street2":
					return ""
				case "City":
					return "Big City"
				case "State":
					return "ST"
				case "Zip":
					return "12345"
				default:
					return nil
				}
			})
		default:
			return nil
		}
	})
	{
		var t Person
		err = set.V(&t).Fill(getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedByTag(t *testing.T) {
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
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	getter := set.GetterFunc(func(key string) interface{} {
		switch key {
		case "name":
			return "Bob"
		case "age":
			return "42"
		case "address":
			return set.GetterFunc(func(key string) interface{} {
				switch key {
				case "street1":
					return "97531 Some Street"
				case "street2":
					return ""
				case "city":
					return "Big City"
				case "state":
					return "ST"
				case "zip":
					return "12345"
				default:
					return nil
				}
			})
		default:
			return nil
		}
	})
	{
		var t Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedByMap(t *testing.T) {
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
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
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
	{
		var t Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedStructSlices(t *testing.T) {
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
		Name    string  `key:"name"`
		Age     uint    `key:"age"`
		Address Address `key:"address"`
	}
	type Company struct {
		Name         string   `key:"name"`
		Employees    []Person `key:"employees"`
		LastEmployee Person   `key:"employees"`
		Slice        []Person `key:"slice"`
	}
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
	{
		var t Company
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		// chk.Equal("Some Company", t.Name)
		// //
		// chk.Equal(2, len(t.Employees))
		// //
		// chk.Equal("Bob", t.Employees[0].Name)
		// chk.Equal(uint(42), t.Employees[0].Age)
		// chk.Equal("97531 Some Street", t.Employees[0].Address.Street1)
		// chk.Equal("", t.Employees[0].Address.Street2)
		// chk.Equal("Big City", t.Employees[0].Address.City)
		// chk.Equal("ST", t.Employees[0].Address.State)
		// chk.Equal("12345", t.Employees[0].Address.Zip)
		// //
		// chk.Equal("Sally", t.Employees[1].Name)
		// chk.Equal(uint(48), t.Employees[1].Age)
		// chk.Equal("555 Small Lane", t.Employees[1].Address.Street1)
		// chk.Equal("", t.Employees[1].Address.Street2)
		// chk.Equal("Other City", t.Employees[1].Address.City)
		// chk.Equal("OO", t.Employees[1].Address.State)
		// chk.Equal("54321", t.Employees[1].Address.Zip)
		// //
		// chk.Equal("Sally", t.LastEmployee.Name)
		// chk.Equal(uint(48), t.LastEmployee.Age)
		// chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
		// chk.Equal("", t.LastEmployee.Address.Street2)
		// chk.Equal("Other City", t.LastEmployee.Address.City)
		// chk.Equal("OO", t.LastEmployee.Address.State)
		// chk.Equal("54321", t.LastEmployee.Address.Zip)
		// //
		// chk.Equal(1, len(t.Slice))
		// chk.Equal("Slice", t.Slice[0].Name)
		// chk.Equal("Slice Street", t.Slice[0].Address.Street1)
		// chk.Equal("", t.Slice[0].Address.Street2)
		// chk.Equal("Slice City", t.Slice[0].Address.City)
		// chk.Equal("SL", t.Slice[0].Address.State)
		// chk.Equal("99999", t.Slice[0].Address.Zip)
	}
}

func TestValue_fillNestedStructSlicesAsPointers(t *testing.T) {
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
		Name         string    `key:"name"`
		Employees    []*Person `key:"employees"`
		LastEmployee *Person   `key:"employees"`
		Slice        []*Person `key:"slice"`
	}
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
	{
		var t *Company
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Some Company", t.Name)
		//
		chk.Equal(2, len(t.Employees))
		//
		chk.Equal("Bob", t.Employees[0].Name)
		chk.Equal(uint(42), t.Employees[0].Age)
		chk.Equal("97531 Some Street", t.Employees[0].Address.Street1)
		chk.Equal("", t.Employees[0].Address.Street2)
		chk.Equal("Big City", t.Employees[0].Address.City)
		chk.Equal("ST", t.Employees[0].Address.State)
		chk.Equal("12345", t.Employees[0].Address.Zip)
		//
		chk.Equal("Sally", t.Employees[1].Name)
		chk.Equal(uint(48), t.Employees[1].Age)
		chk.Equal("555 Small Lane", t.Employees[1].Address.Street1)
		chk.Equal("", t.Employees[1].Address.Street2)
		chk.Equal("Other City", t.Employees[1].Address.City)
		chk.Equal("OO", t.Employees[1].Address.State)
		chk.Equal("54321", t.Employees[1].Address.Zip)
		//
		chk.Equal("Sally", t.LastEmployee.Name)
		chk.Equal(uint(48), t.LastEmployee.Age)
		chk.Equal("555 Small Lane", t.LastEmployee.Address.Street1)
		chk.Equal("", t.LastEmployee.Address.Street2)
		chk.Equal("Other City", t.LastEmployee.Address.City)
		chk.Equal("OO", t.LastEmployee.Address.State)
		chk.Equal("54321", t.LastEmployee.Address.Zip)
		//
		chk.Equal(1, len(t.Slice))
		chk.Equal("Slice", t.Slice[0].Name)
		chk.Equal("Slice Street", t.Slice[0].Address.Street1)
		chk.Equal("", t.Slice[0].Address.Street2)
		chk.Equal("Slice City", t.Slice[0].Address.City)
		chk.Equal("SL", t.Slice[0].Address.State)
		chk.Equal("99999", t.Slice[0].Address.Zip)
	}
}

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
		Employees    *[]*Person `key:"employees"`
		LastEmployee *Person    `key:"employees"`
		Slice        *[]*Person `key:"slice"`
	}
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
	{
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
}

func TestValue_fillNestedPointersByMap(t *testing.T) {
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
	{
		var t *Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.NotNil(t.Address)
		chk.Equal("97531 Some Street", t.Address.Street1)
		chk.Equal("", t.Address.Street2)
		chk.Equal("Big City", t.Address.City)
		chk.Equal("ST", t.Address.State)
		chk.Equal("12345", t.Address.Zip)
	}
}

func TestValue_fillNestedPointersByMapWithNils(t *testing.T) {
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
	m := map[string]interface{}{
		"name": "Bob",
		"age":  42,
	}
	getter := set.MapGetter(m)
	{
		var t *Person
		err = set.V(&t).FillByTag("key", getter)
		chk.NoError(err)
		chk.Equal("Bob", t.Name)
		chk.Equal(uint(42), t.Age)
		chk.NotNil(t.Address)
	}
}

func TestValue_fieldByIndex(t *testing.T) {
	chk := assert.New(t)
	var field *set.Value
	var err error
	{ // No pointers.
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
		type Combined struct {
			Child     Person
			Parent    Person
			Emergency Person `db:"Emergency"`
		}
		var data Combined
		outer := set.V(&data)
		//
		// Combined.Child.Pk
		field, err = outer.FieldByIndex([]int{0, 0, 0})
		chk.NoError(err)
		chk.NotNil(field)
		err = field.To(15)
		chk.NoError(err)
		chk.Equal(15, data.Child.Pk)
		// Combined.Emergency.Name
		field, err = outer.FieldByIndex([]int{2, 1})
		chk.NoError(err)
		chk.NotNil(field)
		err = field.To("Bob")
		chk.NoError(err)
		chk.Equal("Bob", data.Emergency.Name)
	}
	{ // Pointers.
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
		type Combined struct {
			Child     *Person
			Parent    *Person
			Emergency *Person `db:"Emergency"`
		}
		var data Combined
		outer := set.V(&data)
		set := func(indeces []int, arg interface{}) {
			field, err = outer.FieldByIndex(indeces)
			chk.NoError(err)
			chk.NotNil(field)
			err = field.To(arg)
			chk.NoError(err)
		}
		//
		set([]int{0, 0, 0}, 1)
		set([]int{0, 0, 1}, "0c")
		set([]int{0, 0, 2}, "0u")
		set([]int{0, 1}, "Bob")
		set([]int{0, 2}, 5)
		//
		set([]int{1, 0, 0}, 5)
		set([]int{1, 0, 1}, "5c")
		set([]int{1, 0, 2}, "5u")
		set([]int{1, 1}, "Sally")
		set([]int{1, 2}, 30)
		//
		set([]int{2, 0, 0}, 90)
		set([]int{2, 0, 1}, "90c")
		set([]int{2, 0, 2}, "90u")
		set([]int{2, 1}, "Suzy")
		set([]int{2, 2}, 25)
		//
		chk.Equal(1, data.Child.Pk)
		chk.Equal("0c", data.Child.CreatedTime)
		chk.Equal("0u", data.Child.UpdatedTime)
		chk.Equal("Bob", data.Child.Name)
		chk.Equal(5, data.Child.Age)
		//
		chk.Equal(5, data.Parent.Pk)
		chk.Equal("5c", data.Parent.CreatedTime)
		chk.Equal("5u", data.Parent.UpdatedTime)
		chk.Equal("Sally", data.Parent.Name)
		chk.Equal(30, data.Parent.Age)
		//
		chk.Equal(90, data.Emergency.Pk)
		chk.Equal("90c", data.Emergency.CreatedTime)
		chk.Equal("90u", data.Emergency.UpdatedTime)
		chk.Equal("Suzy", data.Emergency.Name)
		chk.Equal(25, data.Emergency.Age)
		//
		// Index out of bounds
		field, err = outer.FieldByIndex([]int{0, 0, 4})
		chk.Nil(field)
		chk.Error(err)
		field, err = outer.FieldByIndex([]int{0, 6})
		chk.Nil(field)
		chk.Error(err)
		field, err = outer.FieldByIndex([]int{10})
		chk.Nil(field)
		chk.Error(err)
	}
}

func TestValue_fieldByIndexCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	var err error
	var value, field *set.Value
	type A struct {
		A string
	}
	type B struct {
		B string
		A
	}
	var a A

	field, err = value.FieldByIndex(nil)
	chk.Error(err)
	chk.Nil(field)
	//
	value = set.V(map[string]string{})
	field, err = value.FieldByIndex(nil)
	chk.Error(err)
	chk.Nil(field)
	//
	value = set.V(a)
	field, err = value.FieldByIndex([]int{0})
	chk.Error(err)
	chk.Nil(field)
	//
	value = set.V(&a)
	field, err = value.FieldByIndex(nil)
	chk.Error(err)
	chk.Nil(field)
	field, err = value.FieldByIndex([]int{})
	chk.Error(err)
	chk.Nil(field)
	//
	{ // Test scalar, aka something not indexable.
		var b bool
		value = set.V(&b)
		field, err := value.FieldByIndex([]int{1, 2})
		chk.Error(err)
		chk.Nil(field)
	}
}

func TestValue_fillCodeCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	m := map[string]interface{}{
		"Nested": map[string]interface{}{
			"String": "Hello, World!",
		},
		"Slice": []map[string]interface{}{
			{
				"String": "Hello, World!",
			},
			{
				"String": "Goodbye, World!",
			},
		},
	}
	getter := set.MapGetter(m)
	{
		type T struct {
			Nested struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested []struct {
				String int
			}
		}
		var t T
		err = set.V(t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested []struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Nested string
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice []struct {
				String int
			}
		}
		var t T
		err = set.V(t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice []struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice struct {
				String int
			}
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
	{
		type T struct {
			Slice int
		}
		var t T
		err = set.V(&t).Fill(getter)
		chk.Error(err)
	}
}

func TestValue_appendCodeCoverageErrors(t *testing.T) {
	chk := assert.New(t)
	//
	var err error
	{
		var b []bool
		err = set.V(b).Append(42)
		chk.Error(err)
	}
}

func TestValue_newElemCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	//
	{ // Tests NewElem when *Value is nil
		var v *set.Value
		elem, err := v.NewElem()
		chk.Error(err)
		chk.Nil(elem)
	}
	{ // Tests NewElem when *Value is not nil but not a map
		var b bool
		v := set.V(&b)
		elem, err := v.NewElem()
		chk.Error(err)
		chk.Nil(elem)
	}
}
