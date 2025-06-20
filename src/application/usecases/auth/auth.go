package auth

import (
	"errors"
	"time"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"

	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"

	errorsDomain "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	Email    string
	Password string
}

type DataUserAuthenticated struct {
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    bool   `json:"status"`
	ID        int    `json:"id"`
}

type DataSecurityAuthenticated struct {
	JWTAccessToken            string    `json:"jwtAccessToken"`
	JWTRefreshToken           string    `json:"jwtRefreshToken"`
	ExpirationAccessDateTime  time.Time `json:"expirationAccessDateTime"`
	ExpirationRefreshDateTime time.Time `json:"expirationRefreshDateTime"`
}

type SecurityAuthenticatedUser struct {
	Data     DataUserAuthenticated     `json:"data"`
	Security DataSecurityAuthenticated `json:"security"`
}

func secAuthUserMapper(domainUser *userDomain.User, authInfo *Auth) *SecurityAuthenticatedUser {
	return &SecurityAuthenticatedUser{
		Data: DataUserAuthenticated{
			UserName:  domainUser.UserName,
			Email:     domainUser.Email,
			FirstName: domainUser.FirstName,
			LastName:  domainUser.LastName,
			ID:        domainUser.ID,
			Status:    domainUser.Status,
		},
		Security: DataSecurityAuthenticated{
			JWTAccessToken:            authInfo.AccessToken,
			JWTRefreshToken:           authInfo.RefreshToken,
			ExpirationAccessDateTime:  authInfo.ExpirationAccessDateTime,
			ExpirationRefreshDateTime: authInfo.ExpirationRefreshDateTime,
		},
	}
}

type IAuthUseCase interface {
	Login(user LoginUser) (*SecurityAuthenticatedUser, error)
	AccessTokenByRefreshToken(refreshToken string) (*SecurityAuthenticatedUser, error)
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

	accessTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "access")
	if err != nil {
		return &SecurityAuthenticatedUser{}, err
	}
	refreshTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "refresh")
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
	claimsMap, err := s.jwtService.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		return nil, err
	}
	userMap := map[string]interface{}{"id": claimsMap["id"]}
	domainUser, err := s.userRepository.GetOneByMap(userMap)
	if err != nil {
		return nil, err
	}

	accessTokenClaims, err := s.jwtService.GenerateJWTToken(domainUser.ID, "access")
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
