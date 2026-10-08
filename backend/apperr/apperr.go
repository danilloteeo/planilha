package apperr

import "errors"

type Error struct {
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(message string, err error) *Error {
	return &Error{Message: message, Err: err}
}

func Message(err error, fallback string) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return fallback
}
