package catmatch

import "errors"

var (
	ErrCatMatchNotFound      = errors.New("cat match not found")
	ErrCatMatchNoLongerValid = errors.New("cat match no longer valid")
	ErrCatHasMatched         = errors.New("cat has matched before")
	ErrCatSameSex            = errors.New("cat has same sex")
	ErrCatSameUser           = errors.New("cat has same user")
	ErrValidationFailed      = errors.New("validation failed")
)
