package path_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set/path"
)

func TestPath(t *testing.T) {
	t.Run("empty indeces", func(t *testing.T) {
		var i path.PathIndeces
		_ = i.String()
	})
	t.Run("empty offsets", func(t *testing.T) {
		var o path.PathOffsets
		_ = o.String()
	})
}

func TestPath_TypeAll(t *testing.T) {
	type This struct {
		A int
		B int
	}
	type That struct {
		M string
		N string
	}
	type More struct {
		Beep string
		Bip  int
		Boop string
	}
	type Other struct {
		X float32
		Y float32
		Z *float32
		More
	}
	type All struct {
		This
		T That
		*Other

		ThreeP ***Other
	}
	//
	// One Other with all pointers filled in.
	pOther := &Other{
		X: math.Pi,
		Y: math.Pi,
		Z: new(float32),
		More: More{
			Beep: "ThreeP.BeepBeepBeep",
			Bip:  1000,
			Boop: "ThreeP.BoopBoopBoop",
		},
	}
	ppOther := &pOther
	pppOther := &ppOther
	v := All{
		This: This{
			A: 0,
			B: 0,
		},
		T: That{
			M: "MMMMMM",
			N: "NNNNNN",
		},
		Other: &Other{
			X: math.Pi,
			Y: math.Pi,
			Z: new(float32),
			More: More{
				Beep: "BeepBeepBeep",
				Bip:  9000,
				Boop: "BoopBoopBoop",
			},
		},
		ThreeP: pppOther,
	}
	*pOther.Z = math.Pi
	*v.Other.Z = math.Pi
	//
	// And another All that is empty.
	empty := &All{}
	//
	all := map[string]path.Path{}
	{
		vtree := path.Stat(v)
		for key, leaf := range vtree.Leaves {
			all[key] = leaf
		}
		for key, br := range vtree.Branches {
			all[key] = br
		}
	}
	t.Run("reflect", func(t *testing.T) {
		var origin, p reflect.Value
		origin = reflect.ValueOf(v)
		chk := assert.New(t)
		//
		p = all["This"].Value(origin)
		chk.Equal(reflect.TypeOf(This{}), p.Type())
		chk.Equal(v.This, p.Interface())
		p = all["This.A"].Value(origin)
		chk.Equal(0, p.Interface())
		p = all["This.B"].Value(origin)
		chk.Equal(0, p.Interface())

		p = all["T"].Value(origin)
		chk.Equal(reflect.TypeOf(That{}), p.Type())
		chk.Equal(v.T, p.Interface())
		p = all["T.M"].Value(origin)
		chk.Equal("MMMMMM", p.Interface())
		p = all["T.N"].Value(origin)
		chk.Equal("NNNNNN", p.Interface())

		p = all["Other"].Value(origin)
		chk.Equal(reflect.TypeOf(Other{}), p.Type())
		chk.Equal(*v.Other, p.Interface())
		p = all["Other.X"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["Other.Y"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["Other.Z"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["Other.More"].Value(origin)
		chk.Equal(reflect.TypeOf(More{}), p.Type())
		chk.Equal(v.Other.More, p.Interface())
		p = all["Other.More.Beep"].Value(origin)
		chk.Equal("BeepBeepBeep", p.Interface())
		p = all["Other.More.Bip"].Value(origin)
		chk.Equal(9000, p.Interface())
		p = all["Other.More.Boop"].Value(origin)
		chk.Equal("BoopBoopBoop", p.Interface())

		p = all["ThreeP"].Value(origin)
		chk.Equal(reflect.TypeOf(Other{}), p.Type())
		chk.Equal(***v.ThreeP, p.Interface())
		p = all["ThreeP.X"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["ThreeP.Y"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["ThreeP.Z"].Value(origin)
		chk.InDelta(math.Pi, p.Interface(), 0.0001)
		p = all["ThreeP.More"].Value(origin)
		chk.Equal(reflect.TypeOf(More{}), p.Type())
		chk.Equal((***v.ThreeP).More, p.Interface())
		p = all["ThreeP.More.Beep"].Value(origin)
		chk.Equal("ThreeP.BeepBeepBeep", p.Interface())
		p = all["ThreeP.More.Bip"].Value(origin)
		chk.Equal(1000, p.Interface())
		p = all["ThreeP.More.Boop"].Value(origin)
		chk.Equal("ThreeP.BoopBoopBoop", p.Interface())
	})
	t.Run("reflect empty", func(t *testing.T) {
		var origin, p reflect.Value
		origin = reflect.ValueOf(empty)
		chk := assert.New(t)
		//
		p = all["This"].Value(origin)
		chk.Equal(reflect.TypeOf(This{}), p.Type())
		chk.Equal(empty.This, p.Interface())
		p = all["This.A"].Value(origin)
		chk.Equal(0, p.Interface())
		p = all["This.B"].Value(origin)
		chk.Equal(0, p.Interface())

		p = all["T"].Value(origin)
		chk.Equal(reflect.TypeOf(That{}), p.Type())
		chk.Equal(empty.T, p.Interface())
		p = all["T.M"].Value(origin)
		chk.Equal("", p.Interface())
		p = all["T.N"].Value(origin)
		chk.Equal("", p.Interface())

		p = all["Other"].Value(origin)
		chk.Equal(reflect.TypeOf(Other{}), p.Type())
		chk.Equal(*empty.Other, p.Interface())
		p = all["Other.X"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["Other.Y"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["Other.Z"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["Other.More"].Value(origin)
		chk.Equal(reflect.TypeOf(More{}), p.Type())
		chk.Equal(empty.Other.More, p.Interface())
		p = all["Other.More.Beep"].Value(origin)
		chk.Equal("", p.Interface())
		p = all["Other.More.Bip"].Value(origin)
		chk.Equal(0, p.Interface())
		p = all["Other.More.Boop"].Value(origin)
		chk.Equal("", p.Interface())

		p = all["ThreeP"].Value(origin)
		chk.Equal(reflect.TypeOf(Other{}), p.Type())
		chk.Equal(***empty.ThreeP, p.Interface())
		p = all["ThreeP.X"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["ThreeP.Y"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["ThreeP.Z"].Value(origin)
		chk.Equal(float32(0), p.Interface())
		p = all["ThreeP.More"].Value(origin)
		chk.Equal(reflect.TypeOf(More{}), p.Type())
		chk.Equal((***empty.ThreeP).More, p.Interface())
		p = all["ThreeP.More.Beep"].Value(origin)
		chk.Equal("", p.Interface())
		p = all["ThreeP.More.Bip"].Value(origin)
		chk.Equal(0, p.Interface())
		p = all["ThreeP.More.Boop"].Value(origin)
		chk.Equal("", p.Interface())
	})
}

func TestPath_TypeVendor(t *testing.T) {
	type C struct {
		Id           int
		CreatedTime  time.Time
		ModifiedTime time.Time
	}

	type Person struct {
		C
		First string
		Last  string
	}

	type Vendor struct {
		C
		Name    string
		Contact Person
	}

	v := Vendor{
		C: C{
			Id:           10,
			CreatedTime:  time.Now().Add(-10 * time.Hour),
			ModifiedTime: time.Now().Add(-10 * time.Hour),
		},
		Name: "Widget Co.",
		Contact: Person{
			C: C{
				Id:           20,
				CreatedTime:  time.Now().Add(-20 * time.Hour),
				ModifiedTime: time.Now().Add(-5 * time.Hour),
			},
			First: "Bob",
			Last:  "Smith",
		},
	}
	all := map[string]path.Path{}
	{
		vtree := path.Stat(v)
		for key, leaf := range vtree.Leaves {
			all[key] = leaf
		}
		for key, br := range vtree.Branches {
			all[key] = br
		}
	}
	//
	t.Run("reflect", func(t *testing.T) {
		var origin, p reflect.Value
		origin = reflect.ValueOf(v)
		chk := assert.New(t)
		//
		p = all["C"].Value(origin)
		chk.Equal(reflect.TypeOf(C{}), p.Type())
		chk.Equal(v.C, p.Interface())
		p = all["C.Id"].Value(origin)
		chk.Equal(10, p.Interface())
		p = all["C.CreatedTime"].Value(origin)
		chk.Equal(v.CreatedTime, p.Interface())
		p = all["C.ModifiedTime"].Value(origin)
		chk.Equal(v.ModifiedTime, p.Interface())

		p = all["Name"].Value(origin)
		chk.Equal("Widget Co.", p.Interface())

		p = all["Contact"].Value(origin)
		chk.Equal(reflect.TypeOf(Person{}), p.Type())
		chk.Equal(v.Contact, p.Interface())
		p = all["Contact.C"].Value(origin)
		chk.Equal(reflect.TypeOf(C{}), p.Type())
		chk.Equal(v.Contact.C, p.Interface())
		p = all["Contact.C.Id"].Value(origin)
		chk.Equal(20, p.Interface())
		p = all["Contact.C.CreatedTime"].Value(origin)
		chk.Equal(v.Contact.CreatedTime, p.Interface())
		p = all["Contact.C.ModifiedTime"].Value(origin)
		chk.Equal(v.Contact.ModifiedTime, p.Interface())
		p = all["Contact.First"].Value(origin)
		chk.Equal("Bob", p.Interface())
		p = all["Contact.Last"].Value(origin)
		chk.Equal("Smith", p.Interface())
	})
}
