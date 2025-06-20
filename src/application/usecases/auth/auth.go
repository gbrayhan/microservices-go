package auth

import (
	"errors"
	"time"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"

	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUseCase interface {
	Login(email, password string) (*userDomain.User, *AuthTokens, error)
	AccessTokenByRefreshToken(refreshToken string) (*userDomain.User, *AuthTokens, error)
}

type AuthUseCase struct {
	userRepository user.UserRepositoryInterface
	jwtService     security.IJWTService
}

func NewAuthUseCase(userRepository user.UserRepositoryInterface, jwtService security.IJWTService) IAuthUseCase {
	return &AuthUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

type AuthTokens struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

func (s *AuthUseCase) Login(email, password string) (*userDomain.User, *AuthTokens, error) {
	userMap := map[string]interface{}{"email": email}
	domainUser, err := s.userRepository.GetOneByMap(userMap)
	if err != nil {
		return nil, nil, err
	}
	if domainUser.ID == 0 {
		return nil, nil, domainErrors.NewAppError(errors.New("email or password does not match"), domainErrors.NotAuthorized)
	}

	isAuthenticated := checkPasswordHash(password, domainUser.HashPassword)
	if !isAuthenticated {
		return nil, nil, domainErrors.NewAppError(errors.New("email or password does not match"), domainErrors.NotAuthorized)
	}

	accessTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return nil, nil, err
	}
	refreshTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "refresh")
	if err != nil {
		return nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}

	return domainUser, authTokens, nil
}

func (s *AuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*userDomain.User, *AuthTokens, error) {
	claimsMap, err := s.jwtService.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, nil, err
	}
	userMap := map[string]interface{}{"id": claimsMap["id"]}
	domainUser, err := s.userRepository.GetOneByMap(userMap)
	if err != nil {
		return nil, nil, err
	}

	accessTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return nil, nil, err
	}

	var expTime = int64(claimsMap["exp"].(float64))

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		RefreshToken:              refreshToken,
		ExpirationRefreshDateTime: time.Unix(expTime, 0),
	}

	return domainUser, authTokens, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
