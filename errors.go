package set

import (
	"errors"
	"fmt"
)

var (
	// ErrIndexOutOfBounds is returned when an index operation exceeds a bounds check.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrInvalidSlicePtr is returned by NewSlicePtr when the in coming value is not pointer-to-slice.
	ErrInvalidSlicePtr = errors.New("set: expected pointer-to-slice *[]T")

	// ErrNoPlan is returned when a PreparedMapping does not have a valid access plan.
	ErrNoPlan = errors.New("no plan")

	// ErrPlanOutOfBounds is returned when an access to a PreparedMapping exceeds the
	// fields specified by the earlier call to Plan.
	ErrPlanOutOfBounds = errors.New("attempted access extends plan")

	// ErrReadOnly is returned when an incoming argument is expected to be passed by address
	// but is passed by value instead.
	ErrReadOnly = errors.New("read only value")

	// ErrUnknownField is returned by BoundMapping and PreparedMapping when given field
	// has no correlating mapping within the struct hierarchy.
	ErrUnknownField = errors.New("unknown field")

	// ErrUnsupported is returned when an assignment or coercion is incompatible due to the
	// destination and source type(s).
	ErrUnsupported = errors.New("unsupported")
)

// pkgerr is a custom error type to provide more context for sentinal errors.
type pkgerr struct {
	Err      error
	Context  string
	CallSite string
	Hint     string
}

func (e pkgerr) Error() string {
	var c, h string = e.Context, e.Hint
	if c != "" {
		c = ": " + c
	}
	if h != "" {
		h = ": hint=[" + h + "]"
	}
	return fmt.Sprintf("set: %v: %v%v%v", e.CallSite, e.Err.Error(), c, h)
}

func (e pkgerr) Unwrap() error {
	return e.Err
}

// WithCallSite returns a copy of pkgerr with Type and Method set to the arguments.
func (e pkgerr) WithCallSite(callsite string) pkgerr {
	rv := e
	rv.CallSite = callsite
	return rv
}

// WithContext returns a copy of pkgerr with Context set to the argument.
func (e pkgerr) WithContext(context string) pkgerr {
	rv := e
	rv.Context = context
	return rv
}

// WithContextf returns a copy of pkgerr with Context set to the arguments passed to fmt.Sprintf.
func (e pkgerr) WithContextf(context string, args ...interface{}) pkgerr {
	rv := e
	rv.Context = fmt.Sprintf(context, args...)
	return rv
}
