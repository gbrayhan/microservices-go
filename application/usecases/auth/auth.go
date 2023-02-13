// Package auth provides the use case for authentication
package auth

import (
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

	isAuthenticated := CheckPasswordHash(user.Password, domainUser.HashPassword)
	if !isAuthenticated {
		err = errorsDomain.NewAppError(err, errorsDomain.NotAuthorized)
		return &SecurityAuthenticatedUser{}, err
	}

	authInfo, err := jwt.GenerateJWTTokens(domainUser.ID)
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}

	return secAuthUserMapper(domainUser, authInfo), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
