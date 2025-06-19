package auth

import (
	"errors"
	"testing"
	"time"

	errorsDomain "github.com/gbrayhan/microservices-go/src/domain/errors"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type mockUserService struct {
	getOneByMapFn         func(map[string]interface{}) (*userDomain.User, error)
	callGetOneByMapCalled bool
}

func (m *mockUserService) GetAll() (*[]userDomain.User, error) {
	return nil, nil
}
func (m *mockUserService) GetByID(id int) (*userDomain.User, error) {
	return nil, nil
}
func (m *mockUserService) Create(newUser *userDomain.User) (*userDomain.User, error) {
	return nil, nil
}
func (m *mockUserService) GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error) {
	m.callGetOneByMapCalled = true
	return m.getOneByMapFn(userMap)
}
func (m *mockUserService) Delete(id int) error {
	return nil
}
func (m *mockUserService) Update(id int, userMap map[string]interface{}) (*userDomain.User, error) {
	return nil, nil
}

type mockJWTService struct {
	generateTokenFn func(int, string) (*security.AppToken, error)
	verifyTokenFn   func(string, string) (jwt.MapClaims, error)
}

func (m *mockJWTService) GenerateJWTToken(userID int, tokenType string) (*security.AppToken, error) {
	return m.generateTokenFn(userID, tokenType)
}

func (m *mockJWTService) GetClaimsAndVerifyToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	return m.verifyTokenFn(tokenString, tokenType)
}

func HashPasswordForTest(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecretPass"
	hashed, err := HashPasswordForTest(password)
	if err != nil {
		t.Fatalf("failed to generate hash for test: %v", err)
	}

	if ok := checkPasswordHash(password, hashed); !ok {
		t.Errorf("checkPasswordHash() = false, want true")
	}

	if ok := checkPasswordHash("wrongPassword", hashed); ok {
		t.Errorf("checkPasswordHash() = true, want false")
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	tests := []struct {
		name                   string
		mockGetOneByMapFn      func(map[string]interface{}) (*userDomain.User, error)
		mockGenerateTokenFn    func(int, string) (*security.AppToken, error)
		inputLogin             LoginUser
		wantErr                bool
		wantErrType            string
		wantEmptySecurity      bool
		wantSuccessAccessToken bool
	}{
		{
			name: "Error fetching user from DB",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				return nil, errors.New("db error")
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputLogin: LoginUser{Email: "test@example.com", Password: "123456"},
			wantErr:    true,
		},
		{
			name: "User not found (ID=0)",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{ID: 0}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputLogin:  LoginUser{Email: "test@example.com", Password: "123456"},
			wantErr:     true,
			wantErrType: errorsDomain.NotAuthorized,
		},
		{
			name: "Incorrect password",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				hashed, _ := HashPasswordForTest("someOtherPass")
				return &userDomain.User{ID: 10, HashPassword: hashed}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputLogin:        LoginUser{Email: "test@example.com", Password: "wrong"},
			wantErr:           true,
			wantErrType:       errorsDomain.NotAuthorized,
			wantEmptySecurity: true,
		},
		{
			name: "Access token generation fails",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				hashed, _ := HashPasswordForTest("somePass")
				return &userDomain.User{ID: 10, HashPassword: hashed}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return nil, errors.New("token generation failed")
			},
			inputLogin: LoginUser{Email: "test@example.com", Password: "somePass"},
			wantErr:    true,
		},
		{
			name: "OK - everything correct",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				hashed, _ := HashPasswordForTest("mySecretPass")
				return &userDomain.User{
					ID:           10,
					Email:        "test@example.com",
					HashPassword: hashed,
				}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{
					Token:          "test_token_" + tokenType,
					ExpirationTime: time.Now().Add(time.Hour),
				}, nil
			},
			inputLogin:             LoginUser{Email: "test@example.com", Password: "mySecretPass"},
			wantErr:                false,
			wantSuccessAccessToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := &mockUserService{
				getOneByMapFn: tt.mockGetOneByMapFn,
			}

			jwtMock := &mockJWTService{
				generateTokenFn: tt.mockGenerateTokenFn,
			}

			uc := NewAuthUseCase(userRepoMock, jwtMock)

			result, err := uc.Login(tt.inputLogin)
			if (err != nil) != tt.wantErr {
				t.Fatalf("[%s] got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErrType != "" && err != nil {
				appErr, ok := err.(*errorsDomain.AppError)
				if !ok || appErr.Type != tt.wantErrType {
					t.Errorf("[%s] expected error type = %s, got = %v", tt.name, tt.wantErrType, err)
				}
			}

			if !tt.wantErr && tt.wantSuccessAccessToken {
				if result.Security.JWTAccessToken == "" {
					t.Errorf("[%s] expected a non-empty AccessToken, got empty", tt.name)
				}
			} else if tt.wantErr && tt.wantEmptySecurity {
				if result.Security.JWTAccessToken != "" {
					t.Errorf("[%s] expected empty AccessToken, but got a non-empty one", tt.name)
				}
			}
		})
	}
}

func TestAuthUseCase_AccessTokenByRefreshToken(t *testing.T) {
	tests := []struct {
		name                string
		mockVerifyTokenFn   func(string, string) (jwt.MapClaims, error)
		mockGetOneByMapFn   func(map[string]interface{}) (*userDomain.User, error)
		mockGenerateTokenFn func(int, string) (*security.AppToken, error)
		inputToken          string
		wantErr             bool
		wantErrType         string
	}{
		{
			name: "Invalid token string -> claims error",
			mockVerifyTokenFn: func(tokenString, tokenType string) (jwt.MapClaims, error) {
				return nil, errors.New("invalid token")
			},
			inputToken: "invalid_token",
			wantErr:    true,
		},
		{
			name: "DB error retrieving user",
			mockVerifyTokenFn: func(tokenString, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(1)}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				return nil, errors.New("db error")
			},
			inputToken: "valid_token",
			wantErr:    true,
		},
		{
			name: "OK - valid refresh token",
			mockVerifyTokenFn: func(tokenString, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{
					"id":  float64(1),
					"exp": float64(time.Now().Add(time.Hour).Unix()),
				}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{ID: 1}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{
					Token:          "new_access_token",
					ExpirationTime: time.Now().Add(time.Hour),
				}, nil
			},
			inputToken: "valid_refresh_token",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := &mockUserService{
				getOneByMapFn: tt.mockGetOneByMapFn,
			}

			jwtMock := &mockJWTService{
				verifyTokenFn:   tt.mockVerifyTokenFn,
				generateTokenFn: tt.mockGenerateTokenFn,
			}

			uc := NewAuthUseCase(userRepoMock, jwtMock)

			result, err := uc.AccessTokenByRefreshToken(tt.inputToken)
			if (err != nil) != tt.wantErr {
				t.Fatalf("[%s] got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErrType != "" && err != nil {
				appErr, ok := err.(*errorsDomain.AppError)
				if !ok || appErr.Type != tt.wantErrType {
					t.Errorf("[%s] expected error type = %s, got = %v", tt.name, tt.wantErrType, err)
				}
			}

			if !tt.wantErr {
				if result.Security.JWTAccessToken == "" {
					t.Errorf("[%s] expected a non-empty AccessToken, got empty", tt.name)
				}
			}
		})
	}
}
