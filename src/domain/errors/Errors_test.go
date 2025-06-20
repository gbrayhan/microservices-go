package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppError(t *testing.T) {
	originalError := errors.New("original error")
	appError := NewAppError(originalError, ValidationError)

	assert.NotNil(t, appError)
	assert.Equal(t, originalError, appError.Err)
	assert.Equal(t, ValidationError, appError.Type)
}

func TestNewAppErrorWithType_NotFound(t *testing.T) {
	appError := NewAppErrorWithType(NotFound)

	assert.NotNil(t, appError)
	assert.Equal(t, NotFound, appError.Type)
	assert.Equal(t, "record not found", appError.Error())
}

func TestNewAppErrorWithType_ValidationError(t *testing.T) {
	appError := NewAppErrorWithType(ValidationError)

	assert.NotNil(t, appError)
	assert.Equal(t, ValidationError, appError.Type)
	assert.Equal(t, "validation error", appError.Error())
}

func TestNewAppErrorWithType_ResourceAlreadyExists(t *testing.T) {
	appError := NewAppErrorWithType(ResourceAlreadyExists)

	assert.NotNil(t, appError)
	assert.Equal(t, ResourceAlreadyExists, appError.Type)
	assert.Equal(t, "resource already exists", appError.Error())
}

func TestNewAppErrorWithType_RepositoryError(t *testing.T) {
	appError := NewAppErrorWithType(RepositoryError)

	assert.NotNil(t, appError)
	assert.Equal(t, RepositoryError, appError.Type)
	assert.Equal(t, "error in repository operation", appError.Error())
}

func TestNewAppErrorWithType_NotAuthenticated(t *testing.T) {
	appError := NewAppErrorWithType(NotAuthenticated)

	assert.NotNil(t, appError)
	assert.Equal(t, NotAuthenticated, appError.Type)
	assert.Equal(t, "not Authenticated", appError.Error())
}

func TestNewAppErrorWithType_NotAuthorized(t *testing.T) {
	appError := NewAppErrorWithType(NotAuthorized)

	assert.NotNil(t, appError)
	assert.Equal(t, NotAuthorized, appError.Type)
	assert.Equal(t, "not authorized", appError.Error())
}

func TestNewAppErrorWithType_TokenGeneratorError(t *testing.T) {
	appError := NewAppErrorWithType(TokenGeneratorError)

	assert.NotNil(t, appError)
	assert.Equal(t, TokenGeneratorError, appError.Type)
	assert.Equal(t, "error in token generation", appError.Error())
}

func TestNewAppErrorWithType_UnknownError(t *testing.T) {
	appError := NewAppErrorWithType(UnknownError)

	assert.NotNil(t, appError)
	assert.Equal(t, UnknownError, appError.Type)
	assert.Equal(t, "something went wrong", appError.Error())
}

func TestNewAppErrorWithType_InvalidType(t *testing.T) {
	appError := NewAppErrorWithType("InvalidType")

	assert.NotNil(t, appError)
	assert.Equal(t, ErrorType("InvalidType"), appError.Type)
	assert.Equal(t, "something went wrong", appError.Error()) // Should default to unknown error
}

func TestAppError_Error(t *testing.T) {
	originalError := errors.New("test error message")
	appError := &AppError{
		Err:  originalError,
		Type: ValidationError,
	}

	assert.Equal(t, "test error message", appError.Error())
}

func TestAppErrorToHTTP_NotFound(t *testing.T) {
	appError := NewAppErrorWithType(NotFound)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, "record not found", message)
}

func TestAppErrorToHTTP_ValidationError(t *testing.T) {
	appError := NewAppErrorWithType(ValidationError)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "validation error", message)
}

func TestAppErrorToHTTP_RepositoryError(t *testing.T) {
	appError := NewAppErrorWithType(RepositoryError)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "error in repository operation", message)
}

func TestAppErrorToHTTP_NotAuthenticated(t *testing.T) {
	appError := NewAppErrorWithType(NotAuthenticated)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "not Authenticated", message)
}

func TestAppErrorToHTTP_NotAuthorized(t *testing.T) {
	appError := NewAppErrorWithType(NotAuthorized)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusForbidden, statusCode)
	assert.Equal(t, "not authorized", message)
}

func TestAppErrorToHTTP_UnknownError(t *testing.T) {
	appError := NewAppErrorWithType(UnknownError)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "Internal Server Error", message)
}

func TestAppErrorToHTTP_ResourceAlreadyExists(t *testing.T) {
	appError := NewAppErrorWithType(ResourceAlreadyExists)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "Internal Server Error", message)
}

func TestAppErrorToHTTP_TokenGeneratorError(t *testing.T) {
	appError := NewAppErrorWithType(TokenGeneratorError)
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "Internal Server Error", message)
}

func TestAppErrorToHTTP_CustomError(t *testing.T) {
	appError := NewAppError(errors.New("custom error"), "CustomError")
	statusCode, message := AppErrorToHTTP(appError)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "Internal Server Error", message)
}

func TestErrorTypeConstants(t *testing.T) {
	// Test that all error type constants are defined
	assert.Equal(t, ErrorType("NotFound"), NotFound)
	assert.Equal(t, ErrorType("ValidationError"), ValidationError)
	assert.Equal(t, ErrorType("ResourceAlreadyExists"), ResourceAlreadyExists)
	assert.Equal(t, ErrorType("RepositoryError"), RepositoryError)
	assert.Equal(t, ErrorType("NotAuthenticated"), NotAuthenticated)
	assert.Equal(t, ErrorType("NotAuthorized"), NotAuthorized)
	assert.Equal(t, ErrorType("TokenGeneratorError"), TokenGeneratorError)
	assert.Equal(t, ErrorType("UnknownError"), UnknownError)
}
