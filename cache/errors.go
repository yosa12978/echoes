package cache

import "errors"

var (
	ErrNotFound        = errors.New("key doesn't exist")
	ErrInternalFailure = errors.New("cache internal failure")
)

func newNotFound(err error) error {
	return errors.Join(err, ErrNotFound)
}

func newInternalFailure(err error) error {
	return errors.Join(err, ErrInternalFailure)
}
