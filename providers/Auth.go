package providers

import (
  "fmt"
  "github.com/golang-jwt/jwt/v4"
  "github.com/spf13/viper"
  "time"
)

type Auth struct {
  Token          string
  ExpirationTime time.Time
}

func (auth *Auth) GenerateJWTToken() (tokenString string, err error) {
  viper.SetConfigFile("config.json")
  if err := viper.ReadInConfig(); err != nil {
    _ = fmt.Errorf("fatal error in config file: %s \n", err.Error())
  }

  SecretKeyJWT := viper.GetString("datastore.metric.host")

  nowTime := time.Now()
  auth.ExpirationTime = nowTime.Add(time.Second * 10)

  // Create a new token object, specifying signing method and the claims
  // you would like it to contain.
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": "23",
    "nbf":     auth.ExpirationTime.Unix(),
  })

  // Sign and get the complete encoded token as a string using the secret
  auth.Token, err = token.SignedString(SecretKeyJWT)

  return
}
