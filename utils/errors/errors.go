package errors

import "net/http"

// RestError struct. Has a status code and a custom message
type RestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e RestError) Error() string {
	return e.Message
}

// NewInternalServerError returns error with status code 500
func NewInternalServerError(message string) *RestError {
	return &RestError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

// NewConflictError returns error with status code 409
func NewConflictError(message string) *RestError {
	return &RestError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// NewNotFoundError returns error with status code 404
func NewNotFoundError(message string) *RestError {
	return &RestError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

// NewBadRequestError returns error with status code 400
func NewBadRequestError(message string) *RestError {
	return &RestError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// NewUnauthorizedError returns error with status code 401
func NewUnauthorizedError(message string) *RestError {
	return &RestError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
