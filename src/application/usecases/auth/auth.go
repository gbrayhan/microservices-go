package auth

import (
	"errors"
	"time"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IAuthUseCase interface {
	Login(email, password string) (*domainUser.User, *AuthTokens, error)
	AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error)
}

type AuthUseCase struct {
	UserRepository user.UserRepositoryInterface
	JWTService     security.IJWTService
	Logger         *logger.Logger
}

func NewAuthUseCase(userRepository user.UserRepositoryInterface, jwtService security.IJWTService, loggerInstance *logger.Logger) IAuthUseCase {
	return &AuthUseCase{
		UserRepository: userRepository,
		JWTService:     jwtService,
		Logger:         loggerInstance,
	}
}

type AuthTokens struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

func (s *AuthUseCase) Login(email, password string) (*domainUser.User, *AuthTokens, error) {
	s.Logger.Info("User login attempt", zap.String("email", email))
	user, err := s.UserRepository.GetByEmail(email)
	if err != nil {
		s.Logger.Error("Error getting user for login", zap.Error(err), zap.String("email", email))
		return nil, nil, err
	}
	if user.ID == 0 {
		s.Logger.Warn("Login failed: user not found", zap.String("email", email))
		return nil, nil, domainErrors.NewAppError(errors.New("email or password does not match"), domainErrors.NotAuthenticated)
	}

	isAuthenticated := checkPasswordHash(password, user.HashPassword)
	if !isAuthenticated {
		s.Logger.Warn("Login failed: invalid password", zap.String("email", email))
		return nil, nil, domainErrors.NewAppError(errors.New("email or password does not match"), domainErrors.NotAuthenticated)
	}

	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, "access")
	if err != nil {
		s.Logger.Error("Error generating access token", zap.Error(err), zap.Int("userID", user.ID))
		return nil, nil, err
	}
	refreshTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, "refresh")
	if err != nil {
		s.Logger.Error("Error generating refresh token", zap.Error(err), zap.Int("userID", user.ID))
		return nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}

	s.Logger.Info("User login successful", zap.String("email", email), zap.Int("userID", user.ID))
	return user, authTokens, nil
}

func (s *AuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error) {
	s.Logger.Info("Refreshing access token")
	claimsMap, err := s.JWTService.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		s.Logger.Error("Error verifying refresh token", zap.Error(err))
		return nil, nil, err
	}
	userID := int(claimsMap["id"].(float64))
	user, err := s.UserRepository.GetByID(userID)
	if err != nil {
		s.Logger.Error("Error getting user for token refresh", zap.Error(err), zap.Int("userID", userID))
		return nil, nil, err
	}

	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, "access")
	if err != nil {
		s.Logger.Error("Error generating new access token", zap.Error(err), zap.Int("userID", user.ID))
		return nil, nil, err
	}

	var expTime = int64(claimsMap["exp"].(float64))

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		RefreshToken:              refreshToken,
		ExpirationRefreshDateTime: time.Unix(expTime, 0),
	}

	s.Logger.Info("Access token refreshed successfully", zap.Int("userID", user.ID))
	return user, authTokens, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
