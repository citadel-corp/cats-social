package cat

import "errors"

var (
	ErrCatNotFound      = errors.New("cat not found")
	ErrCatHasMatched    = errors.New("cat has matched befor")
	ErrValidationFailed = errors.New("validation failed")
)
