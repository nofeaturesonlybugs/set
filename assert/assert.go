package assert

import (
	"reflect"
	"testing"

	"github.com/nofeaturesonlybugs/errors"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) *Assert {
	return &Assert{t}
}

func (me *Assert) Equal(expected, actual interface{}) {
	if expected != actual {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed Equal(): %T( %v ) != %T( %v ) @ %v:%v", expected, expected, actual, actual, stack.Line, stack.Function)
	}
}

func (me *Assert) NotEqual(expected, actual interface{}) {
	if expected == actual {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed NotEqual(): %T( %v ) == %T( %v ) @ %v:%v", expected, expected, actual, actual, stack.Line, stack.Function)
	}
}

func (me *Assert) Error(err error) {
	if err == nil {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed Error() @ %v:%v", stack.Line, stack.Function)
	}
}

func (me *Assert) NoError(err error) {
	if err != nil {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed NoError() @ %v:%v with %v", stack.Line, stack.Function, err.Error())
	}
}

func (me *Assert) InDelta(expected interface{}, actual interface{}, delta interface{}) {
	var e, a, d float64
	switch t := expected.(type) {
	case float32:
		e = float64(t)
	case float64:
		e = t
	default:
		me.t.Errorf("Failed InDelta(): expected is not a float %T( %v )", expected, expected)
	}
	switch t := actual.(type) {
	case float32:
		a = float64(t)
	case float64:
		a = t
	default:
		me.t.Errorf("Failed InDelta(): actual is not a float %T( %v )", actual, actual)
	}
	switch t := delta.(type) {
	case float32:
		d = float64(t)
	case float64:
		d = t
	default:
		me.t.Errorf("Failed InDelta(): delta is not a float %T( %v )", delta, delta)
	}
	if a < (e-d) || a > (e+d) {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed InDelta(): actual not within expected += delta; %v <> %v += %v @ %v:%v", a, e, d, stack.Line, stack.Function)
	}
}

func (me *Assert) NotNil(v interface{}) {
	if v == nil {
		stack := errors.Stack()[1]
		me.t.Errorf("Failed NotNil() @ %v:%v", stack.Line, stack.Function)
	}
}

func (me *Assert) Nil(v interface{}) {
	if v == nil {
		return
	} else {
		value := reflect.ValueOf(v)
		k := value.Kind()
		nillable := k == reflect.Chan || k == reflect.Func || k == reflect.Interface || k == reflect.Map || k == reflect.Ptr || k == reflect.Slice
		if nillable && value.IsNil() {
			return
		}
		stack := errors.Stack()[1]
		me.t.Errorf("Failed Nil() @ %v:%v", stack.Line, stack.Function)
	}
}
