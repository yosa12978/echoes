package types

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInternalFailure  = errors.New("internal failure")
	ErrBadRequest       = errors.New("bad request")
	ErrValidationFailed = errors.New("validation failed")
)

func NewErrNotFound(err error) error {
	return errors.Join(err, ErrNotFound)
}

func NewErrInternalFailure(err error) error {
	return errors.Join(err, ErrInternalFailure)
}

func NewErrBadRequest(err error) error {
	return errors.Join(err, ErrBadRequest)
}

func NewErrValidationFailed(err error) error {
	return errors.Join(err, ErrValidationFailed)
}
