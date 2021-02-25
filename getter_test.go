package set_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
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
