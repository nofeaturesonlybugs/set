package path_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set/path"
)

func TestTree(t *testing.T) {
	type S struct {
		Str string
		Num int
	}
	type Nested struct {
		S
		Next  *S
		Other S
	}
	t.Run("string", func(t *testing.T) {
		n := &Nested{}
		tree := path.Stat(n)
		_ = tree.String()
	})
	t.Run("slice", func(t *testing.T) {
		tree := path.Stat(Nested{})
		tree.Slice()
	})
	t.Run("empty", func(t *testing.T) {
		type Empty struct{}
		tree := path.Stat(Empty{})
		tree.Slice()
	})
}
