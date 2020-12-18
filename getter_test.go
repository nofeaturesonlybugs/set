package set_test

import (
	"testing"

	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/set/assert"
)

func TestMapGetter_codeCoverage(t *testing.T) {
	chk := assert.New(t)
	//
	{
		chk.NotNil(set.MapGetter(42))
	}
	{
		g := set.MapGetter(map[int]bool{})
		chk.NotNil(g)
		chk.Nil(g.Get("foo"))
	}
}
