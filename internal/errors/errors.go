package errors

import (
	"fmt"
	"strconv"
)

type Error struct {
	err         error
	message     string
	recoverable bool
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (recoverable=%s)", e.err.Error(), strconv.FormatBool(e.recoverable))
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) IsRecoverable() bool {
	return e.recoverable
}

func NewError(err error, message string, recoverable bool) Error {
	return Error{
		err:         err,
		message:     message,
		recoverable: recoverable,
	}
}
