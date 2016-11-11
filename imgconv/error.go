package imgconv

import "fmt"

type ErrorType uint

const (
	ErrUnsupportedFormat ErrorType = iota
)

type Error struct {
	Type    ErrorType
	Message string
}

// Error returns the error's message
func (e *Error) Error() string {
	return e.Message
}

func newError(tp ErrorType, message string) *Error {
	return &Error{
		Type:    tp,
		Message: message,
	}
}

func newErrorf(tp ErrorType, format string, args ...interface{}) *Error {
	return newError(tp, fmt.Sprintf(format, args...))
}
