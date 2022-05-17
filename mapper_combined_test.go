package set_test

import (
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

	// New is the factory method that creates a new instance of the value to Bind or Prepare.
	New func() interface{}

	// If non-nil the expected error from Bind or Prepare.
	MapperError error

	// The mapped Fields to set or access.
	Fields []MapField

	AssignablesFn func(value interface{}) []interface{}
}

// BoundMappingMapperError tests Mapper.Bind when the Bind method returns an error.
func (mt MapStructTest) BoundMappingMapperError(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var b set.BoundMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			b, err = m.Bind(mt.New())
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(mt.New())
		}
		//
		ptrs, err := b.Assignables(fieldNames, nil)
		chk.ErrorIs(err, set.ErrReadOnly)
		chk.Nil(ptrs)
		//
		for _, field := range fieldNames {
			v, err := b.Field(field)
			chk.ErrorIs(err, set.ErrReadOnly)
			chk.False(v.TopValue.IsValid())
		}
		//
		values, err := b.Fields(fieldNames, nil)
		chk.ErrorIs(err, set.ErrReadOnly)
		chk.Nil(values)
		//
		for _, field := range mt.Fields {
			err := b.Set(field.Field, field.To)
			chk.ErrorIs(err, set.ErrReadOnly)
		}
	}
}

// BoundMappingUnknownField tests BoundMapping methods when an unknown field is present.
func (mt MapStructTest) BoundMappingUnknownField(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var b set.BoundMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			b, err = m.Bind(mt.New())
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(mt.New())
		}
		//
		_, err = b.Assignables(fieldNames, nil)
		chk.ErrorIs(err, set.ErrUnknownField)
		//
		for _, field := range mt.Fields {
			v, err := b.Field(field.Field)
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
			err = v.To(field.To)
			chk.NoError(err)
			chk.Equal(field.Expect, v.WriteValue.Interface())
		}
		//
		_, err = b.Fields(fieldNames, nil)
		chk.ErrorIs(err, set.ErrUnknownField)
		//
		for _, field := range mt.Fields {
			err = b.Set(field.Field, field.To)
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
			v, _ := b.Field(field.Field)
			chk.Equal(field.Expect, v.WriteValue.Interface())
		}
	}
}

// BoundMappingAssignables tests BoundMapping.Assignables method for this test.
func (mt MapStructTest) BoundMappingAssignables(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var b set.BoundMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		v := mt.New()
		if k == 0 {
			b, err = m.Bind(v)
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(v)
		}
		//
		ptrs, err := b.Assignables(fieldNames, nil)
		chk.NoError(err)
		chk.Equal(len(fieldNames), len(ptrs))
		if fn := mt.AssignablesFn; fn != nil {
			chk.Equal(fn(v), ptrs)
		}
	}
}

// BoundMappingField tests BoundMapping.Field method for this test.
func (mt MapStructTest) BoundMappingField(t *testing.T, m *set.Mapper) {
	chk := assert.New(t)
	var b set.BoundMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			b, err = m.Bind(mt.New())
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(mt.New())
		}
		//
		for _, field := range mt.Fields {
			v, err := b.Field(field.Field)
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
			err = v.To(field.To)
			chk.NoError(err)
			chk.Equal(field.Expect, v.WriteValue.Interface())
		}
	}
}

// BoundMappingFields tests BoundMapping.Fields method for this test.
func (mt MapStructTest) BoundMappingFields(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var b set.BoundMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			b, err = m.Bind(mt.New())
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(mt.New())
		}
		//
		values, err := b.Fields(fieldNames, nil)
		chk.NoError(err)
		chk.Equal(len(fieldNames), len(values))
	}
}

// BoundMappingSet tests BoundMapping.Set method for this test.
func (mt MapStructTest) BoundMappingSet(t *testing.T, m *set.Mapper) {
	chk := assert.New(t)
	var b set.BoundMapping
	var v set.Value
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			b, err = m.Bind(mt.New())
			chk.ErrorIs(err, mt.MapperError)
		} else {
			b.Rebind(mt.New())
		}
		for _, field := range mt.Fields {
			err = b.Set(field.Field, field.To)
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
			v, err = b.Field(field.Field)
			chk.NoError(err)
			chk.Equal(field.Expect, v.WriteValue.Interface())
		}
	}
}

// PreparedMappingMapperError tests Mapper.Prepare when an error is expected during Prepare.
func (mt MapStructTest) PreparedMappingMapperError(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			p, err = m.Prepare(mt.New())
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.ErrorIs(err, set.ErrReadOnly)
		} else {
			p.Rebind(mt.New())
		}
		//
		ptrs, err := p.Assignables(nil)
		chk.ErrorIs(err, set.ErrReadOnly)
		chk.Nil(ptrs)
		//
		for range fieldNames {
			v, err := p.Field()
			chk.ErrorIs(err, set.ErrReadOnly)
			chk.False(v.TopValue.IsValid())
		}
		//
		values, err := p.Fields(nil)
		chk.ErrorIs(err, set.ErrReadOnly)
		chk.Nil(values)
		//
		for _, field := range mt.Fields {
			err = p.Set(field.To)
			chk.ErrorIs(err, set.ErrReadOnly)
		}
	}
}

// PreparedMappingUnknownField tests PreparedMapping methods when an unknown field is present.
func (mt MapStructTest) PreparedMappingUnknownField(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			p, err = m.Prepare(mt.New())
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.ErrorIs(err, set.ErrUnknownField)
		} else {
			p.Rebind(mt.New())
		}
		//
		ptrs, err := p.Assignables(nil)
		chk.ErrorIs(err, set.ErrNoPlan)
		chk.Nil(ptrs)
		//
		for range fieldNames {
			v, err := p.Field()
			chk.ErrorIs(err, set.ErrNoPlan)
			chk.False(v.TopValue.IsValid())
		}
		//
		values, err := p.Fields(nil)
		chk.ErrorIs(err, set.ErrNoPlan)
		chk.Nil(values)
		//
		for _, field := range mt.Fields {
			err = p.Set(field.To)
			chk.ErrorIs(err, set.ErrNoPlan)
		}
	}
}

// PreparedMappingAssignables tests PreparedMapping.Assignables method for this test.
func (mt MapStructTest) PreparedMappingAssignables(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		v := mt.New()
		if k == 0 {
			p, err = m.Prepare(v)
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err)
		} else {
			p.Rebind(v)
		}
		//
		ptrs, err := p.Assignables(nil)
		chk.NoError(err)
		chk.Equal(len(fieldNames), len(ptrs))
		if fn := mt.AssignablesFn; fn != nil {
			chk.Equal(fn(v), ptrs)
		}
	}
}

// PreparedMappingField tests PreparedMapping.Field method for this test.
func (mt MapStructTest) PreparedMappingField(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			p, err = m.Prepare(mt.New())
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err)
		} else {
			p.Rebind(mt.New())
		}
		//
		for _, field := range mt.Fields {
			v, err := p.Field()
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
			err = v.To(field.To)
			chk.NoError(err)
			chk.Equal(field.Expect, v.WriteValue.Interface())
		}
	}
}

// PreparedMappingFields tests PreparedMapping.Fields method for this test.
func (mt MapStructTest) PreparedMappingFields(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var err error
	//
	for k := 0; k < 2; k++ {
		if k == 0 {
			p, err = m.Prepare(mt.New())
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err)
		} else {
			p.Rebind(mt.New())
		}
		//
		values, err := p.Fields(nil)
		chk.NoError(err)
		chk.Equal(len(fieldNames), len(values))
	}
}

// PreparedMappingSet tests PreparedMapping.Set method for this test.
func (mt MapStructTest) PreparedMappingSet(t *testing.T, m *set.Mapper, fieldNames ...string) {
	chk := assert.New(t)
	var p set.PreparedMapping
	var fv set.Value
	var err error
	//
	for k := 0; k < 2; k++ {
		v := mt.New()
		if k == 0 {
			p, err = m.Prepare(v)
			chk.ErrorIs(err, mt.MapperError)
			//
			err = p.Plan(fieldNames...)
			chk.NoError(err)
		} else {
			p.Rebind(v)
		}
		//
		for _, field := range mt.Fields {
			err = p.Set(field.To)
			chk.ErrorIs(err, field.Error)
			if err != nil {
				continue
			}
		}
		// TODO Need better way to reset plan counter
		p.Rebind(v) // Resets plan counter
		for _, field := range mt.Fields {
			fv, err = p.Field()
			chk.NoError(err)
			chk.Equal(field.Expect, fv.WriteValue.Interface())
		}
	}
}

// MapStructTests is a slice of tests to run.
type MapStructTests []MapStructTest

// Run runs the tests in MapStructTests.
func (tests MapStructTests) Run(t *testing.T, m *set.Mapper) {
	if m == nil {
		// Default if not provided is basic mapper joining with DOT
		m = &set.Mapper{
			TreatAsScalar: set.NewTypeList(Scalar{}),
			Join:          ".",
		}
	}
	//
	for _, test := range tests {
		fieldNames := make([]string, len(test.Fields))
		hasUnknownField := false
		for k, field := range test.Fields {
			fieldNames[k] = field.Field
			hasUnknownField = field.Error == set.ErrUnknownField
		}
		//
		if test.MapperError != nil {
			t.Run("Bind "+test.Name, func(t *testing.T) {
				test.BoundMappingMapperError(t, m, fieldNames...)
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				test.PreparedMappingMapperError(t, m, fieldNames...)
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			t.Run("Bind "+test.Name, func(t *testing.T) {
				test.BoundMappingUnknownField(t, m, fieldNames...)
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				test.PreparedMappingUnknownField(t, m, fieldNames...)
			})
			continue
		}
		t.Run("Bind Assignables "+test.Name, func(t *testing.T) {
			test.BoundMappingAssignables(t, m, fieldNames...)
		})
		t.Run("Prepare Assignables "+test.Name, func(t *testing.T) {
			test.PreparedMappingAssignables(t, m, fieldNames...)
		})
		t.Run("Bind Field "+test.Name, func(t *testing.T) {
			test.BoundMappingField(t, m)
		})
		t.Run("Prepare Field "+test.Name, func(t *testing.T) {
			test.PreparedMappingField(t, m, fieldNames...)
		})
		t.Run("Bind Fields "+test.Name, func(t *testing.T) {
			test.BoundMappingFields(t, m, fieldNames...)
		})
		t.Run("Prepare Fields "+test.Name, func(t *testing.T) {
			test.PreparedMappingFields(t, m, fieldNames...)
		})
		t.Run("Bind Set "+test.Name, func(t *testing.T) {
			test.BoundMappingSet(t, m)
		})
		t.Run("Prepare Set "+test.Name, func(t *testing.T) {
			test.PreparedMappingSet(t, m, fieldNames...)
		})
	}
}

// RunParallel runs the tests in MapStructTests in parallel.
func (tests MapStructTests) RunParallel(t *testing.T, m *set.Mapper) {
	if m == nil {
		// Default if not provided is basic mapper joining with DOT
		m = &set.Mapper{
			TreatAsScalar: set.NewTypeList(Scalar{}),
			Join:          ".",
		}
	}
	//
	for _, test := range tests {
		test := test // Capture the test variable; this is necessary during parallel tests.
		fieldNames := make([]string, len(test.Fields))
		hasUnknownField := false
		for k, field := range test.Fields {
			fieldNames[k] = field.Field
			hasUnknownField = field.Error == set.ErrUnknownField
		}
		//
		if test.MapperError != nil {
			t.Run("Bind "+test.Name, func(t *testing.T) {
				t.Parallel()
				test.BoundMappingMapperError(t, m, fieldNames...)
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				t.Parallel()
				test.PreparedMappingMapperError(t, m, fieldNames...)
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			t.Run("Bind "+test.Name, func(t *testing.T) {
				t.Parallel()
				test.BoundMappingUnknownField(t, m, fieldNames...)
			})
			t.Run("Prepare "+test.Name, func(t *testing.T) {
				t.Parallel()
				test.PreparedMappingUnknownField(t, m, fieldNames...)
			})
			continue
		}
		t.Run("Bind Assignables "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.BoundMappingAssignables(t, m, fieldNames...)
		})
		t.Run("Prepare Assignables "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.PreparedMappingAssignables(t, m, fieldNames...)
		})
		t.Run("Bind Field "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.BoundMappingField(t, m)
		})
		t.Run("Prepare Field "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.PreparedMappingField(t, m, fieldNames...)
		})
		t.Run("Bind Fields "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.BoundMappingFields(t, m, fieldNames...)
		})
		t.Run("Prepare Fields "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.PreparedMappingFields(t, m, fieldNames...)
		})
		t.Run("Bind Set "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.BoundMappingSet(t, m)
		})
		t.Run("Prepare Set "+test.Name, func(t *testing.T) {
			t.Parallel()
			test.PreparedMappingSet(t, m, fieldNames...)
		})
	}
}

// Benchmark runs the tests as benchmarks.
func (tests MapStructTests) Benchmark(B *testing.B) {
	m := &set.Mapper{
		Join: ".",
	}
	//
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
				bound, _ := m.Bind(test.New())
				for n := 0; n < b.N; n++ {
					bound.Rebind(test.New())
					_, _ = bound.Assignables(fieldNames, nil)
					_, _ = bound.Fields(fieldNames, nil)
					for _, field := range test.Fields {
						_, _ = bound.Field(field.Field)
						_ = bound.Set(field.Field, field.To)
					}
				}
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				for n := 0; n < b.N; n++ {
					p.Rebind(test.New())
					_, _ = p.Assignables(nil)
					_, _ = p.Fields(nil)
					for _, field := range test.Fields {
						_, _ = p.Field()
						_ = p.Set(field.To)
					}
				}
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			B.Run("Bind "+test.Name, func(b *testing.B) {
				bound, _ := m.Bind(test.New())
				slice := make([]interface{}, len(fieldNames))
				for n := 0; n < b.N; n++ {
					bound.Rebind(test.New())
					_, _ = bound.Assignables(fieldNames, slice)
					_, _ = bound.Fields(fieldNames, slice)
					for _, field := range test.Fields {
						_, _ = bound.Field(field.Field)
						_ = bound.Set(field.Field, field.To)
					}
				}
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				for n := 0; n < b.N; n++ {
					p.Rebind(test.New())
					_, _ = p.Assignables(nil)
					_, _ = p.Fields(nil)
					for range fieldNames {
						_, _ = p.Field()
					}
					for _, field := range test.Fields {
						_ = p.Set(field.To)
					}
				}
			})
			continue
		}
		B.Run("Bind Assignables "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(test.New())
			ptrs := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				bound.Rebind(test.New())
				_, _ = bound.Assignables(fieldNames, ptrs)
			}
		})
		B.Run("Prepare Assignables "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(test.New())
			_ = p.Plan(fieldNames...)
			ptrs := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				p.Rebind(test.New())
				_, _ = p.Assignables(ptrs)
			}
		})
		B.Run("Bind Fields "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(test.New())
			values := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				bound.Rebind(test.New())
				_, _ = bound.Fields(fieldNames, values)
			}
		})
		B.Run("Prepare Fields "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(test.New())
			_ = p.Plan(fieldNames...)
			values := make([]interface{}, len(fieldNames))
			for n := 0; n < b.N; n++ {
				p.Rebind(test.New())
				_, _ = p.Fields(values)
			}
		})
		B.Run("Bind Field "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(test.New())
			for n := 0; n < b.N; n++ {
				bound.Rebind(test.New())
				for _, field := range test.Fields {
					_, _ = bound.Field(field.Field)
				}
			}
		})
		B.Run("Prepare Field "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(test.New())
			_ = p.Plan(fieldNames...)
			for n := 0; n < b.N; n++ {
				p.Rebind(test.New())
				for range test.Fields {
					_, _ = p.Field()
				}
			}
		})
		B.Run("Bind Set "+test.Name, func(b *testing.B) {
			bound, _ := m.Bind(test.New())
			for n := 0; n < b.N; n++ {
				bound.Rebind(test.New())
				for _, field := range test.Fields {
					_ = bound.Set(field.Field, field.To)
				}
			}
		})
		B.Run("Prepare Set "+test.Name, func(b *testing.B) {
			p, _ := m.Prepare(test.New())
			_ = p.Plan(fieldNames...)
			for n := 0; n < b.N; n++ {
				p.Rebind(test.New())
				for _, field := range test.Fields {
					_ = p.Set(field.To)
				}
			}
		})
	}
}

// Benchmark runs the tests as benchmarks.
func (tests MapStructTests) BenchmarkParallel(B *testing.B) {
	m := &set.Mapper{
		Join: ".",
	}
	//
	for _, test := range tests {
		test := test // Capture the table test row
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
				b.RunParallel(func(pb *testing.PB) {
					bound, _ := m.Bind(test.New())
					for pb.Next() {
						bound.Rebind(test.New())
						_, _ = bound.Assignables(fieldNames, nil)
						_, _ = bound.Fields(fieldNames, nil)
						for _, field := range test.Fields {
							_, _ = bound.Field(field.Field)
							_ = bound.Set(field.Field, field.To)
						}
					}
				})
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					p, _ := m.Prepare(test.New())
					_ = p.Plan(fieldNames...)
					for pb.Next() {
						p.Rebind(test.New())
						_, _ = p.Assignables(nil)
						_, _ = p.Fields(nil)
						for _, field := range test.Fields {
							_, _ = p.Field()
							_ = p.Set(field.To)
						}
					}
				})
			})
			continue
		} else if hasUnknownField {
			// When an unknown field is present we expect specific behaviors.
			B.Run("Bind "+test.Name, func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					bound, _ := m.Bind(test.New())
					slice := make([]interface{}, len(fieldNames))
					for pb.Next() {
						bound.Rebind(test.New())
						_, _ = bound.Assignables(fieldNames, slice)
						_, _ = bound.Fields(fieldNames, slice)
						for _, field := range test.Fields {
							_, _ = bound.Field(field.Field)
							_ = bound.Set(field.Field, field.To)
						}
					}
				})
			})
			B.Run("Prepare "+test.Name, func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					p, _ := m.Prepare(test.New())
					_ = p.Plan(fieldNames...)
					for pb.Next() {
						p.Rebind(test.New())
						_, _ = p.Assignables(nil)
						_, _ = p.Fields(nil)
						for range fieldNames {
							_, _ = p.Field()
						}
						for _, field := range test.Fields {
							_ = p.Set(field.To)
						}
					}
				})
			})
			continue
		}
		B.Run("Bind Assignables "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				bound, _ := m.Bind(test.New())
				ptrs := make([]interface{}, len(fieldNames))
				for pb.Next() {
					bound.Rebind(test.New())
					_, _ = bound.Assignables(fieldNames, ptrs)
				}
			})
		})
		B.Run("Prepare Assignables "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				ptrs := make([]interface{}, len(fieldNames))
				for pb.Next() {
					p.Rebind(test.New())
					_, _ = p.Assignables(ptrs)
				}
			})
		})
		B.Run("Bind Fields "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				bound, _ := m.Bind(test.New())
				values := make([]interface{}, len(fieldNames))
				for pb.Next() {
					bound.Rebind(test.New())
					_, _ = bound.Fields(fieldNames, values)
				}
			})
		})
		B.Run("Prepare Fields "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				values := make([]interface{}, len(fieldNames))
				for pb.Next() {
					p.Rebind(test.New())
					_, _ = p.Fields(values)
				}
			})
		})
		B.Run("Bind Field "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				bound, _ := m.Bind(test.New())
				for pb.Next() {
					bound.Rebind(test.New())
					for _, field := range test.Fields {
						_, _ = bound.Field(field.Field)
					}
				}
			})
		})
		B.Run("Prepare Field "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				for pb.Next() {
					p.Rebind(test.New())
					for range test.Fields {
						_, _ = p.Field()
					}
				}
			})
		})
		B.Run("Bind Set "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				bound, _ := m.Bind(test.New())
				for pb.Next() {
					bound.Rebind(test.New())
					for _, field := range test.Fields {
						_ = bound.Set(field.Field, field.To)
					}
				}
			})
		})
		B.Run("Prepare Set "+test.Name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				p, _ := m.Prepare(test.New())
				_ = p.Plan(fieldNames...)
				for pb.Next() {
					p.Rebind(test.New())
					for _, field := range test.Fields {
						_ = p.Set(field.To)
					}
				}
			})
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

// Scalar is for testing Mapper.TreatAsScalar
type Scalar struct {
	N int
	S string
}

// ScalarParent is the mapped struct for testing Mapper.TreatAsScalar.
type ScalarParent struct {
	X, Y int
	S    Scalar
}

func Test_Mapper_BindPrepare(t *testing.T) {
	tests := append(CreateCombinedMapperTests(), CreateCombinedMapperSaletests()...)
	MapStructTests(tests).Run(t, nil)
}

func Test_Mapper_BindPrepareParallel(t *testing.T) {
	tests := append(CreateCombinedMapperTests(), CreateCombinedMapperSaletests()...)
	MapStructTests(tests).RunParallel(t, nil)
}

func Benchmark_Mapper_BindPrepare(b *testing.B) {
	tests := append(CreateCombinedMapperTests(), CreateCombinedMapperSaletests()...)
	MapStructTests(tests).Benchmark(b)
}

func Benchmark_Mapper_BindPrepareParallel(b *testing.B) {
	tests := append(CreateCombinedMapperTests(), CreateCombinedMapperSaletests()...)
	MapStructTests(tests).BenchmarkParallel(b)
}

func CreateCombinedMapperTests() []MapStructTest {
	return []MapStructTest{
		//
		// Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable Unaddressable
		//
		{
			Name:        "unaddr1",
			New:         func() interface{} { return SimpleStruct{} },
			MapperError: set.ErrReadOnly,
			Fields: []MapField{
				{Field: "Str", To: "Hi", Error: set.ErrReadOnly},
			},
		},
		{
			Name:        "unaddr2",
			New:         func() interface{} { return NestedStruct{} },
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
			New:  func() interface{} { return &SimpleStruct{} },
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		{
			Name: "unknown2",
			New:  func() interface{} { return &NestedStruct{} },
			Fields: []MapField{
				{Field: "FIELD DNE", To: "Hi", Error: set.ErrUnknownField},
			},
		},
		//
		// Successful Successful Successful Successful Successful Successful Successful Successful
		//
		{
			Name: "simple",
			New:  func() interface{} { return &SimpleStruct{} },
			Fields: []MapField{
				{Field: "Str", To: "Hi", Expect: "Hi"},
				{Field: "Int", To: "42", Expect: 42},
			},
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*SimpleStruct)
				return []interface{}{&v.Str, &v.Int}
			},
		},
		{
			Name: "nested",
			New:  func() interface{} { return &NestedStruct{} },
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*NestedStruct)
				return []interface{}{&v.SimpleStruct.Str, &v.SimpleStruct.Int, &v.Next.Str, &v.Next.Int}
			},
		},
		{
			Name: "ptr nested",
			New:  func() interface{} { return &NestedPtrStruct{} },
			Fields: []MapField{
				{Field: "SimpleStruct.Str", To: "Bye", Expect: "Bye"},
				{Field: "SimpleStruct.Int", To: "100", Expect: 100},
				{Field: "Next.Str", To: 42, Expect: "42"},
				{Field: "Next.Int", To: float32(999), Expect: 999},
			},
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*NestedPtrStruct)
				return []interface{}{&v.SimpleStruct.Str, &v.SimpleStruct.Int, &v.Next.Str, &v.Next.Int}
			},
		},
		{
			// BoundMapping.Set() and PreparedMapping.Set() should use the "fast paths" type switch
			// to set these values.
			Name: "primitives",
			New:  func() interface{} { return &PrimitivesStruct{} },
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
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*PrimitivesStruct)
				return []interface{}{
					&v.B,
					&v.I, &v.I8, &v.I16, &v.I32, &v.I64,
					&v.U, &v.U8, &v.U16, &v.U32, &v.U64,
					&v.F32, &v.F64,
					&v.S,
				}
			},
		},
		{
			// Testing TreatAsScalar
			Name: "treat as scalar",
			New:  func() interface{} { return &ScalarParent{} },
			Fields: []MapField{
				{Field: "X", To: 10, Expect: 10},
				{Field: "Y", To: 20, Expect: 20},
				{
					Field: "S",
					To: Scalar{
						N: 100,
						S: "Treated as scalar!",
					},
					Expect: Scalar{
						N: 100,
						S: "Treated as scalar!",
					},
				},
			},
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*ScalarParent)
				return []interface{}{&v.X, &v.Y, &v.S}
			},
		},
	}
}

func CreateCombinedMapperSaletests() []MapStructTest {
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
	created := time.Now().Add(-20 * time.Minute)
	modified := created.Add(5 * time.Minute)
	test := []MapStructTest{
		{
			Name: "sale",
			New:  func() interface{} { return &Sale{} },
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
			AssignablesFn: func(value interface{}) []interface{} {
				v := value.(*Sale)
				return []interface{}{
					&v.Id, &v.CreatedTime, &v.ModifiedTime,
					&v.Price, &v.Quantity, &v.Total,
					&v.Customer.Id, &v.Customer.First, &v.Customer.Last,
					&v.Vendor.Id, &v.Vendor.Name, &v.Vendor.Description, &v.Vendor.Contact.Id, &v.Vendor.Contact.First, &v.Vendor.Contact.Last,
				}
			},
		},
	}
	return test
}
