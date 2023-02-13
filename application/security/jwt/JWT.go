// Package jwt implements the JWT authentication
package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type Auth struct {
	AccessToken           string
	RefreshToken          string
	ExpirationAccessTime  time.Time
	ExpirationRefreshTime time.Time
}

type Claims struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateJWTTokens(userID int) (authData *Auth, err error) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s", err.Error())
	}

	JWTAccessSecure := viper.GetString("Secure.JWTAccessSecure")
	JWTRefreshSecure := viper.GetString("Secure.JWTRefreshSecure")

	JWTAccessTimeMinute := viper.GetString("Secure.JWTAccessTimeMinute")
	JWTRefreshTimeHour := viper.GetString("Secure.JWTRefreshTimeHour")

	AccessTimeMinute, err := strconv.ParseInt(JWTAccessTimeMinute, 10, 64)
	if err != nil {
		return
	}
	RefreshTimeHour, err := strconv.ParseInt(JWTRefreshTimeHour, 10, 64)
	if err != nil {
		return
	}

	AccessTimeDurationMinute := time.Minute * time.Duration(AccessTimeMinute)
	RefreshTimeDurationHour := time.Hour * time.Duration(RefreshTimeHour)

	nowTime := time.Now()
	expirationAccessTime := nowTime.Add(AccessTimeDurationMinute)
	expirationRefreshTime := nowTime.Add(RefreshTimeDurationHour)

	accessClaims := &Claims{
		ID:   userID,
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationAccessTime),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshClaims := &Claims{
		ID:   userID,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationRefreshTime),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenStr, err := accessToken.SignedString([]byte(JWTAccessSecure))
	if err != nil {
		return
	}

	refreshTokenStr, err := refreshToken.SignedString([]byte(JWTRefreshSecure))
	if err != nil {
		return
	}

	authData = &Auth{
		ExpirationAccessTime:  expirationAccessTime,
		ExpirationRefreshTime: expirationRefreshTime,
		AccessToken:           accessTokenStr,
		RefreshToken:          refreshTokenStr,
	}

	return
}

func GetClaimsAndVerifyAccessToken(tokenString string) (claims jwt.MapClaims, err error) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s", err.Error())
	}
	JWTAccessSecure := viper.GetString("Secure.JWTAccessSecure")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTAccessSecure), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func GetClaimsAndVerifyRefreshToken(tokenString string) (claims jwt.MapClaims, err error) {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("fatal error in config file: %s", err.Error())
	}
	JWTRefreshSecure := viper.GetString("Secure.JWTRefreshSecure")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTRefreshSecure), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
