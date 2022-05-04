package set_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

// MapField indicates which struct field to set and what to expect.
type MapField struct {
	// The field to access or set.
	Field string

	// The value to set.
	To interface{}

	// The expected result or error.
	Expect interface{}
	Error  error
}

// MapStructTest helps build table driven tests for testing our mapper on various inputs.
type MapStructTest struct {
	// Name of the test.
	Name string

	// The value to map with Bind and Prepare.
	V interface{}

	// If non-nil the expected error from Bind or Prepare.
	MapperError error

	// The mapped Fields to set or access.
	Fields []MapField
}

// MapStructTests is a slice of tests to run.
type MapStructTests []MapStructTest

// Run runs the tests in MapStructTests.
func (tests MapStructTests) Run(t *testing.T, m *set.Mapper) {
	if m == nil {
		// Default if not provided is basic mapper joining with DOT
		m = &set.Mapper{
			Join: ".",
		}
	}
	//
	// NewV creates a copy of the test.V so that each test runs on a separate instance.
	NewV := func(v interface{}) interface{} {
		var rv interface{}
		typeOf := reflect.TypeOf(v)
		switch typeOf.Kind() {
		case reflect.Struct:
			rv = reflect.Zero(typeOf).Interface()
		case reflect.Ptr:
			rv = reflect.New(typeOf.Elem()).Interface()
		}
		return rv
	}
	for _, test := range tests {
		describe := fmt.Sprintf("test=%v type=%v", test.Name, reflect.TypeOf(test.V))
		fieldNames := make([]string, len(test.Fields))
		hasUnknownField := false
		for k, field := range test.Fields {
			fieldNames[k] = field.Field
			hasUnknownField = field.Error == set.ErrUnknownField
		}
		//
		if test.MapperError != nil {
			t.Run("Bind "+test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				b, err := m.Bind(NewV(test.V))
				chk.ErrorIs(err, test.MapperError, describe)
				//
				ptrs, err := b.Assignables(fieldNames, nil)
				chk.ErrorIs(err, set.ErrReadOnly, describe)
				chk.Nil(ptrs, describe)
				//
				for _, field := range fieldNames {
					v, err := b.Field(field)
					chk.ErrorIs(err, set.ErrReadOnly, describe)
					chk.Nil(v, describe)
				}
				//
				values, err := b.Fields(fieldNames, nil)
				chk.ErrorIs(err, set.ErrReadOnly, describe)
				chk.Nil(values, describe)
				//
				for _, field := range test.Fields {
					err := b.Set(field.Field, field.To)
					chk.ErrorIs(err, set.ErrReadOnly, describe)
				}
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				p, err := m.Prepare(NewV(test.V))
				chk.ErrorIs(err, test.MapperError, describe)
				//
				err = p.Plan(fieldNames...)
				chk.ErrorIs(err, set.ErrReadOnly, describe)
				//
				ptrs, err := p.Assignables(nil)
				chk.ErrorIs(err, set.ErrReadOnly, describe)
				chk.Nil(ptrs, describe)
				//
				for range fieldNames {
					v, err := p.Field()
					chk.ErrorIs(err, set.ErrReadOnly, describe)
					chk.Nil(v, describe)
				}
				//
				values, err := p.Fields(nil)
				chk.ErrorIs(err, set.ErrReadOnly, describe)
				chk.Nil(values, describe)
				//
				for _, field := range test.Fields {
					err = p.Set(field.To)
					chk.ErrorIs(err, set.ErrReadOnly, describe)
				}
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			t.Run("Bind "+test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				b, err := m.Bind(NewV(test.V))
				chk.ErrorIs(err, test.MapperError, describe)
				//
				_, err = b.Assignables(fieldNames, nil)
				chk.ErrorIs(err, set.ErrUnknownField, describe)
				//
				for _, field := range test.Fields {
					v, err := b.Field(field.Field)
					chk.ErrorIs(err, field.Error, describe)
					if err != nil {
						continue
					}
					err = v.To(field.To)
					chk.NoError(err, describe)
					chk.Equal(field.Expect, v.WriteValue.Interface(), describe)
				}
				//
				_, err = b.Fields(fieldNames, nil)
				chk.ErrorIs(err, set.ErrUnknownField, describe)
				//
				for _, field := range test.Fields {
					err = b.Set(field.Field, field.To)
					chk.ErrorIs(err, field.Error, describe)
					if err != nil {
						continue
					}
					v, _ := b.Field(field.Field)
					chk.Equal(field.Expect, v.WriteValue.Interface(), describe)
				}
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				p, err := m.Prepare(NewV(test.V))
				chk.ErrorIs(err, test.MapperError, describe)
				//
				err = p.Plan(fieldNames...)
				chk.ErrorIs(err, set.ErrUnknownField, describe)
				//
				ptrs, err := p.Assignables(nil)
				chk.ErrorIs(err, set.ErrNoPlan, describe)
				chk.Nil(ptrs, describe)
				//
				for range fieldNames {
					v, err := p.Field()
					chk.ErrorIs(err, set.ErrNoPlan, describe)
					chk.Nil(v, describe)
				}
				//
				values, err := p.Fields(nil)
				chk.ErrorIs(err, set.ErrNoPlan, describe)
				chk.Nil(values, describe)
				//
				for _, field := range test.Fields {
					err = p.Set(field.To)
					chk.ErrorIs(err, set.ErrNoPlan, describe)
				}
			})
			continue
		}
		t.Run("Bind Assignables "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			b, err := m.Bind(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			ptrs, err := b.Assignables(fieldNames, nil)
			chk.NoError(err, describe)
			chk.Equal(len(fieldNames), len(ptrs), describe)
		})
		t.Run("Prepare Assignables "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			p, err := m.Prepare(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err, describe)
			//
			ptrs, err := p.Assignables(nil)
			chk.NoError(err, describe)
			chk.Equal(len(fieldNames), len(ptrs), describe)
		})
		t.Run("Bind Field "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			b, err := m.Bind(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			for _, field := range test.Fields {
				v, err := b.Field(field.Field)
				chk.ErrorIs(err, field.Error, describe)
				if err != nil {
					continue
				}
				err = v.To(field.To)
				chk.NoError(err, describe)
				chk.Equal(field.Expect, v.WriteValue.Interface(), describe)
			}
		})
		t.Run("Prepare Field "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			p, err := m.Prepare(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err, describe)
			//
			for _, field := range test.Fields {
				v, err := p.Field()
				chk.ErrorIs(err, field.Error, describe)
				if err != nil {
					continue
				}
				err = v.To(field.To)
				chk.NoError(err, describe)
				chk.Equal(field.Expect, v.WriteValue.Interface(), describe)
			}
		})
		t.Run("Bind Fields "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			b, err := m.Bind(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			values, err := b.Fields(fieldNames, nil)
			chk.NoError(err, describe)
			chk.Equal(len(fieldNames), len(values), describe)
		})
		t.Run("Prepare Fields "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			p, err := m.Prepare(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err, describe)
			//
			values, err := p.Fields(nil)
			chk.NoError(err, describe)
			chk.Equal(len(fieldNames), len(values), describe)
		})
		t.Run("Bind Set "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			b, err := m.Bind(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			for _, field := range test.Fields {
				err = b.Set(field.Field, field.To)
				chk.ErrorIs(err, field.Error, describe)
				if err != nil {
					continue
				}
				v, _ := b.Field(field.Field)
				chk.Equal(field.Expect, v.WriteValue.Interface(), describe)
			}
		})
		t.Run("Prepare Set "+test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			p, err := m.Prepare(NewV(test.V))
			chk.ErrorIs(err, test.MapperError, describe)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err)
			//
			for _, field := range test.Fields {
				err = p.Set(field.To)
				chk.ErrorIs(err, field.Error, describe)
				if err != nil {
					continue
				}
				// TODO Need a way to confirm field was set correctly
			}
		})
	}
}

// Benchmark runs the tests as benchmarks.
func (tests MapStructTests) Benchmark(B *testing.B) {
	m := &set.Mapper{
		Join: ".",
	}
	//
	// NewV creates a copy of the test.V so that each benchmark runs on a separate instance.
	NewV := func(v interface{}) interface{} {
		return v
		// NB  We don't want extra allocations during these benchmarks.
		// var rv interface{}
		// typeOf := reflect.TypeOf(v)
		// switch typeOf.Kind() {
		// case reflect.Struct:
		// 	rv = reflect.Zero(typeOf).Interface()
		// case reflect.Ptr:
		// 	rv = reflect.New(typeOf.Elem()).Interface()
		// }
		// return rv
	}
	for _, test := range tests {
		//
		fieldNames := make([]string, len(test.Fields))
		hasUnknownField := false
		for k, field := range test.Fields {
			fieldNames[k] = field.Field
			hasUnknownField = field.Error == set.ErrUnknownField
		}
		//
		if test.MapperError != nil {
			B.Run("Bind "+test.Name, func(b *testing.B) {
				bound, _ := m.Bind(NewV(test.V))
				for n := 0; n < b.N; n++ {
					bound.Rebind(NewV(test.V))
					bound.Assignables(fieldNames, nil)
					bound.Fields(fieldNames, nil)
					for _, field := range test.Fields {
						bound.Field(field.Field)
						bound.Set(field.Field, field.To)
					}
				}
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				p, _ := m.Prepare(NewV(test.V))
				p.Plan(fieldNames...)
				for n := 0; n < b.N; n++ {
					p.Rebind(NewV(test.V))
					p.Assignables(nil)
					p.Fields(nil)
					for _, field := range test.Fields {
						p.Field()
						p.Set(field.To)
					}
				}
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			B.Run("Bind "+test.Name, func(b *testing.B) {
				bound, _ := m.Bind(NewV(test.V))
				slice := make([]interface{}, len(fieldNames))
				for n := 0; n < b.N; n++ {
					bound.Rebind(NewV(test.V))
					bound.Assignables(fieldNames, slice)
					bound.Fields(fieldNames, slice)
					for _, field := range test.Fields {
						bound.Field(field.Field)
						bound.Set(field.Field, field.To)
					}
				}
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				p, _ := m.Prepare(NewV(test.V))
				p.Plan(fieldNames...)
				for n := 0; n < b.N; n++ {
					p.Rebind(NewV(test.V))
					p.Assignables(nil)
					p.Fields(nil)
					for range fieldNames {
						p.Field()
					}
					for _, field := range test.Fields {
						p.Set(field.To)
					}
				}
			})
			continue
		}
		B.Run("Bind Assignables "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(NewV(test.V))
			ptrs := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				bound.Rebind(NewV(test.V))
				bound.Assignables(fieldNames, ptrs)
			}
		})
		B.Run("Prepare Assignables "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(NewV(test.V))
			p.Plan(fieldNames...)
			ptrs := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				p.Rebind(NewV(test.V))
				p.Assignables(ptrs)
			}
		})
		B.Run("Bind Fields "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(NewV(test.V))
			values := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				bound.Rebind(NewV(test.V))
				bound.Fields(fieldNames, values)
			}
		})
		B.Run("Prepare Fields "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(NewV(test.V))
			p.Plan(fieldNames...)
			values := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				p.Rebind(NewV(test.V))
				p.Fields(values)
			}
		})
		B.Run("Bind Field "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(NewV(test.V))
			for n := 0; n < b.N; n++ {
				bound.Rebind(NewV(test.V))
				for _, field := range test.Fields {
					bound.Field(field.Field)
				}
			}
		})
		B.Run("Prepare Field "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(NewV(test.V))
			p.Plan(fieldNames...)
			for n := 0; n < b.N; n++ {
				p.Rebind(NewV(test.V))
				for range test.Fields {
					p.Field()
				}
			}
		})
		B.Run("Bind Set "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(NewV(test.V))
			for n := 0; n < b.N; n++ {
				bound.Rebind(NewV(test.V))
				for _, field := range test.Fields {
					bound.Set(field.Field, field.To)
				}
			}
		})
		B.Run("Prepare Set "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(NewV(test.V))
			p.Plan(fieldNames...)
			for n := 0; n < b.N; n++ {
				p.Rebind(NewV(test.V))
				for _, field := range test.Fields {
					p.Set(field.To)
				}
			}
		})
	}
}

// SimpleStruct is a simple type to test our mapping types.
type SimpleStruct struct {
	Str string
	Int int
}

// NestedStruct contains a level of nesting for testing our mapping types.
type NestedStruct struct {
	SimpleStruct
	Next SimpleStruct
}

// NestedPtrStruct contains a level of nesting via pointers for testing our mapping types.
type NestedPtrStruct struct {
	*SimpleStruct
	Next *SimpleStruct
}

// PrimitivesStruct contains all of the basic primitives as fields.
type PrimitivesStruct struct {
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

func Test_Mapper_BindPrepare(t *testing.T) {
	var simple SimpleStruct
	var nested NestedStruct
	var nestedPtr NestedPtrStruct
	var primitives PrimitivesStruct
	//
	tests := []MapStructTest{
		//
		// Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable
		//
		{
			Name:        "unaddr1",
			V:           SimpleStruct{},
			MapperError: set.ErrReadOnly,
			Fields: []MapField{
				{Field: "Str", To: "Hi", Error: set.ErrReadOnly},
			},
		},
		{
			Name:        "unaddr2",
			V:           NestedStruct{},
			MapperError: set.ErrReadOnly,
			Fields: []MapField{
				{Field: "Next.Str", To: "Hi", Error: set.ErrReadOnly},
			},
		},
		//
		// Unknown-Field Unknown-Field Unknown-Field Unknown-Field Unknown-Field Unknown-Field
		//
		{
			Name: "unknown1",
			V:    &simple,
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		{
			Name: "unknown2",
			V:    &nested,
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		//
		// Successful Successful Successful Successful Successful Successful Successful Successful
		//
		{
			Name: "simple",
			V:    &simple,
			Fields: []MapField{
				{Field: "Str", To: "Hi", Expect: "Hi"},
				{Field: "Int", To: "42", Expect: 42},
			},
		},
		{
			Name: "nested",
			V:    &nested,
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
		},
		{
			Name: "ptr nested",
			V:    &nestedPtr,
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
		},
		{
			// BoundMapping.Set() and PreparedMapping.Set() should use the "fast paths" type switch
			// to set these values.
			Name: "primitives",
			V:    &primitives,
			Fields: []MapField{
				{Field: "B", To: true, Expect: true},
				{Field: "I", To: int(-42), Expect: int(-42)},
				{Field: "I8", To: int8(-8), Expect: int8(-8)},
				{Field: "I16", To: int16(-16), Expect: int16(-16)},
				{Field: "I32", To: int32(-32), Expect: int32(-32)},
				{Field: "I64", To: int64(-64), Expect: int64(-64)},
				{Field: "U", To: uint(42), Expect: uint(42)},
				{Field: "U8", To: uint8(8), Expect: uint8(8)},
				{Field: "U16", To: uint16(16), Expect: uint16(16)},
				{Field: "U32", To: uint32(32), Expect: uint32(32)},
				{Field: "U64", To: uint64(64), Expect: uint64(64)},
				{Field: "F32", To: float32(3.14), Expect: float32(3.14)},
				{Field: "F64", To: float64(6.28), Expect: float64(6.28)},
				{Field: "S", To: "Cheerio", Expect: "Cheerio"},
			},
		},
	}
	//
	tests = append(tests, CreateTestsSale()...)
	//
	MapStructTests(tests).Run(t, nil)
}

func Benchmark_Mapper_BindPrepare(b *testing.B) {
	var simple SimpleStruct
	var nested NestedStruct
	var nestedPtr NestedPtrStruct
	var primitives PrimitivesStruct
	tests := []MapStructTest{
		//
		// Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable
		//
		{
			Name:        "unaddr1",
			V:           SimpleStruct{},
			MapperError: set.ErrReadOnly,
			Fields: []MapField{
				{Field: "Str", To: "Hi", Error: set.ErrReadOnly},
			},
		},
		{
			Name:        "unaddr2",
			V:           NestedStruct{},
			MapperError: set.ErrReadOnly,
			Fields: []MapField{
				{Field: "Next.Str", To: "Hi", Error: set.ErrReadOnly},
			},
		},
		//
		// Unknown-Field Unknown-Field Unknown-Field Unknown-Field Unknown-Field Unknown-Field
		//
		{
			Name: "unknown1",
			V:    &simple,
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		{
			Name: "unknown2",
			V:    &nested,
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		//
		// Successful Successful Successful Successful Successful Successful Successful Successful
		//
		{
			Name: "simple",
			V:    &simple,
			Fields: []MapField{
				{Field: "Str", To: "Hi", Expect: "Hi"},
				{Field: "Int", To: "42", Expect: 42},
			},
		},
		{
			Name: "nested",
			V:    &nested,
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
		},
		{
			Name: "ptr nested",
			V:    &nestedPtr,
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
		},
		{
			// BoundMapping.Set() and PreparedMapping.Set() should use the "fast paths" type switch
			// to set these values.
			Name: "primitives",
			V:    &primitives,
			Fields: []MapField{
				{Field: "B", To: true, Expect: true},
				{Field: "I", To: int(-42), Expect: int(-42)},
				{Field: "I8", To: int8(-8), Expect: int8(-8)},
				{Field: "I16", To: int16(-16), Expect: int16(-16)},
				{Field: "I32", To: int32(-32), Expect: int32(-32)},
				{Field: "I64", To: int64(-64), Expect: int64(-64)},
				{Field: "U", To: uint(42), Expect: uint(42)},
				{Field: "U8", To: uint8(8), Expect: uint8(8)},
				{Field: "U16", To: uint16(16), Expect: uint16(16)},
				{Field: "U32", To: uint32(32), Expect: uint32(32)},
				{Field: "U64", To: uint64(64), Expect: uint64(64)},
				{Field: "F32", To: float32(3.14), Expect: float32(3.14)},
				{Field: "F64", To: float64(6.28), Expect: float64(6.28)},
				{Field: "S", To: "Cheerio", Expect: "Cheerio"},
			},
		},
	}
	//
	tests = append(tests, CreateTestsSale()...)
	//
	MapStructTests(tests).Benchmark(b)
}

func CreateTestsSale() []MapStructTest {
	//
	type Common struct {
		Id int
	}
	type Timestamps struct {
		CreatedTime  string
		ModifiedTime string
	}
	type Person struct {
		Common
		Timestamps // Not used but present anyways
		First      string
		Last       string
	}
	type Vendor struct {
		Common
		Timestamps  // Not used but present anyways
		Name        string
		Description string
		Contact     Person
	}
	type Sale struct {
		Common
		Timestamps
		//
		Price    int
		Quantity int
		Total    int
		//
		Customer Person
		Vendor   Vendor
	}
	var sale Sale
	created := time.Now().Add(-20 * time.Minute)
	modified := created.Add(5 * time.Minute)
	test := []MapStructTest{
		{
			Name: "sale",
			V:    &sale,
			Fields: []MapField{
				{Field: "Common.Id", To: 4, Expect: 4},
				{Field: "Timestamps.CreatedTime", To: created.Format("2006-01-02 15:04:05"), Expect: created.Format("2006-01-02 15:04:05")},
				{Field: "Timestamps.ModifiedTime", To: modified.Format("2006-01-02 15:04:05"), Expect: modified.Format("2006-01-02 15:04:05")},
				{Field: "Price", To: "10.00", Expect: 10},
				{Field: "Quantity", To: "5", Expect: 5},
				{Field: "Total", To: "50.00", Expect: 50},
				{Field: "Customer.Common.Id", To: 42, Expect: 42},
				{Field: "Customer.First", To: "John", Expect: "John"},
				{Field: "Customer.Last", To: "Smith", Expect: "Smith"},
				{Field: "Vendor.Common.Id", To: 4242, Expect: 4242},
				{Field: "Vendor.Name", To: "Neat Widgets Inc.", Expect: "Neat Widgets Inc."},
				{Field: "Vendor.Description", To: "Sales neat widgets.", Expect: "Sales neat widgets."},
				{Field: "Vendor.Contact.Common.Id", To: 424242, Expect: 424242},
				{Field: "Vendor.Contact.First", To: "Jane", Expect: "Jane"},
				{Field: "Vendor.Contact.Last", To: "Doe", Expect: "Doe"},
			},
		},
	}
	return test
}
