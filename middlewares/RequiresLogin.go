package middlewares

import "github.com/golang-jwt/jwt/v4"

import (
  "fmt"
  "time"
)

// Actions to authenticate client

func GenerateJWTToken() (string, error) {

  // Create a new token object, specifying signing method and the claims
  // you would like it to contain.
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "foo": "bar",
    "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
  })

  // Sign and get the complete encoded token as a string using the secret
  tokenString, err := token.SignedString("SecretKeyJWT")

  fmt.Println(tokenString, err)

}
