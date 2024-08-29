package repos

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInternalFailure = errors.New("internal failure")
)
