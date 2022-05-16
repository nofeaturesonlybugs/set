package set_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestBoundMapping_Assignables(t *testing.T) {
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
			b, err := set.DefaultMapper.Bind(test.V)
			chk.NoError(err)
			//
			ptrs, err := b.Assignables(test.Fields, nil)
			chk.ErrorIs(err, test.Error)
			if test.Error == nil { // Do not expect error
				chk.NoError(err)
				chk.Equal(test.Expect, ptrs)
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
				b, err := set.DefaultMapper.Bind(test.V)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				_, err = b.Assignables(test.Fields, nil)
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
			bound, err := set.DefaultMapper.Bind(&data[n])
			chk.NoError(err)
			ptrs, err := bound.Assignables([]string{"A_B", "A_A", "Field_A", "Field_B"}, nil)
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

func TestBoundMapping_Copy(t *testing.T) {
	chk := assert.New(t)
	//
	type S struct {
		A string
		B int
	}
	var first, copied S
	var err error

	b, err := set.DefaultMapper.Bind(&first)
	chk.NoError(err)
	err = b.Set("A", "Hello")
	chk.NoError(err)
	chk.Nil(b.Err())
	err = b.Set("B", 42)
	chk.NoError(err)
	chk.Nil(b.Err())

	cp := b.Copy()
	cp.Rebind(&copied)
	err = cp.Set("A", "Copied!")
	chk.NoError(err)
	chk.Nil(cp.Err())
	err = cp.Set("B", 100)
	chk.NoError(err)
	chk.Nil(cp.Err())

	chk.Equal("Hello", first.A)
	chk.Equal(42, first.B)
	chk.Equal("Copied!", copied.A)
	chk.Equal(100, copied.B)
}

func TestBoundMapping_Err(t *testing.T) {
	chk := assert.New(t)
	//
	type S struct {
		A int
	}
	mapper := &set.Mapper{}
	var s, o S
	var err, errOrig error
	b, err := mapper.Bind(&s)
	chk.NoError(err)

	t.Run("err is set", func(t *testing.T) {
		chk := assert.New(t)
		//
		err = b.Set("A", 42)
		chk.NoError(err)
		chk.Equal(err, b.Err())
		err = b.Set("B", "does not exist")
		chk.ErrorIs(err, set.ErrUnknownField)
		chk.Equal(err, b.Err())
		errOrig = err
		err = b.Set("A", 48)
		chk.NoError(err)
		chk.Equal(errOrig, b.Err())
	})

	t.Run("rebind clears error", func(t *testing.T) {
		chk := assert.New(t)
		//
		b.Rebind(&o)
		chk.Nil(b.Err())
	})

}

func TestBoundMapping_Fields(t *testing.T) {
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
	type WithTime struct {
		T time.Time
		S []int
	}
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
			b, err := m.Bind(test.V)
			chk.NoError(err)
			//
			values, err := b.Fields(test.Fields, nil)
			chk.ErrorIs(err, test.Error)
			if test.Error == nil { // Do not expect error
				chk.Equal(test.Expect, values)
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
				b, err := set.DefaultMapper.Bind(test.V)
				chk.ErrorIs(err, set.ErrReadOnly)
				//
				_, err = b.Fields(test.Fields, nil)
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
			bound, err := set.DefaultMapper.Bind(&data[n])
			chk.NoError(err)
			values, err := bound.Fields([]string{"A_B", "A_A", "Field_A", "Field_B"}, nil)
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
