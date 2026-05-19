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

func NewError(err error, message string, recoverable bool, args ...any) Error {
	if err == nil {
		err = fmt.Errorf(message, args...)
	} else if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return Error{
		err:         err,
		message:     message,
		recoverable: recoverable,
	}
}

func NewNotImplementedError(feature string) Error {
	return NewError(nil, "feature '%s' not implemented", false, feature)
}

func NewConfigurationError(message string, args ...any) Error {
	return NewError(nil, "configuration error: "+message, false, args...)
}
