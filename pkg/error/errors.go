package error

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// UnknownError error
	UnknownError ErrorType = iota
	DBError
	ValidationError
	NotFound
)

type AppError struct {
	ErrorType     ErrorType
	OriginalError error
	Context       errorContext
}

type errorContext struct {
	Field   string `json:"field"`
	Message string `json:"msg"`
}

func (e AppError) MarshalJSON() ([]byte, error) {
	jsonErr := struct {
		ErrorType string       `json:"type"`
		Cause     string       `json:"cause"`
		Context   errorContext `json:"context,omitempty"`
	}{
		ErrorType: e.ErrorType.String(),
		Cause:     e.OriginalError.Error(),
		Context:   e.Context,
	}

	return json.Marshal(jsonErr)
}

func (e AppError) error() error {
	return e.OriginalError
}

// AddErrorContext adds a Context to an error
func (e AppError) AddContext(field, message string) AppError {
	context := errorContext{Field: field, Message: message}
	e.Context = context
	return e
}

func (errorType ErrorType) String() string {
	return [...]string{"UnknownError", "DBError", "ValidationError", "NotFound"}[errorType]
}

// New creates a new AppError
func (errorType ErrorType) New(msg string) AppError {
	return AppError{ErrorType: errorType, OriginalError: errors.New(msg)}
}

// New creates a new AppError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) AppError {
	return AppError{ErrorType: errorType, OriginalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) AppError {
	return errorType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) AppError {
	return AppError{ErrorType: errorType, OriginalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the mssage of a AppError
func (e AppError) Error() string {
	return e.OriginalError.Error()
}

// New creates a no type error
func New(msg string) AppError {
	return AppError{ErrorType: UnknownError, OriginalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) AppError {
	return AppError{ErrorType: UnknownError, OriginalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(AppError); ok {
		return AppError{
			ErrorType:     customErr.ErrorType,
			OriginalError: wrappedError,
			Context:       customErr.Context,
		}
	}

	return AppError{ErrorType: UnknownError, OriginalError: wrappedError}
}

// AddErrorContext adds a Context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(AppError); ok {
		return AppError{ErrorType: customErr.ErrorType, OriginalError: customErr.OriginalError, Context: context}
	}

	return AppError{ErrorType: UnknownError, OriginalError: err, Context: context}
}

// GetErrorContext returns the error Context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(AppError); ok || customErr.Context != emptyContext {

		return map[string]string{"field": customErr.Context.Field, "message": customErr.Context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(AppError); ok {
		return customErr.ErrorType
	}

	return UnknownError
}
