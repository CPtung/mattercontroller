package server

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError interface {
	Error() string
	Status() int
	Code() int
}

type httpError struct {
	status  int
	message string
	code    int
}

func New(s string) error {
	return errors.New(s)
}

func (a httpError) Status() int {
	return a.status
}

func (a httpError) Error() string {
	return a.message
}

func (a httpError) Code() int {
	return a.code
}

// HTTPErrorItemNotFound is a HttpError
func HTTPErrorItemNotFound(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusNotFound,
		fmt.Sprintf(format, msg...),
		http.StatusNotFound,
	}
}

// HTTPErrorBadRequest is a HttpError
func HTTPErrorBadRequest(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusBadRequest,
		fmt.Sprintf(format, msg...),
		http.StatusBadRequest,
	}
}

// HTTPErrorCodeBadRequest is a HttpError
func HTTPErrorCodeBadRequest(code int, format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusBadRequest,
		fmt.Sprintf(format, msg...),
		code,
	}
}

// HTTPErrorInternal is a HttpError
func HTTPErrorInternal(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusInternalServerError,
		fmt.Sprintf(format, msg...),
		http.StatusInternalServerError,
	}
}

// HTTPErrorUnauthorized is a HttpError
func HTTPErrorUnauthorized(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusUnauthorized,
		fmt.Sprintf(format, msg...),
		http.StatusUnauthorized,
	}
}

// HTTPErrorMethodNotAllowed is a HttpError
func HTTPErrorMethodNotAllowed(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusMethodNotAllowed,
		fmt.Sprintf(format, msg...),
		http.StatusMethodNotAllowed,
	}
}

// HTTPErrorTimeout is a HttpError
func HTTPErrorTimeout(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusRequestTimeout,
		fmt.Sprintf(format, msg...),
		http.StatusRequestTimeout,
	}
}

// HTTPErrorForbidden is a HttpError
func HTTPErrorForbidden(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusForbidden,
		fmt.Sprintf(format, msg...),
		http.StatusForbidden,
	}
}

// HTTPErrorConflict is a HttpError
func HTTPErrorConflict(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusConflict,
		fmt.Sprintf(format, msg...),
		http.StatusConflict,
	}
}

// HTTPError is a HttpError
func HTTPError(status, code int, format string, msg ...interface{}) HttpError {
	return httpError{
		status,
		fmt.Sprintf(format, msg...),
		code,
	}
}

// HTTPErrorServiceUnavailable is a HttpError
func HTTPErrorServiceUnavailable(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusServiceUnavailable,
		fmt.Sprintf(format, msg...),
		http.StatusServiceUnavailable,
	}
}

// HTTPErrorRequestEntityTooLarge is a HttpError
func HTTPErrorRequestEntityTooLarge(format string, msg ...interface{}) HttpError {
	return httpError{
		http.StatusRequestEntityTooLarge,
		fmt.Sprintf(format, msg...),
		http.StatusRequestEntityTooLarge,
	}
}
