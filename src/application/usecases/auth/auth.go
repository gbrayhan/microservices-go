package auth

import (
	"errors"
	"time"

	jwtInfrastructure "github.com/gbrayhan/microservices-go/src/infrastructure/security"

	errorsDomain "github.com/gbrayhan/microservices-go/src/domain/errors"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUseCase interface {
	Login(user LoginUser) (*SecurityAuthenticatedUser, error)
	AccessTokenByRefreshToken(refreshToken string) (*SecurityAuthenticatedUser, error)
}

type AuthUseCase struct {
	userRepository userDomain.IUserService
}

func NewAuthUseCase(userRepository userDomain.IUserService) IAuthUseCase {
	return &AuthUseCase{
		userRepository: userRepository,
	}
}

type Auth struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

func (s *AuthUseCase) Login(user LoginUser) (*SecurityAuthenticatedUser, error) {
	userMap := map[string]interface{}{"email": user.Email}
	domainUser, err := s.userRepository.GetOneByMap(userMap)
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}
	if domainUser.ID == 0 {
		return &SecurityAuthenticatedUser{}, errorsDomain.NewAppError(errors.New("email or password does not match"), errorsDomain.NotAuthorized)
	}

	isAuthenticated := checkPasswordHash(user.Password, domainUser.HashPassword)
	if !isAuthenticated {
		return &SecurityAuthenticatedUser{}, errorsDomain.NewAppError(errors.New("email or password does not match"), errorsDomain.NotAuthorized)
	}

	accessTokenClaims, err := jwtInfrastructure.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}
	refreshTokenClaims, err := jwtInfrastructure.GenerateJWTToken(domainUser.ID, "refresh")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}

	return secAuthUserMapper(domainUser, &Auth{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}), nil
}

func (s *AuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*SecurityAuthenticatedUser, error) {
	claimsMap, err := jwtInfrastructure.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}
	userMap := map[string]interface{}{"id": claimsMap["id"]}
	domainUser, err := s.userRepository.GetOneByMap(userMap)
	if err != nil {
		return nil, err
	}

	accessTokenClaims, err := jwtInfrastructure.GenerateJWTToken(domainUser.ID, "access")
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
