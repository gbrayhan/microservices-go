// Package auth provides the use case for authentication
package auth

import (
	"errors"
	"github.com/gbrayhan/microservices-go/application/security/jwt"
	errorsDomain "github.com/gbrayhan/microservices-go/domain/errors"
	userRepository "github.com/gbrayhan/microservices-go/infrastructure/repository/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Auth contains the data of the authentication
type Auth struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

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

	accessTokenClaims, err := jwt.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}
	refreshTokenClaims, err := jwt.GenerateJWTToken(domainUser.ID, "refresh")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}

	return secAuthUserMapper(domainUser, &Auth{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}), err
}

// AccessTokenByRefreshToken implements the Access Token By Refresh Token use case
func (s *Service) AccessTokenByRefreshToken(refreshToken string) (*SecurityAuthenticatedUser, error) {
	claimsMap, err := jwt.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	userMap := map[string]interface{}{"id": claimsMap["id"]}
	domainUser, err := s.UserRepository.GetOneByMap(userMap)
	if err != nil {
		return nil, err

	}

	accessTokenClaims, err := jwt.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}

	var expTime = int64(claimsMap["exp"].(float64))

	return secAuthUserMapper(domainUser, &Auth{
		AccessToken:               accessTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		RefreshToken:              refreshToken,
		ExpirationRefreshDateTime: time.Unix(expTime, 0),
	}), nil
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
