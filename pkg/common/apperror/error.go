package apperror

import (
	"fmt"
)

type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	RootCause  error  `json:"-"`
	HTTPStatus int    `json:"-"`
}

// Error() for internal logging
func (e *AppError) Error() string {
	if e.RootCause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.RootCause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// ClientMessage only for client response — never contain RootCause
func (e *AppError) ClientMessage() string {
	return e.Message
}

// New creates an AppError with no underlying cause — use when the error
// originates here (validation, business rule), not from a wrapped call.
func New(code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap creates an AppError around an existing error — use when propagating
// a failure from a lower layer (DB, external API, another package).
func Wrap(err error, code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		RootCause:  err,
	}
}

func (e *AppError) Unwrap() error {
	return e.RootCause
}
