// Package jwt implements the JWT authentication
package jwt

import (
	"errors"
	"fmt"
	domainErrors "github.com/gbrayhan/microservices-go/domain/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

const (
	Access  = "access"
	Refresh = "refresh"
)

// AppToken is a struct that contains the JWT token
type AppToken struct {
	Token          string    `json:"token"`
	TokenType      string    `json:"type"`
	ExpirationTime time.Time `json:"expitationTime"`
}

// Claims is a struct that contains the claims of the JWT
type Claims struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// TokenTypeKeyName is a map that contains the key name of the JWT in config.json
var TokenTypeKeyName = map[string]string{
	Access:  "Secure.JWTAccessSecure",
	Refresh: "Secure.JWTRefreshSecure",
}

// TokenTypeExpTime is a map that contains the expiration time of the JWT
var TokenTypeExpTime = map[string]string{
	Access:  "Secure.JWTAccessTimeMinute",
	Refresh: "Secure.JWTRefreshTimeHour",
}

// GenerateJWTToken generates a JWT token (refresh or access)
func GenerateJWTToken(userID int, tokenType string) (appToken *AppToken, err error) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s", err.Error())
	}

	JWTSecureKey := viper.GetString(TokenTypeKeyName[tokenType])
	JWTExpTime := viper.GetString(TokenTypeExpTime[tokenType])

	tokenTimeConverted, err := strconv.ParseInt(JWTExpTime, 10, 64)
	if err != nil {
		return
	}

	tokenTimeUnix := time.Duration(tokenTimeConverted)
	switch tokenType {
	case Refresh:
		tokenTimeUnix *= time.Hour
	case Access:
		tokenTimeUnix *= time.Minute
	default:
		err = errors.New("invalid token type")
	}

	if err != nil {
		return
	}
	nowTime := time.Now()
	expirationTokenTime := nowTime.Add(tokenTimeUnix)

	tokenClaims := &Claims{
		ID:   userID,
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTokenTime),
		},
	}
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := tokenWithClaims.SignedString([]byte(JWTSecureKey))
	if err != nil {
		return
	}

	appToken = &AppToken{
		Token:          tokenStr,
		TokenType:      tokenType,
		ExpirationTime: expirationTokenTime,
	}

	return
}

// GetClaimsAndVerifyToken verifies the token and returns the claims
func GetClaimsAndVerifyToken(tokenString string, tokenType string) (claims jwt.MapClaims, err error) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s", err.Error())
	}
	JWTRefreshSecure := viper.GetString(TokenTypeKeyName[tokenType])
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domainErrors.NewAppError(errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])), domainErrors.NotAuthenticated)
		}

		return []byte(JWTRefreshSecure), nil
	})

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
