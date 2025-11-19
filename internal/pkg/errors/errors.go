package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeConflict      ErrorType = "CONFLICT"
	ErrorTypeInternal      ErrorType = "INTERNAL"
	ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeTimeout       ErrorType = "TIMEOUT"
	ErrorTypeRateLimit     ErrorType = "RATE_LIMIT"
)

// AppError represents a custom application error
type AppError struct {
	Type    ErrorType              `json:"type"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(errorType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// NotFoundError creates a new not found error
func NotFoundError(resource string, identifier interface{}) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: map[string]interface{}{
			"resource":   resource,
			"identifier": identifier,
		},
	}
}

// ValidationError creates a new validation error
func ValidationError(message string, field string) *AppError {
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: map[string]interface{}{},
	}
	if field != "" {
		err.Details["field"] = field
	}
	return err
}

// ValidationErrors creates a validation error with multiple field errors
func ValidationErrors(fieldErrors map[string]string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: "Validation failed",
		Details: map[string]interface{}{
			"field_errors": fieldErrors,
		},
	}
}

// UnauthorizedError creates a new unauthorized error
func UnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// ForbiddenError creates a new forbidden error
func ForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// ConflictError creates a new conflict error
func ConflictError(resource string, message string) *AppError {
	if message == "" {
		message = fmt.Sprintf("%s already exists", resource)
	}
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Details: map[string]interface{}{
			"resource": resource,
		},
	}
}

// BadRequestError creates a new bad request error
func BadRequestError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// InternalError creates a new internal server error
func InternalError(message string, err error) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// TimeoutError creates a new timeout error
func TimeoutError(operation string) *AppError {
	return &AppError{
		Type:    ErrorTypeTimeout,
		Message: fmt.Sprintf("Operation timed out: %s", operation),
		Details: map[string]interface{}{
			"operation": operation,
		},
	}
}

// RateLimitError creates a new rate limit error
func RateLimitError(retryAfter int) *AppError {
	return &AppError{
		Type:    ErrorTypeRateLimit,
		Message: "Rate limit exceeded",
		Details: map[string]interface{}{
			"retry_after_seconds": retryAfter,
		},
	}
}

// GetHTTPStatus returns the HTTP status code for the error type
func (e *AppError) GetHTTPStatus() int {
	switch e.Type {
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case ErrorTypeRateLimit:
		return http.StatusTooManyRequests
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError attempts to extract an AppError from an error chain
func GetAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsUnauthorized checks if an error is an unauthorized error
func IsUnauthorized(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Type == ErrorTypeUnauthorized
	}
	return false
}

// IsConflict checks if an error is a conflict error
func IsConflict(err error) bool {
	if appErr, ok := GetAppError(err); ok {
		return appErr.Type == ErrorTypeConflict
	}
	return false
}
