// Package auth provides the use case for authentication
package auth

import (
	"errors"
	"github.com/gbrayhan/microservices-go/application/security/jwt"
	errorsDomain "github.com/gbrayhan/microservices-go/domain/errors"
	userRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/user"
	"golang.org/x/crypto/bcrypt"
)

// Service is a struct that contains the repository implementation for auth use case
type Service struct {
	UserRepository userRepository.Repository
}

// Login implements the login use case
func (s *Service) Login(user LoginUser) (*SecurityAuthenticatedUser, error) {
	userMap := map[string]interface{}{"email": user.Email}
	domainUser, err := s.UserRepository.GetOneByMap(userMap)
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}
	if domainUser.ID == 0 {
		return &SecurityAuthenticatedUser{}, errorsDomain.NewAppError(errors.New("email or password does not match"), errorsDomain.NotAuthorized)
	}

	isAuthenticated := CheckPasswordHash(user.Password, domainUser.HashPassword)
	if !isAuthenticated {
		err = errorsDomain.NewAppError(err, errorsDomain.NotAuthorized)
		return &SecurityAuthenticatedUser{}, errorsDomain.NewAppError(errors.New("email or password does not match"), errorsDomain.NotAuthorized)
	}

	authInfo, err := jwt.GenerateJWTTokens(domainUser.ID)
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}

	return secAuthUserMapper(domainUser, authInfo), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
