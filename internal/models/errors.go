package models

import "errors"

var (
	ErrTypeAssertionFailed  = errors.New("type assertion failed")
	ErrInvalidRequestParams = errors.New("invalid request params")
)
