package errors

import "errors"

var (
	ErrPostIDInvalid  = errors.New("post ID is invalid")
	ErrPostIDNotFound = errors.New("post ID not found")
	ErrPostIDExpired  = errors.New("post ID expired")
	ErrPostIDDeleted  = errors.New("post ID deleted")
)
