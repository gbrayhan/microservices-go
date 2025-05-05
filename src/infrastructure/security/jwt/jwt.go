package jwt

import (
	"errors"
	"fmt"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strconv"
	"time"
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

func GenerateJWTToken(userID int, tokenType string) (*AppToken, error) {
	var secretKey string
	var expStr string

	switch tokenType {
	case Access:
		secretKey = os.Getenv("JWT_ACCESS_SECRET")
		expStr = os.Getenv("JWT_ACCESS_TIME_MINUTE")
	case Refresh:
		secretKey = os.Getenv("JWT_REFRESH_SECRET")
		expStr = os.Getenv("JWT_REFRESH_TIME_HOUR")
	default:
		return nil, errors.New("invalid token type")
	}

	if secretKey == "" || expStr == "" {
		return nil, errors.New("missing token environment variables")
	}

	tokenTimeConverted, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return nil, err
	}

	var duration time.Duration
	switch tokenType {
	case Refresh:
		duration = time.Duration(tokenTimeConverted) * time.Hour
	case Access:
		duration = time.Duration(tokenTimeConverted) * time.Minute
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

func GetClaimsAndVerifyToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	var secretKey string
	if tokenType == Refresh {
		secretKey = os.Getenv("JWT_REFRESH_SECRET")
	} else {
		secretKey = os.Getenv("JWT_ACCESS_SECRET")
	}
	if secretKey == "" {
		return nil, domainErrors.NewAppError(errors.New("missing token environment variable"), domainErrors.NotAuthenticated)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domainErrors.NewAppError(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]), domainErrors.NotAuthenticated)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != tokenType {
			return nil, domainErrors.NewAppError(errors.New("invalid token type"), domainErrors.NotAuthenticated)
		}
		var timeExpire = claims["exp"].(float64)
		if time.Now().Unix() > int64(timeExpire) {
			return nil, domainErrors.NewAppError(errors.New("token expired"), domainErrors.NotAuthenticated)
		}
		return claims, nil
	}
	return nil, err
}
