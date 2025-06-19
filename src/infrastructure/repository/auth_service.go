package repository

import (
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct{}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}

func (a *AuthServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
