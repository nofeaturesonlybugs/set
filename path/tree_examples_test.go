package path_test

import (
	"fmt"
	"time"

	"github.com/nofeaturesonlybugs/set/path"
)

func ExampleTree() {
	// The following fields become leaves:
	//	- T
	//	- Str
	//	- Int
	//	- X
	//	- Y
	//
	// There is only one branch:
	//	- A (embedded in Foo)

	type A struct {
		T   time.Time
		Str string
		Int int
	}
	type Foo struct {
		X float64
		Y float64
		A
	}

	tree := path.Stat(Foo{})
	fmt.Println(tree.StringIndent("    "))

	// Output: Branches
	//     A = A 2 16 Type=path_test.A Pathway[A][2] Parent[] Offsets= +16 ∴ path_test.A
	// Leaves
	//     X = X 0 0 Type=float64 Pathway[X][0] Parent[] Offsets= +0 ∴ float64
	//     Y = Y 1 8 Type=float64 Pathway[Y][1] Parent[] Offsets= +8 ∴ float64
	//         A.T = T 0 0 Type=time.Time Pathway[A.T][2 0] Parent[A] Offsets= +16 ∴ time.Time
	//         A.Str = Str 1 24 Type=string Pathway[A.Str][2 1] Parent[A] Offsets= +40 ∴ string
	//         A.Int = Int 2 40 Type=int Pathway[A.Int][2 2] Parent[A] Offsets= +56 ∴ int
}
