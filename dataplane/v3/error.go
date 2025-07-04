package v3

import (
	"errors"
	"fmt"
)

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Message)
}

// InvalidResponseError represents an error parsing the API response
type InvalidResponseError struct {
	Message string
}

func (e *InvalidResponseError) Error() string {
	return fmt.Sprintf("invalid response: %s", e.Message)
}

// UnauthorizedError represents a 401 Unauthorized error (Auth Failed)
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

// ConflictError represents a 409 Conflict error
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// CommitFailedError represents a transaction commit failure
type CommitFailedError struct {
	Message       string
	TransactionID string
}

func (e *CommitFailedError) Error() string {
	return fmt.Sprintf("commit failed for transaction %s: %s", e.TransactionID, e.Message)
}

// UnknownError represents any other non-2xx status code with HTTP status information
type UnknownError struct {
	Message    string
	StatusCode int
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("unknown error (status %d): %s", e.StatusCode, e.Message)
}

// InternalError represents internal client errors (not from API)
type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Message)
}

// Helper functions to check error types
func IsNotFound(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}

func IsUnauthorized(err error) bool {
	var unauthorizedErr *UnauthorizedError
	return errors.As(err, &unauthorizedErr)
}

func IsConflict(err error) bool {
	var conflictErr *ConflictError
	return errors.As(err, &conflictErr)
}

func IsCommitFailed(err error) bool {
	var commitFailedErr *CommitFailedError
	return errors.As(err, &commitFailedErr)
}

func IsBadRequest(err error) bool {
	var badRequestErr *BadRequestError
	return errors.As(err, &badRequestErr)
}

func IsUnknownError(err error) bool {
	var unknownErr *UnknownError
	return errors.As(err, &unknownErr)
}

// GetHTTPStatusCode returns the HTTP status code from an UnknownError, or 0 if not applicable
func GetHTTPStatusCode(err error) int {
	var unknownErr *UnknownError
	if errors.As(err, &unknownErr) {
		return unknownErr.StatusCode
	}
	return 0
}
