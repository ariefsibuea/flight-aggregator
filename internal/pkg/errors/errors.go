package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error that can be converted to an HTTP response.
type APIError struct {
	StatusCode int
	Message    string
}

// GetErrorCode returns the HTTP status code carried by err if it wraps an APIError, or 500 otherwise.
func GetErrorCode(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code()
	}
	return http.StatusInternalServerError
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) Code() int {
	return e.StatusCode
}

func BadRequestError(msg string) *APIError {
	return &APIError{StatusCode: http.StatusBadRequest, Message: msg}
}

func BadRequestErrorf(format string, args ...any) *APIError {
	return &APIError{StatusCode: http.StatusBadRequest, Message: fmt.Sprintf(format, args...)}
}

func ValidationError(msg string) *APIError {
	return &APIError{StatusCode: http.StatusUnprocessableEntity, Message: msg}
}

func ValidationErrorf(format string, args ...any) *APIError {
	return &APIError{StatusCode: http.StatusUnprocessableEntity, Message: fmt.Sprintf(format, args...)}
}
