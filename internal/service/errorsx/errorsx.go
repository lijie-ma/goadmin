package errorsx

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid parameter")
	ErrReqired  = errors.New("required parameter")
)
