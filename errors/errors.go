package errors

import (
	"fmt"
	"github.com/subchen/gstack/runtime/stack"
	"io"
)

// _error is an error implementation returned by New and Newf
// that implements its own fmt.Formatter.
type _error struct {
	msg   string
	stack stack.Stack
	cause error
}

// Error implements error interface
func (e *_error) Error() string {
	return e.msg
}

// Cause return e.cause
func (e *_error) Cause() error {
	return e.cause
}

// Format implements fmt.Formatter interface
func (e *_error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			fmt.Fprintf(s, "\n%+v", e.stack)
			if e.cause != nil {
				fmt.Fprintf(s, "Caused by: %+v", e.cause)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

// New returns an error with the supplied message.
func New(message string) error {
	return &_error{message, stack.Callers(3), nil}
}

// Newf returns an error with te supplied message.
func Newf(format string, args ...interface{}) error {
	return &_error{fmt.Sprintf(format, args...), stack.Callers(3), nil}
}

func Wrap(cause error, message string) error {
	return &_error{message, stack.Callers(3), cause}
}

func Wrapf(cause error, format string, args ...interface{}) error {
	return &_error{fmt.Sprintf(format, args), stack.Callers(3), cause}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            error
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		error
		Cause() error
	}

	if err == nil {
		return nil
	}

	for {
		cause, ok := err.(causer)
		if !ok {
			return err
		}
		err = cause.Cause()
		if err == nil {
			return cause
		}
	}
}
