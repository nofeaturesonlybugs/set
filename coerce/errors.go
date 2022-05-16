package coerce

import (
	"fmt"
)

// The following errors are returned by this package.
//
// They are typically wrapped and can be checked with errors.Is.
var (
	// ErrInvalid occurs when an attempted coercion is invalid.
	ErrInvalid = fmt.Errorf("coerce: invalid")

	// ErrOverflow occurs when an attempted coercion overflows the destination.
	ErrOverflow = fmt.Errorf("coerce: overflow")

	// ErrUnsupported occurs when a source type for coercion is unsupported.
	ErrUnsupported = fmt.Errorf("coerce: unsupported")
)
