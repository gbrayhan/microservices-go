package errors

import (
	"errors"
	"testing"
)

func TestNewAppError(t *testing.T) {
	originalErr := errors.New("test error")
	appErr := NewAppError(originalErr, ErrorType("TestError"))

	if appErr.Err != originalErr {
		t.Errorf("Expected error to be %v, got %v", originalErr, appErr.Err)
	}

	if appErr.Type != ErrorType("TestError") {
		t.Errorf("Expected type to be 'TestError', got %s", appErr.Type)
	}
}

func TestNewAppErrorWithType_NotFound(t *testing.T) {
	appErr := NewAppErrorWithType(NotFound)

	if appErr.Type != NotFound {
		t.Errorf("Expected type to be %s, got %s", NotFound, appErr.Type)
	}

	if appErr.Error() != string(notFoundMessage) {
		t.Errorf("Expected error message to be %s, got %s", notFoundMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_ValidationError(t *testing.T) {
	appErr := NewAppErrorWithType(ValidationError)

	if appErr.Type != ValidationError {
		t.Errorf("Expected type to be %s, got %s", ValidationError, appErr.Type)
	}

	if appErr.Error() != string(validationErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", validationErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_ResourceAlreadyExists(t *testing.T) {
	appErr := NewAppErrorWithType(ResourceAlreadyExists)

	if appErr.Type != ResourceAlreadyExists {
		t.Errorf("Expected type to be %s, got %s", ResourceAlreadyExists, appErr.Type)
	}

	if appErr.Error() != string(alreadyExistsErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", alreadyExistsErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_RepositoryError(t *testing.T) {
	appErr := NewAppErrorWithType(RepositoryError)

	if appErr.Type != RepositoryError {
		t.Errorf("Expected type to be %s, got %s", RepositoryError, appErr.Type)
	}

	if appErr.Error() != string(repositoryErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", repositoryErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_NotAuthenticated(t *testing.T) {
	appErr := NewAppErrorWithType(NotAuthenticated)

	if appErr.Type != NotAuthenticated {
		t.Errorf("Expected type to be %s, got %s", NotAuthenticated, appErr.Type)
	}

	if appErr.Error() != string(notAuthenticatedErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", notAuthenticatedErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_NotAuthorized(t *testing.T) {
	appErr := NewAppErrorWithType(NotAuthorized)

	if appErr.Type != NotAuthorized {
		t.Errorf("Expected type to be %s, got %s", NotAuthorized, appErr.Type)
	}

	if appErr.Error() != string(notAuthorizedErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", notAuthorizedErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_TokenGeneratorError(t *testing.T) {
	appErr := NewAppErrorWithType(TokenGeneratorError)

	if appErr.Type != TokenGeneratorError {
		t.Errorf("Expected type to be %s, got %s", TokenGeneratorError, appErr.Type)
	}

	if appErr.Error() != string(tokenGeneratorErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", tokenGeneratorErrorMessage, appErr.Error())
	}
}

func TestNewAppErrorWithType_UnknownError(t *testing.T) {
	appErr := NewAppErrorWithType(ErrorType("UnknownType"))

	if appErr.Type != ErrorType("UnknownType") {
		t.Errorf("Expected type to be 'UnknownType', got %s", appErr.Type)
	}

	if appErr.Error() != string(unknownErrorMessage) {
		t.Errorf("Expected error message to be %s, got %s", unknownErrorMessage, appErr.Error())
	}
}

func TestAppError_Error(t *testing.T) {
	originalErr := errors.New("custom error message")
	appErr := &AppError{
		Err:  originalErr,
		Type: ErrorType("CustomError"),
	}

	if appErr.Error() != "custom error message" {
		t.Errorf("Expected error message to be 'custom error message', got %s", appErr.Error())
	}
}
