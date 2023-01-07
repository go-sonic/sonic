package xerr

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Status int

const (
	StatusBadRequest          = http.StatusBadRequest
	StatusInternalServerError = http.StatusInternalServerError
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
)

type ErrorType uint

const (
	NoType ErrorType = iota
	BadParam
	NoRecord
	Forbidden
	DB
	Email
)

type customError struct {
	errorType  ErrorType
	cause      error
	httpStatus int
	// msg used to return to the response
	msg    string
	errMsg string
}

func (errorType ErrorType) New(errMsg string, args ...interface{}) *customError {
	err := &customError{errorType: errorType, cause: errors.Errorf(errMsg, args...), httpStatus: -1}
	return err
}

func (errorType ErrorType) Wrapf(err error, errMsg string, args ...interface{}) *customError {
	return &customError{errorType: errorType, cause: errors.Wrapf(err, errMsg, args...), httpStatus: -1}
}

func (errorType ErrorType) Wrap(err error) *customError {
	return &customError{errorType: errorType, cause: errors.Wrap(err, ""), httpStatus: -1}
}

func (ce *customError) Error() string {
	if ce.errMsg != "" {
		if ce.cause != nil {
			return ce.errMsg + " : " + ce.cause.Error()
		} else {
			return ce.errMsg
		}
	}
	if ce.cause != nil {
		return ce.cause.Error()
	}
	return ce.msg
}

func (ce *customError) Cause() error {
	return ce.cause
}

func (ce *customError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", ce.Cause())
			_, _ = io.WriteString(s, ce.msg)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, ce.Error())
	}
}

func WithStatus(err error, status int) *customError {
	//nolint:errorlint
	ee, ok := err.(*customError)
	if ok {
		return &customError{errorType: ee.errorType, cause: ee, httpStatus: status}
	}
	return &customError{errorType: NoType, cause: err, httpStatus: status}
}

func WithMsg(err error, msg string) *customError {
	//nolint:errorlint
	ee, ok := err.(*customError)
	if ok {
		return &customError{errorType: ee.errorType, cause: ee, httpStatus: -1, msg: msg}
	}
	return &customError{errorType: NoType, cause: err, httpStatus: -1, msg: msg}
}

func WithErrMsgf(err error, errMsg string, args ...interface{}) *customError {
	//nolint:errorlint
	ee, ok := err.(*customError)
	if ok {
		return &customError{errorType: ee.errorType, cause: err, httpStatus: -1, errMsg: fmt.Sprintf(errMsg, args...)}
	}
	return &customError{errorType: NoType, cause: err, httpStatus: -1, errMsg: fmt.Sprintf(errMsg, args...)}
}

func (ce *customError) WithErrMsgf(errMsg string, args ...interface{}) *customError {
	return &customError{errorType: ce.errorType, cause: ce, httpStatus: -1, errMsg: fmt.Sprintf(errMsg, args...)}
}

func (ce *customError) WithStatus(status int) *customError {
	return &customError{errorType: ce.errorType, cause: ce, httpStatus: status}
}

func (ce *customError) WithMsg(msg string) *customError {
	return &customError{errorType: ce.errorType, cause: ce, httpStatus: ce.httpStatus, msg: msg}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	//nolint:errorlint
	if ee, ok := err.(*customError); ok {
		return ee.errorType
	}
	return NoType
}

func GetHTTPStatus(err error) int {
	for err != nil {
		//nolint:errorlint
		if e, ok := err.(*customError); ok {
			if e.httpStatus != -1 {
				return e.httpStatus
			} else {
				err = e.cause
			}
		} else {
			break
		}
	}
	return StatusInternalServerError
}

func GetMessage(err error) string {
	for err != nil {
		//nolint:errorlint
		if e, ok := err.(*customError); ok {
			if e.msg != "" {
				return e.msg
			} else {
				err = e.cause
			}
		} else {
			break
		}
	}
	return http.StatusText(http.StatusInternalServerError)
}
