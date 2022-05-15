package set_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestPreparedMapping_Assignables(t *testing.T) {
	Addr := func(p interface{}) string {
		return fmt.Sprintf("%p", p)
	}
	type Test struct {
		Name   string
		V      interface{}
		Fields []string
		Expect []interface{}
		Error  error
	}
	type S struct {
		A string
		B int
	}
	type Nested struct {
		S
		Next S
	}
	var s S
	var n Nested
	tests := []Test{
		{
			Name:   "&s one",
			V:      &s,
			Fields: []string{"A", "B"},
			Expect: []interface{}{&s.A, &s.B},
		},
		{
			Name:   "&s two",
			V:      &s,
			Fields: []string{"B", "A"},
			Expect: []interface{}{&s.B, &s.A},
		},
		// Nesting
		{
			Name:   "n one",
			V:      &n,
			Fields: []string{"S_A", "S_B", "Next_B", "Next_A"},
			Expect: []interface{}{&n.S.A, &n.S.B, &n.Next.B, &n.Next.A},
		},
		{
			Name:   "n two",
			V:      &n,
			Fields: []string{"Next_A", "Next_B", "S_B", "S_A"},
			Expect: []interface{}{&n.Next.A, &n.Next.B, &n.S.B, &n.S.A},
		},
		// Unrecognized field
		{
			Name:   "unknown field",
			V:      &n,
			Fields: []string{"Next_A", "Next_B", "Unknown Field"},
			Error:  set.ErrUnknownField,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			p, err := set.DefaultMapper.Prepare(test.V)
			chk.NoError(err)
			//
			err = p.Plan(test.Fields...)
			chk.ErrorIs(err, test.Error)
			//
			ptrs, err := p.Assignables(nil)
			chk.Equal(test.Expect, ptrs)
			if errors.Is(test.Error, set.ErrUnknownField) {
				chk.ErrorIs(err, set.ErrNoPlan)
			} else {
				chk.NoError(err)
				for k := range test.Expect {
					chk.Equal(Addr(test.Expect[k]), Addr(ptrs[k]))
				}
			}
		})
	}
	//
	t.Run("readonly", func(t *testing.T) {
		// Unaddressable S results in error
		tests := []Test{
			{
				Name:   "s one",
				V:      s,
				Fields: []string{"A", "B"},
			},
			{
				Name:   "s two",
				V:      s,
				Fields: []string{"B", "A"},
			},
		}
		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				p, err := set.DefaultMapper.Prepare(test.V)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				err = p.Plan(test.Fields...)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				_, err = p.Assignables(nil)
				chk.ErrorIs(err, set.ErrReadOnly)
			})
		}
	})
	//
	t.Run("make ptrs", func(t *testing.T) {
		type A struct {
			A string
			B string
		}
		type T struct {
			*A
			Field *A
		}
		//
		chk := assert.New(t)
		//
		data := []T{
			{
				A: &A{
					A: "0.A",
					B: "0.B",
				},
			},
			{
				A: &A{
					A: "1.A",
					B: "1.B",
				},
			},
		}
		// We pass over each element in data twice in this test.
		// The first pass for each element creates the Field *A pointer with zero values
		// but assigns the values from A to Field.
		// The second pass checks these are the correct values.
		for k := 0; k < 4; k++ {
			n := k % 2
			p, err := set.DefaultMapper.Prepare(&data[n])
			chk.NoError(err)
			err = p.Plan([]string{"A_B", "A_A", "Field_A", "Field_B"}...)
			chk.NoError(err)
			ptrs, err := p.Assignables(nil)
			chk.NoError(err)
			chk.NotNil(ptrs)
			chk.Equal(4, len(ptrs))
			// embedded + field
			e, f := data[n].A, data[n].Field
			chk.Equal(Addr(&e.B), Addr(ptrs[0]))
			chk.Equal(Addr(&e.A), Addr(ptrs[1]))
			chk.Equal(Addr(&f.A), Addr(ptrs[2]))
			chk.Equal(Addr(&f.B), Addr(ptrs[3]))
			//
			chk.Equal(&e.B, ptrs[0])
			chk.Equal(&e.A, ptrs[1])
			chk.Equal(&f.A, ptrs[2])
			chk.Equal(&f.B, ptrs[3])
			//
			chk.Equal(fmt.Sprintf("%v.B", n), reflect.Indirect(reflect.ValueOf(ptrs[0])).Interface())
			chk.Equal(fmt.Sprintf("%v.A", n), reflect.Indirect(reflect.ValueOf(ptrs[1])).Interface())
			//
			if k < 2 {
				chk.Equal("", reflect.Indirect(reflect.ValueOf(ptrs[2])).Interface())
				chk.Equal("", reflect.Indirect(reflect.ValueOf(ptrs[3])).Interface())
				data[k].Field = data[k].A
			} else {
				chk.Equal(fmt.Sprintf("%v.A", n), reflect.Indirect(reflect.ValueOf(ptrs[2])).Interface())
				chk.Equal(fmt.Sprintf("%v.B", n), reflect.Indirect(reflect.ValueOf(ptrs[3])).Interface())
			}
		}
	})
}

func TestPreparedMapping_Copy(t *testing.T) {
	chk := assert.New(t)
	//
	type S struct {
		A string
		B int
	}
	var first, copied S
	var err error

	p, err := set.DefaultMapper.Prepare(&first)
	chk.NoError(err)
	err = p.Plan("A", "B")
	chk.NoError(err)

	err = p.Set("Hello")
	chk.NoError(err)
	chk.Nil(p.Err())
	err = p.Set(42)
	chk.NoError(err)
	chk.Nil(p.Err())

	cp := p.Copy()
	cp.Rebind(&copied)
	err = cp.Set("Copied!")
	chk.NoError(err)
	chk.Nil(cp.Err())
	err = cp.Set(100)
	chk.NoError(err)
	chk.Nil(cp.Err())

	chk.Equal("Hello", first.A)
	chk.Equal(42, first.B)
	chk.Equal("Copied!", copied.A)
	chk.Equal(100, copied.B)
}

func TestPreparedMapping_Err(t *testing.T) {
	chk := assert.New(t)
	type S struct {
		A int
	}
	mapper := &set.Mapper{}
	var s, o S
	var err error
	p, err := mapper.Prepare(&s)
	chk.NoError(err)
	err = p.Plan("A")
	chk.NoError(err)

	t.Run("err is set", func(t *testing.T) {
		chk := assert.New(t)
		//
		err = p.Set(42)
		chk.NoError(err)
		// cause error
		err = p.Set("does not exist")
		chk.ErrorIs(err, set.ErrPlanOutOfBounds)
	})
	// This test depends on p's state from the previous test.
	t.Run("plan clears error", func(t *testing.T) {
		chk.ErrorIs(p.Err(), set.ErrPlanOutOfBounds)
		err = p.Plan("A") // Should clear error
		chk.NoError(err)
		chk.Nil(p.Err())
		err = p.Set(48)
		chk.NoError(err)
		chk.Nil(p.Err())
		// cause error
		err = p.Set("does not exist")
		chk.ErrorIs(err, set.ErrPlanOutOfBounds)
	})
	// This test depends on p's state from the previous test.
	t.Run("rebind clears error", func(t *testing.T) {
		chk := assert.New(t)
		//
		chk.ErrorIs(p.Err(), set.ErrPlanOutOfBounds)
		p.Rebind(&o)
		chk.Nil(p.Err())
	})
	// This test does not depend on previous state.
	t.Run("invalid value", func(t *testing.T) {
		chk := assert.New(t)
		//
		err := p.Plan("A")
		chk.NoError(err)
		err = p.Set("Hello")
		chk.Error(err)
	})
}

func TestPreparedMapping_Field(t *testing.T) {
	type S struct {
		A string
		B int
	}
	type Nested struct {
		S
		Next S
	}
	var a, b S
	var m, n Nested

	var v set.Value
	var p set.PreparedMapping
	var err error

	t.Run("S", func(t *testing.T) {
		chk := assert.New(t)
		//
		p, err = set.DefaultMapper.Prepare(&a)
		chk.NoError(err)
		err = p.Plan("B", "A")
		chk.NoError(err)

		v, err = p.Field()
		chk.NoError(err)
		v.To("42")
		v, err = p.Field()
		chk.NoError(err)
		v.To("First")
		v, err = p.Field()
		chk.ErrorIs(err, set.ErrPlanOutOfBounds)
		chk.ErrorIs(p.Err(), set.ErrPlanOutOfBounds)
		chk.False(v.TopValue.IsValid())

		p.Rebind(&b)
		chk.Nil(p.Err())
		v, err = p.Field()
		chk.NoError(err)
		v.To("78")
		v, err = p.Field()
		chk.NoError(err)
		v.To("Second")
		v, err = p.Field()
		chk.ErrorIs(err, set.ErrPlanOutOfBounds)
		chk.ErrorIs(p.Err(), set.ErrPlanOutOfBounds)
		chk.False(v.TopValue.IsValid())

		chk.Equal("First", a.A)
		chk.Equal(42, a.B)
		chk.Equal("Second", b.A)
		chk.Equal(78, b.B)
	})
	//
	t.Run("Nested", func(t *testing.T) {
		chk := assert.New(t)
		//
		p, err = set.DefaultMapper.Prepare(&m)
		chk.NoError(err)
		err = p.Plan("Next_B", "Next_A", "S_A", "S_B")
		chk.NoError(err)

		v, _ = p.Field()
		v.To(100)
		v, _ = p.Field()
		v.To("m.Next.A")
		v, _ = p.Field()
		v.To("m.S.A")
		v, _ = p.Field()
		v.To(10)
		v, err = p.Field()
		chk.ErrorIs(err, set.ErrPlanOutOfBounds)
		chk.ErrorIs(p.Err(), set.ErrPlanOutOfBounds)

		p.Rebind(&n)
		chk.Nil(p.Err())
		v, _ = p.Field()
		v.To(900)
		v, _ = p.Field()
		v.To("n.Next.A")
		v, _ = p.Field()
		v.To("n.S.A")
		v, _ = p.Field()
		v.To(90)
		v, err = p.Field()

		chk.Equal("m.S.A", m.S.A)
		chk.Equal(10, m.S.B)
		chk.Equal("m.Next.A", m.Next.A)
		chk.Equal(100, m.Next.B)

		chk.Equal("n.S.A", n.S.A)
		chk.Equal(90, n.S.B)
		chk.Equal("n.Next.A", n.Next.A)
		chk.Equal(900, n.Next.B)
	})
	//
	t.Run("readonly", func(t *testing.T) {
		chk := assert.New(t)
		//
		p, err := set.DefaultMapper.Prepare(a)
		chk.ErrorIs(err, set.ErrReadOnly)

		err = p.Plan("A", "B")
		chk.ErrorIs(err, set.ErrReadOnly)

		_, err = p.Field()
		chk.ErrorIs(err, set.ErrReadOnly)
	})
	//
	t.Run("invalid", func(t *testing.T) {
		chk := assert.New(t)
		//
		p, err := set.DefaultMapper.Prepare(&a)
		chk.NoError(err)

		_, err = p.Field()
		chk.ErrorIs(err, set.ErrNoPlan)
	})
}

func TestPreparedMapping_Fields(t *testing.T) {
	type Test struct {
		Name   string
		V      interface{}
		Fields []string
		Expect []interface{}
		Error  error
	}
	type S struct {
		A string
		B int
	}
	type Nested struct {
		S
		Next S
	}
	//
	type Interfacer interface {
		Interfacer()
	}
	type WithTime struct {
		T time.Time
		S []int
	}
	//
	s := S{
		A: "S.A",
		B: 12345,
	}
	n := Nested{
		S: S{
			A: "n.S.A",
			B: 54321,
		},
		Next: S{
			A: "n.Next.A",
			B: 9999,
		},
	}
	wt := WithTime{}
	tests := []Test{
		{
			Name:   "&s one",
			V:      &s,
			Fields: []string{"A", "B"},
			Expect: []interface{}{s.A, s.B},
		},
		{
			Name:   "&s two",
			V:      &s,
			Fields: []string{"B", "A"},
			Expect: []interface{}{s.B, s.A},
		},
		// Nesting
		{
			Name:   "n one",
			V:      &n,
			Fields: []string{"S_A", "S_B", "Next_B", "Next_A"},
			Expect: []interface{}{n.S.A, n.S.B, n.Next.B, n.Next.A},
		},
		{
			Name:   "n two",
			V:      &n,
			Fields: []string{"Next_A", "Next_B", "S_B", "S_A"},
			Expect: []interface{}{n.Next.A, n.Next.B, n.S.B, n.S.A},
		},
		// time.Time and S (for default case)
		{
			Name:   "with time",
			V:      &wt,
			Fields: []string{"T", "S"},
			Expect: []interface{}{wt.T, wt.S},
		},
		// Unrecognized field
		{
			Name:   "unknown field",
			V:      &n,
			Fields: []string{"Next_A", "Next_B", "Unknown Field"},
			Error:  set.ErrUnknownField,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			chk := assert.New(t)
			//
			m := set.Mapper{
				TreatAsScalar: set.NewTypeList([]int(nil)),
				Join:          "_",
			}
			p, err := m.Prepare(test.V)
			chk.NoError(err)
			//
			err = p.Plan(test.Fields...)
			chk.ErrorIs(err, test.Error)
			//
			values, err := p.Fields(nil)
			chk.Equal(test.Expect, values)
			if errors.Is(test.Error, set.ErrUnknownField) {
				chk.ErrorIs(err, set.ErrNoPlan)
			} else {
				chk.NoError(err)
				for k := range test.Expect {
					chk.Equal(test.Expect[k], values[k])
				}

			}
		})
	}
	//
	t.Run("readonly", func(t *testing.T) {
		// Unaddressable S results in error
		tests := []Test{
			{
				Name:   "s one",
				V:      s,
				Fields: []string{"A", "B"},
			},
			{
				Name:   "s two",
				V:      s,
				Fields: []string{"B", "A"},
			},
		}
		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				chk := assert.New(t)
				//
				p, err := set.DefaultMapper.Prepare(test.V)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				err = p.Plan(test.Fields...)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				_, err = p.Fields(nil)
				chk.ErrorIs(err, set.ErrReadOnly)
			})
		}
	})
	//
	t.Run("make ptrs", func(t *testing.T) {
		type A struct {
			A string
			B string
		}
		type T struct {
			*A
			Field *A
		}
		//
		chk := assert.New(t)
		//
		data := []T{
			{
				A: &A{
					A: "0.A",
					B: "0.B",
				},
			},
			{
				A: &A{
					A: "1.A",
					B: "1.B",
				},
			},
		}
		// We pass over each element in data twice in this test.
		// The first pass for each element creates the Field *A pointer with zero values
		// but assigns the values from A to Field.
		// The second pass checks these are the correct values.
		for k := 0; k < 4; k++ {
			n := k % 2
			p, err := set.DefaultMapper.Prepare(&data[n])
			chk.NoError(err)
			err = p.Plan([]string{"A_B", "A_A", "Field_A", "Field_B"}...)
			chk.NoError(err)
			values, err := p.Fields(nil)
			chk.NoError(err)
			chk.NotNil(values)
			chk.Equal(4, len(values))
			// embedded + field
			e, f := data[n].A, data[n].Field
			chk.Equal(e.B, values[0])
			chk.Equal(e.A, values[1])
			chk.Equal(f.A, values[2])
			chk.Equal(f.B, values[3])
			//
			chk.Equal(fmt.Sprintf("%v.B", n), values[0])
			chk.Equal(fmt.Sprintf("%v.A", n), values[1])
			//
			if k < 2 {
				chk.Equal("", values[2])
				chk.Equal("", values[3])
				data[k].Field = data[k].A
			} else {
				chk.Equal(fmt.Sprintf("%v.A", n), values[2])
				chk.Equal(fmt.Sprintf("%v.B", n), values[3])
			}
		}
	})
}

func TestPreparedMapping_PlanCorrectlySetsValid(t *testing.T) {
	// ReparedMapping.Plan does not set p.valid=false before iterating the fields and building
	// the plan.  This means if Plan() is called with a valid plan followed by an invalid
	// plan the PreparedMapping instance will still treat the plan as valid.  This tests
	// ensures consecutive calls to Plan() with invalid plans returns the correct error.

	chk := assert.New(t)
	type S struct {
		A int
		B string
	}
	var s S
	p, err := set.DefaultMapper.Prepare(&s)
	chk.NoError(err)

	// Set a valid plan
	err = p.Plan("A", "B")
	chk.NoError(err)

	// Use the plan
	slice, err := p.Fields(nil)
	chk.NoError(err)
	chk.NotNil(slice)
	chk.Len(slice, 2)

	// Set an invalid plan
	err = p.Plan("C")
	chk.ErrorIs(err, set.ErrUnknownField)

	// Use the invalid plan plan
	slice, err = p.Fields(nil)
	chk.ErrorIs(err, set.ErrNoPlan)
	chk.Nil(slice)

	// Set a valid plan
	err = p.Plan("B", "A", "A")
	chk.NoError(err)

	// Use the plan
	slice, err = p.Fields(nil)
	chk.NoError(err)
	chk.NotNil(slice)
	chk.Len(slice, 3)
}

func TestPreparedMapping_ErrNoPlan(t *testing.T) {
	// Mapper.Prepare sets the err field in PreparedMapping to ErrNoPlan.
	// This test covers code blocks in PreparedMapping methods that check for
	// this error and convert it to the pkgerr{} type.
	chk := assert.New(t)
	type S struct {
		A int
		B string
	}
	var s S
	p, err := set.DefaultMapper.Prepare(&s)
	chk.NoError(err)
	slice, err := p.Assignables(nil)
	chk.Nil(slice)
	chk.ErrorIs(err, set.ErrNoPlan)

	chk.ErrorIs(p.Err(), set.ErrNoPlan)

	slice, err = p.Fields(slice)
	chk.Nil(slice)
	chk.ErrorIs(err, set.ErrNoPlan)
}
