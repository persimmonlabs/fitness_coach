package domain

import "errors"

// Common domain errors
var (
	// ErrNotFound indicates a resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates invalid input data
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates unauthorized access
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates forbidden access
	ErrForbidden = errors.New("forbidden")

	// ErrConflict indicates a conflict with existing data
	ErrConflict = errors.New("conflict with existing data")

	// ErrInternal indicates an internal server error
	ErrInternal = errors.New("internal server error")

	// ErrEmailAlreadyExists indicates email is already registered
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidCredentials indicates invalid login credentials
	ErrInvalidCredentials = errors.New("invalid credentials")
)
