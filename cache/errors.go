package cache

import "errors"

var (
	ErrNotFound        = errors.New("key doesn't exist")
	ErrInternalFailure = errors.New("cache internal failure")
)
