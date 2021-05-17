package set_test

import (
	"database/sql"
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
		bound := mapper.Bind(&data)
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
		bound := mapper.Bind(&data)
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

func TestMapper_Bind_Rebind(t *testing.T) {
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
		var a, b Person
		var other CommonDb
		mapper := &set.Mapper{
			Elevated: set.NewTypeList(CommonDb{}),
		}
		bound := mapper.Bind(&a)
		chk.NotNil(bound)
		//
		err := bound.Set("Pk", 10)
		chk.NoError(err)
		chk.Equal(10, a.Pk)
		//
		bound.Rebind(&b)
		err = bound.Set("Pk", 20)
		chk.NoError(err)
		chk.Equal(20, b.Pk)
		//
		chk.NotEqual(a.Pk, b.Pk)
		//
		var didPanic bool
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			bound.Rebind(&other)
		}()
		chk.Equal(true, didPanic)
	}
}

func TestBoundMappingSetFast(t *testing.T) {
	chk := assert.New(t)
	type T struct {
		B   bool
		I   int
		I8  int8
		I16 int16
		I32 int32
		I64 int64
		U   uint
		U8  uint8
		U16 uint16
		U32 uint32
		U64 uint64
		F32 float32
		F64 float64
		S   string
	}
	{
		var t T
		bound := set.DefaultMapper.Bind(&t)
		bound.Set("B", true)
		bound.Set("I", int(-42))
		bound.Set("I8", int8(-8))
		bound.Set("I16", int16(-16))
		bound.Set("I32", int32(-32))
		bound.Set("I64", int64(-64))
		bound.Set("U", uint(42))
		bound.Set("U8", uint8(8))
		bound.Set("U16", uint16(16))
		bound.Set("U32", uint32(32))
		bound.Set("U64", uint64(64))
		bound.Set("F32", float32(3.14))
		bound.Set("F64", float64(6.28))
		bound.Set("S", "string")
		err := bound.Err()
		chk.NoError(err)
		chk.Equal(true, t.B)
		chk.Equal(-42, t.I)
		chk.Equal(int8(-8), t.I8)
		chk.Equal(int16(-16), t.I16)
		chk.Equal(int32(-32), t.I32)
		chk.Equal(int64(-64), t.I64)
		chk.Equal(uint(42), t.U)
		chk.Equal(uint8(8), t.U8)
		chk.Equal(uint16(16), t.U16)
		chk.Equal(uint32(32), t.U32)
		chk.Equal(uint64(64), t.U64)
		chk.Equal(float32(3.14), t.F32)
		chk.Equal(float64(6.28), t.F64)
		chk.Equal("string", t.S)
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
		bound := mapper.Bind(&t)
		bound.Set("A_Id", "15")
		bound.Set("B_Id", 25)
		err := bound.Err()
		chk.NoError(err)
		chk.Equal(15, t.A.Id)
		chk.Equal(25, t.B.Id)
	}
}

func TestMapperCodeCoverage(t *testing.T) {
	chk := assert.New(t)
	{ // Tests case where receiver is nil when calling Mapping.Lookup ~AND~ Mapping.String
		var mapping *set.Mapping
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
		bound := set.DefaultMapper.Bind(&b)
		chk.NotNil(bound)
		err := bound.Set("Huh", false)
		chk.Error(err)
	}
	{ // Tests Mapper.Bind when bound value is already a *set.Value and BoundMapping.Rebind when value is already a *set.Value
		type T struct {
			A string
		}
		var t, u T
		vt, vu := set.V(&t), set.V(&u)
		bound := set.DefaultMapper.Bind(vt)
		chk.NotNil(bound)
		bound.Rebind(vu)
	}
	{ // Tests BoundMapping.Set when the underlying set can not be performed.
		type A struct {
			I int
		}
		bound := set.DefaultMapper.Bind(&A{})
		chk.NotNil(bound)
		err := bound.Set("I", "Hello, World!")
		chk.Error(err)
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

func TestBoundMappingAssignables(t *testing.T) {
	chk := assert.New(t)
	{ // Tests BoundMapping.Assignables
		type A struct {
			A string
			B string
		}
		data := []A{{}, {}}
		for k := 0; k < len(data); k++ {
			bound := set.DefaultMapper.Bind(&data[k])
			chk.NotNil(bound)
			assignables, err := bound.Assignables([]string{"B", "A"}, nil)
			chk.NoError(err)
			chk.NotNil(assignables)
			chk.Equal(2, len(assignables))
			chk.Equal(fmt.Sprintf("%p", &data[k].B), fmt.Sprintf("%p", assignables[0]))
			chk.Equal(fmt.Sprintf("%p", &data[k].A), fmt.Sprintf("%p", assignables[1]))
			chk.Equal(&data[k].B, assignables[0])
			chk.Equal(&data[k].A, assignables[1])
			//
			assignables, err = bound.Assignables([]string{"b", "b"}, nil)
			chk.Error(err)
			chk.Nil(assignables)
			chk.Equal(0, len(assignables))
		}
	}
	{ // Test when bound data is not writable.
		type A struct {
			A string
			B string
		}
		data := A{}
		bound := set.DefaultMapper.Bind(data)
		chk.NotNil(bound)
		assignables, err := bound.Assignables([]string{"B", "A"}, nil)
		chk.Error(err)
		chk.Nil(assignables)
		chk.Equal(0, len(assignables))
	}
	{ // Test pointers for nested/embedded structs instantiated along the way.
		type A struct {
			A string
			B string
		}
		type T struct {
			*A
			Field *A
		}
		data := []T{{}, {}}
		for k := 0; k < len(data); k++ {
			bound := set.DefaultMapper.Bind(&data[k])
			chk.NotNil(bound)
			assignables, err := bound.Assignables([]string{"A_B", "A_A", "Field_A", "Field_B"}, nil)
			chk.NoError(err)
			chk.NotNil(assignables)
			chk.Equal(4, len(assignables))
			// embedded + field
			e, f := data[k].A, data[k].Field
			chk.Equal(fmt.Sprintf("%p", &e.B), fmt.Sprintf("%p", assignables[0]))
			chk.Equal(fmt.Sprintf("%p", &e.A), fmt.Sprintf("%p", assignables[1]))
			chk.Equal(fmt.Sprintf("%p", &f.A), fmt.Sprintf("%p", assignables[2]))
			chk.Equal(fmt.Sprintf("%p", &f.B), fmt.Sprintf("%p", assignables[3]))
			//
			chk.Equal(&e.B, assignables[0])
			chk.Equal(&e.A, assignables[1])
			chk.Equal(&f.A, assignables[2])
			chk.Equal(&f.B, assignables[3])
		}
	}
}

func TestBoundMappingCopy(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct {
		A string
		B string
	}
	var b1, b2 set.BoundMapping
	data := []A{{"a1", "b1"}, {"a2", "b2"}}
	for k := 0; k < len(data); k++ {
		b1 = set.DefaultMapper.Bind(&data[k])
		chk.NotNil(b1)
		fields, err := b1.Fields([]string{"B", "A"}, nil)
		chk.NoError(err)
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		chk.Equal(data[k].B, fields[0])
		chk.Equal(data[k].A, fields[1])
		//
		fields, err = b1.Assignables([]string{"b", "b"}, nil)
		chk.Error(err)
		chk.Nil(fields)
		chk.Equal(0, len(fields))
	}
	//
	b2 = b1.Copy()
	chk.NotNil(b2)
	for k := 0; k < len(data); k++ {
		b2.Rebind(&data[k])
		fields, err := b2.Fields([]string{"B", "A"}, nil)
		chk.NoError(err)
		chk.NotNil(fields)
		chk.Equal(2, len(fields))
		chk.Equal(data[k].B, fields[0])
		chk.Equal(data[k].A, fields[1])
		//
		fields, err = b2.Assignables([]string{"b", "b"}, nil)
		chk.Error(err)
		chk.Nil(fields)
		chk.Equal(0, len(fields))
	}
}

func TestBoundMappingFields(t *testing.T) {
	chk := assert.New(t)
	{ // Tests BoundMapping.Fields
		type A struct {
			A string
			B string
		}
		data := []A{{"a1", "b1"}, {"a2", "b2"}}
		for k := 0; k < len(data); k++ {
			bound := set.DefaultMapper.Bind(&data[k])
			chk.NotNil(bound)
			fields, err := bound.Fields([]string{"B", "A"}, nil)
			chk.NoError(err)
			chk.NotNil(fields)
			chk.Equal(2, len(fields))
			chk.Equal(data[k].B, fields[0])
			chk.Equal(data[k].A, fields[1])
			//
			fields, err = bound.Fields([]string{"b", "b"}, nil)
			chk.Error(err)
			chk.Nil(fields)
			chk.Equal(0, len(fields))
		}
	}
	{ // Test when bound data is not writable.
		type A struct {
			A string
			B string
		}
		data := A{}
		bound := set.DefaultMapper.Bind(data)
		chk.NotNil(bound)
		fields, err := bound.Fields([]string{"B", "A"}, nil)
		chk.Error(err)
		chk.Nil(fields)
		chk.Equal(0, len(fields))
	}
	{ // Test when unknown field is requested.
		type A struct {
			A string
			B string
		}
		data := A{}
		bound := set.DefaultMapper.Bind(&data)
		chk.NotNil(bound)
		fields, err := bound.Fields([]string{"BB", "A"}, nil)
		chk.Error(err)
		chk.Nil(fields)
		chk.Equal(0, len(fields))
	}
	{ // Test pointers for nested/embedded structs instantiated along the way.
		type A struct {
			A string
			B string
		}
		type T struct {
			*A
			Field *A
		}
		data := []T{{}, {}}
		for k := 0; k < len(data); k++ {
			bound := set.DefaultMapper.Bind(&data[k])
			chk.NotNil(bound)
			fields, err := bound.Fields([]string{"A_B", "A_A", "Field_A", "Field_B"}, nil)
			chk.NoError(err)
			chk.NotNil(fields)
			chk.Equal(4, len(fields))
			// embedded + field
			e, f := data[k].A, data[k].Field
			chk.Equal(e.B, fields[0])
			chk.Equal(e.A, fields[1])
			chk.Equal(f.A, fields[2])
			chk.Equal(f.B, fields[3])
			//
			data[k].A.A = fmt.Sprintf("a%v", k)
			data[k].A.B = fmt.Sprintf("b%v", k)
			data[k].Field.A = fmt.Sprintf("field.a%v", k)
			data[k].Field.B = fmt.Sprintf("field.b%v", k)
			fields, err = bound.Fields([]string{"A_B", "A_A", "Field_A", "Field_B"}, nil)
			chk.NoError(err)
			chk.NotNil(fields)
			chk.Equal(4, len(fields))
			chk.Equal(data[k].A.B, fields[0])
			chk.Equal(data[k].A.A, fields[1])
			chk.Equal(data[k].Field.A, fields[2])
			chk.Equal(data[k].Field.B, fields[3])
		}
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
		mapping := mapper.Map(&data)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("lowercase with dot separators")
		mapper := &set.Mapper{
			Join:      ".",
			Transform: strings.ToLower,
		}
		mapping := mapper.Map(&data)
		fmt.Println(strings.Replace(mapping.String(), "\t\t", " ", -1))
	}
	{
		fmt.Println("")
		fmt.Println("specify tags")
		mapper := &set.Mapper{
			Join: "_",
			Tags: []string{"t"},
		}
		mapping := mapper.Map(&data)
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

func ExampleMapper_treatAsScalar() {
	type T struct {
		S string
		T time.Time
		N sql.NullString
	}

	mapping := set.DefaultMapper.Map(T{})
	if mapping.Get("T") != nil {
		fmt.Println("T is mapped because time.Time is automatically treated as a scalar.")
	}
	if mapping.Get("N") == nil {
		fmt.Println("N can not be found because sql.NullString was not treated as a scalar.")
	}
	if mapping.Get("N_Valid") != nil {
		fmt.Println("N_Valid was mapped because the exported fields in sql.NullString were mapped.")
	}

	//
	// Now we'll treat sql.NullString as a scalar when mapping.
	mapper := &set.Mapper{
		TreatAsScalar: set.NewTypeList(sql.NullString{}),
	}
	mapping = mapper.Map(T{})
	if mapping.Get("N") != nil {
		fmt.Println("N is now mapped to the entire sql.NullString member.")
		v, _ := set.V(&T{}).FieldByIndex(mapping.Get("N"))
		fmt.Printf("N's type is %v\n", v.Type())
	}

	// Output: T is mapped because time.Time is automatically treated as a scalar.
	// N can not be found because sql.NullString was not treated as a scalar.
	// N_Valid was mapped because the exported fields in sql.NullString were mapped.
	// N is now mapped to the entire sql.NullString member.
	// N's type is sql.NullString
}

func ExampleMapper_Bind() {
	type CommonDb struct {
		Pk          int
		CreatedTime string
		UpdatedTime string
	}
	type Person struct {
		CommonDb
		Name string
		Age  int
	}
	Print := func(p Person) {
		fmt.Printf("Person: pk=%v created=%v updated=%v name=%v age=%v\n", p.Pk, p.CreatedTime, p.UpdatedTime, p.Name, p.Age)
	}
	data := []Person{{}, {}}
	Print(data[0])
	Print(data[1])
	//
	mapper := &set.Mapper{
		Elevated: set.NewTypeList(CommonDb{}),
	}
	bound := mapper.Bind(&data[0])
	bound.Set("Pk", 10)
	bound.Set("CreatedTime", "-5h")
	bound.Set("UpdatedTime", "-2h")
	bound.Set("Name", "Bob")
	bound.Set("Age", 30)
	if err := bound.Err(); err != nil {
		fmt.Println(err.Error())
	}
	//
	bound.Rebind(&data[1])
	bound.Set("Pk", 20)
	bound.Set("CreatedTime", "-15h")
	bound.Set("UpdatedTime", "-12h")
	bound.Set("Name", "Sally")
	bound.Set("Age", 20)
	if err := bound.Err(); err != nil {
		fmt.Println(err.Error())
	}
	//
	Print(data[0])
	Print(data[1])

	// Output: Person: pk=0 created= updated= name= age=0
	// Person: pk=0 created= updated= name= age=0
	// Person: pk=10 created=-5h updated=-2h name=Bob age=30
	// Person: pk=20 created=-15h updated=-12h name=Sally age=20

}

func ExampleBoundMapping() {
	type Person struct {
		First string
		Last  string
	}
	values := []map[string]string{
		{"first": "Bob", "last": "Smith"},
		{"first": "Sally", "last": "Smith"},
	}
	mapper := &set.Mapper{
		Transform: strings.ToLower,
	}
	var people []Person
	bound := mapper.Bind(&Person{})
	for _, m := range values {
		person := Person{}
		bound.Rebind(&person)
		for fieldName, fieldValue := range m {
			bound.Set(fieldName, fieldValue)
		}
		if err := bound.Err(); err != nil {
			fmt.Println(err.Error())
		}
		people = append(people, person)
	}
	fmt.Println(people[0].First, people[0].Last)
	fmt.Println(people[1].First, people[1].Last)

	// Output: Bob Smith
	// Sally Smith
}
