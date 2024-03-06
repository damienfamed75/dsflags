package dsflags

import "errors"

// Sentinel errors.
var (
	ErrCannotCastToType = errors.New("cannot natively cast to type")
)
