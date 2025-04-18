package errors

import "errors"

const (
	NotFound        = "NotFound"
	notFoundMessage = "record not found"

	ValidationError        = "ValidationError"
	validationErrorMessage = "validation error"

	ResourceAlreadyExists     = "ResourceAlreadyExists"
	alreadyExistsErrorMessage = "resource already exists"

	RepositoryError        = "RepositoryError"
	repositoryErrorMessage = "error in repository operation"

	NotAuthenticated             = "NotAuthenticated"
	notAuthenticatedErrorMessage = "not Authenticated"

	TokenGeneratorError        = "TokenGeneratorError"
	tokenGeneratorErrorMessage = "error in token generation"

	NotAuthorized             = "NotAuthorized"
	notAuthorizedErrorMessage = "not authorized"

	UnknownError        = "UnknownError"
	unknownErrorMessage = "something went wrong"
)

type AppError struct {
	Err  error
	Type string
}

func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType string) *AppError {
	var err error

	switch errType {
	case NotFound:
		err = errors.New(notFoundMessage)
	case ValidationError:
		err = errors.New(validationErrorMessage)
	case ResourceAlreadyExists:
		err = errors.New(alreadyExistsErrorMessage)
	case RepositoryError:
		err = errors.New(repositoryErrorMessage)
	case NotAuthenticated:
		err = errors.New(notAuthenticatedErrorMessage)
	case NotAuthorized:
		err = errors.New(notAuthorizedErrorMessage)
	case TokenGeneratorError:
		err = errors.New(tokenGeneratorErrorMessage)
	default:
		err = errors.New(unknownErrorMessage)
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}
