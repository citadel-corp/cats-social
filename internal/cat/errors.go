package cat

import "errors"

var (
	ErrCatNotFound      = errors.New("cat not found")
	ErrValidationFailed = errors.New("validation failed")
)
