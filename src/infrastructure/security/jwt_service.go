package security

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Access  = "access"
	Refresh = "refresh"
)

type AppToken struct {
	Token          string    `json:"token"`
	TokenType      string    `json:"type"`
	ExpirationTime time.Time `json:"expirationTime"`
}

type Claims struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTime    int64
	RefreshTime   int64
}

// IJWTService defines the interface for JWT operations
type IJWTService interface {
	GenerateJWTToken(userID int, tokenType string) (*AppToken, error)
	GetClaimsAndVerifyToken(tokenString string, tokenType string) (jwt.MapClaims, error)
}

// JWTService implements IJWTService
type JWTService struct {
	config JWTConfig
}

// NewJWTService creates a new JWT service instance
func NewJWTService() IJWTService {
	config := loadJWTConfig()
	return &JWTService{
		config: config,
	}
}

// NewJWTServiceWithConfig creates a new JWT service with custom configuration
func NewJWTServiceWithConfig(config JWTConfig) IJWTService {
	return &JWTService{
		config: config,
	}
}

// loadJWTConfig loads JWT configuration from environment variables
func loadJWTConfig() JWTConfig {
	return JWTConfig{
		AccessSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "default_access_secret"),
		RefreshSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "default_refresh_secret"),
		AccessTime:    getEnvAsInt64OrDefault("JWT_ACCESS_TIME_MINUTE", 60),
		RefreshTime:   getEnvAsInt64OrDefault("JWT_REFRESH_TIME_HOUR", 24),
	}
}

// GenerateJWTToken generates a JWT token for the given user ID and type
func (s *JWTService) GenerateJWTToken(userID int, tokenType string) (*AppToken, error) {
	var secretKey string
	var duration time.Duration

	switch tokenType {
	case Access:
		secretKey = s.config.AccessSecret
		duration = time.Duration(s.config.AccessTime) * time.Minute
	case Refresh:
		secretKey = s.config.RefreshSecret
		duration = time.Duration(s.config.RefreshTime) * time.Hour
	default:
		return nil, errors.New("invalid token type")
	}

	nowTime := time.Now()
	expirationTokenTime := nowTime.Add(duration)

	tokenClaims := &Claims{
		ID:   userID,
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTokenTime),
		},
	}
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	tokenStr, err := tokenWithClaims.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &AppToken{
		Token:          tokenStr,
		TokenType:      tokenType,
		ExpirationTime: expirationTokenTime,
	}, nil
}

// GetClaimsAndVerifyToken verifies a JWT token and returns its claims
func (s *JWTService) GetClaimsAndVerifyToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	var secretKey string

	if tokenType == Refresh {
		secretKey = s.config.RefreshSecret
	} else {
		secretKey = s.config.AccessSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, domainErrors.NewAppError(
				fmt.Errorf("unexpected signing method: %v", token.Header["alg"]),
				domainErrors.NotAuthenticated,
			)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims type or token not valid")
	}

	if claims["type"] != tokenType {
		return nil, domainErrors.NewAppError(errors.New("invalid token type"), domainErrors.NotAuthenticated)
	}
	var timeExpire = claims["exp"].(float64)
	if time.Now().Unix() > int64(timeExpire) {
		return nil, domainErrors.NewAppError(errors.New("token expired"), domainErrors.NotAuthenticated)
	}
	return claims, nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
