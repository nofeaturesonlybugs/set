package set_test

import (
	"fmt"
	"testing"

	"github.com/nofeaturesonlybugs/set"
)

func Benchmark_PerformanceGuide(b *testing.B) {
	mapper := &set.Mapper{
		Join: ".",
	}
	//
	type Single struct {
		A int
	}
	type Double struct {
		Single Single
	}
	type Triple struct {
		Double Double
	}
	type Quad struct {
		Triple Triple
	}
	//
	type Test struct {
		Name  string
		Iters int
	}
	type Type struct {
		Name       string
		Value      interface{}
		MappedName string
	}
	tests := []Test{
		{
			Name:  "1x",
			Iters: 1,
		},
		{
			Name:  "3x",
			Iters: 3,
		},
		{
			Name:  "10x",
			Iters: 10,
		},
	}
	types := []Type{
		{Name: "Single", Value: &Single{}, MappedName: "A"},
		{Name: "Double", Value: &Double{}, MappedName: "Single.A"},
		{Name: "Triple", Value: &Triple{}, MappedName: "Double.Single.A"},
		{Name: "Quad", Value: &Quad{}, MappedName: "Triple.Double.Single.A"},
	}
	//
	const (
		IterWidth = 3
		TestWidth = 9
	)
	nameAdjust := func(s string, n int) string {
		return fmt.Sprintf("%[1]*[2]v", n, s)
	}
	//
	for _, typ := range types {
		b.Run(typ.Name, func(b *testing.B) {
			for _, test := range tests {
				b.Run(nameAdjust(fmt.Sprintf("%vx", test.Iters), IterWidth), func(b *testing.B) {
					b.Run(nameAdjust("go", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								switch v := typ.Value.(type) {
								case *Single:
									v.A = 42
								case *Double:
									v.Single.A = 42
								case *Triple:
									v.Double.Single.A = 42
								case *Quad:
									v.Triple.Double.Single.A = 42
								default:
									b.Fatalf(fmt.Sprintf("unhandled type %T", typ.Value))
								}
							}
						}
					})
					b.Run(nameAdjust("go ptr", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								switch v := typ.Value.(type) {
								case *Single:
									_ = &v.A
								case *Double:
									_ = &v.Single.A
								case *Triple:
									_ = &v.Double.Single.A
								case *Quad:
									_ = &v.Triple.Double.Single.A
								default:
									b.Fatalf(fmt.Sprintf("unhandled type %T", typ.Value))
								}
							}
						}
					})
					b.Run(nameAdjust("bound", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						var err error
						bound := mapper.Bind(typ.Value)
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								err = bound.Set(typ.MappedName, 42)
								if err != nil {
									b.Fatalf("BoundMapping.Set %v %v", typ.MappedName, err)
								}
							}
						}
					})
					b.Run(nameAdjust("bound ptr", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						var err error
						ptr := make([]interface{}, 1)
						fields := []string{typ.MappedName}
						bound := mapper.Bind(typ.Value)
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								_, err = bound.Assignables(fields, ptr)
								if err != nil {
									b.Fatalf("BoundMapping.Assignables %v %v", typ.MappedName, err)
								}
							}
						}
					})
				})
			}
		})
	}
}

func Benchmark_PerformanceGuideDerefs(b *testing.B) {
	mapper := &set.Mapper{
		Join: ".",
	}
	//
	type Single struct {
		A *int
	}
	type Double struct {
		Single *Single
	}
	type Triple struct {
		Double *Double
	}
	type Quad struct {
		Triple *Triple
	}
	var i int
	single := Single{A: &i}
	double := Double{Single: &single}
	triple := Triple{Double: &double}
	quad := Quad{Triple: &triple}
	//
	type Test struct {
		Name  string
		Iters int
	}
	type Type struct {
		Name       string
		Value      interface{}
		MappedName string
	}
	tests := []Test{
		{
			Name:  "1x",
			Iters: 1,
		},
		{
			Name:  "3x",
			Iters: 3,
		},
		{
			Name:  "10x",
			Iters: 10,
		},
	}
	types := []Type{
		{Name: "Single", Value: &single, MappedName: "A"},
		{Name: "Double", Value: &double, MappedName: "Single.A"},
		{Name: "Triple", Value: &triple, MappedName: "Double.Single.A"},
		{Name: "Quad", Value: &quad, MappedName: "Triple.Double.Single.A"},
	}
	//
	const (
		IterWidth = 3
		TestWidth = 9
	)
	nameAdjust := func(s string, n int) string {
		return fmt.Sprintf("%[1]*[2]v", n, s)
	}
	//
	for _, typ := range types {
		b.Run(typ.Name, func(b *testing.B) {
			for _, test := range tests {
				b.Run(nameAdjust(fmt.Sprintf("%vx", test.Iters), IterWidth), func(b *testing.B) {
					b.Run(nameAdjust("go", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								switch v := typ.Value.(type) {
								case *Single:
									*v.A = 42
								case *Double:
									*v.Single.A = 42
								case *Triple:
									*v.Double.Single.A = 42
								case *Quad:
									*v.Triple.Double.Single.A = 42
								default:
									b.Fatalf(fmt.Sprintf("unhandled type %T", typ.Value))
								}
							}
						}
					})
					b.Run(nameAdjust("go ptr", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								switch v := typ.Value.(type) {
								case *Single:
									_ = &v.A
								case *Double:
									_ = &v.Single.A
								case *Triple:
									_ = &v.Double.Single.A
								case *Quad:
									_ = &v.Triple.Double.Single.A
								default:
									b.Fatalf(fmt.Sprintf("unhandled type %T", typ.Value))
								}
							}
						}
					})
					b.Run(nameAdjust("bound", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						var err error
						bound := mapper.Bind(typ.Value)
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								err = bound.Set(typ.MappedName, 42)
								if err != nil {
									b.Fatalf("BoundMapping.Set %v %v", typ.MappedName, err)
								}
							}
						}
					})
					b.Run(nameAdjust("bound ptr", TestWidth), func(b *testing.B) {
						b.ResetTimer()
						var err error
						ptr := make([]interface{}, 1)
						fields := []string{typ.MappedName}
						bound := mapper.Bind(typ.Value)
						for n := 0; n < b.N; n++ {
							for k := 0; k < test.Iters; k++ {
								_, err = bound.Assignables(fields, ptr)
								if err != nil {
									b.Fatalf("BoundMapping.Assignables %v %v", typ.MappedName, err)
								}
							}
						}
					})
				})
			}
		})
	}
}
