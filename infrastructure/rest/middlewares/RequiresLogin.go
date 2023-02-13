// Package middlewares contains the middlewares for the rest api
package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"net/http"
)

// AuthJWTMiddleware is a function that validates the jwt token
func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		viper.SetConfigFile("config.json")
		if err := viper.ReadInConfig(); err != nil {
			_ = fmt.Errorf("fatal error in config file: %s", err.Error())
		}

		JWTAccessSecure := viper.GetString("Secure.JWTAccessSecure")
		tokenString := c.GetHeader("Authorization")
		signature := []byte(JWTAccessSecure)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return signature, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
