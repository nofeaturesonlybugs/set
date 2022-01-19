package set_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/set"
)

func TestBoundMapping_Err(t *testing.T) {
	type S struct {
		A int
	}
	mapper := &set.Mapper{}
	var s, o S
	var err, errOrig error
	b := mapper.Bind(&s)

	t.Run("err is set", func(t *testing.T) {
		chk := assert.New(t)
		//
		err = b.Set("A", 42)
		chk.NoError(err)
		chk.Equal(err, b.Err())
		err = b.Set("B", "does not exist")
		chk.Error(err)
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
