package auth

import (
	"errors"
	"testing"
	"time"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type mockUserService struct {
	getOneByMapFn         func(map[string]interface{}) (*domainUser.User, error)
	callGetOneByMapCalled bool
}

func (m *mockUserService) GetAll() (*[]domainUser.User, error) {
	return nil, nil
}
func (m *mockUserService) GetByID(id int) (*domainUser.User, error) {
	return nil, nil
}
func (m *mockUserService) Create(newUser *domainUser.User) (*domainUser.User, error) {
	return nil, nil
}
func (m *mockUserService) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	m.callGetOneByMapCalled = true
	return m.getOneByMapFn(userMap)
}
func (m *mockUserService) Delete(id int) error {
	return nil
}
func (m *mockUserService) Update(id int, userMap map[string]interface{}) (*domainUser.User, error) {
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

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
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
		mockGetOneByMapFn      func(map[string]interface{}) (*domainUser.User, error)
		mockGenerateTokenFn    func(int, string) (*security.AppToken, error)
		inputEmail             string
		inputPassword          string
		wantErr                bool
		wantErrType            domainErrors.ErrorType
		wantEmptySecurity      bool
		wantSuccessAccessToken bool
	}{
		{
			name: "Error fetching user from DB",
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return nil, errors.New("db error")
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputEmail:    "test@example.com",
			inputPassword: "123456",
			wantErr:       true,
		},
		{
			name: "User not found (ID=0)",
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 0}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputEmail:    "test@example.com",
			inputPassword: "123456",
			wantErr:       true,
			wantErrType:   domainErrors.NotAuthorized,
		},
		{
			name: "Incorrect password",
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				hashed, _ := HashPasswordForTest("someOtherPass")
				return &domainUser.User{ID: 10, HashPassword: hashed}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "test_token"}, nil
			},
			inputEmail:        "test@example.com",
			inputPassword:     "wrong",
			wantErr:           true,
			wantErrType:       domainErrors.NotAuthorized,
			wantEmptySecurity: true,
		},
		{
			name: "Access token generation fails",
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				hashed, _ := HashPasswordForTest("somePass")
				return &domainUser.User{ID: 10, HashPassword: hashed}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return nil, errors.New("token generation failed")
			},
			inputEmail:    "test@example.com",
			inputPassword: "somePass",
			wantErr:       true,
		},
		{
			name: "OK - everything correct",
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				hashed, _ := HashPasswordForTest("mySecretPass")
				return &domainUser.User{
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
			inputEmail:             "test@example.com",
			inputPassword:          "mySecretPass",
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

			logger := setupLogger(t)
			uc := NewAuthUseCase(userRepoMock, jwtMock, logger)

			user, authTokens, err := uc.Login(tt.inputEmail, tt.inputPassword)
			if (err != nil) != tt.wantErr {
				t.Fatalf("[%s] got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErrType != "" && err != nil {
				appErr, ok := err.(*domainErrors.AppError)
				if !ok || appErr.Type != tt.wantErrType {
					t.Errorf("[%s] expected error type = %s, got = %v", tt.name, tt.wantErrType, err)
				}
			}

			if !tt.wantErr && tt.wantSuccessAccessToken {
				if authTokens.AccessToken == "" {
					t.Errorf("[%s] expected a non-empty AccessToken, got empty", tt.name)
				}
				if user == nil {
					t.Errorf("[%s] expected a non-nil user, got nil", tt.name)
				}
			} else if tt.wantErr && tt.wantEmptySecurity {
				if authTokens != nil && authTokens.AccessToken != "" {
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
		mockGetOneByMapFn   func(map[string]interface{}) (*domainUser.User, error)
		mockGenerateTokenFn func(int, string) (*security.AppToken, error)
		inputRefreshToken   string
		wantErr             bool
		wantErrType         domainErrors.ErrorType
		wantSuccess         bool
	}{
		{
			name: "Invalid refresh token",
			mockVerifyTokenFn: func(token, tokenType string) (jwt.MapClaims, error) {
				return nil, errors.New("invalid token")
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 10}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "new_access_token"}, nil
			},
			inputRefreshToken: "invalid_token",
			wantErr:           true,
		},
		{
			name: "User not found after token verification",
			mockVerifyTokenFn: func(token, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(999)}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return nil, errors.New("user not found")
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "new_access_token"}, nil
			},
			inputRefreshToken: "valid_token",
			wantErr:           true,
		},
		{
			name: "New access token generation fails",
			mockVerifyTokenFn: func(token, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(10)}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 10}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return nil, errors.New("token generation failed")
			},
			inputRefreshToken: "valid_token",
			wantErr:           true,
		},
		{
			name: "OK - successful token refresh",
			mockVerifyTokenFn: func(token, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(10), "exp": float64(time.Now().Add(time.Hour).Unix())}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 10, Email: "test@example.com"}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "new.token", TokenType: tokenType, ExpirationTime: time.Now().Add(time.Hour)}, nil
			},
			inputRefreshToken: "valid_refresh_token",
			wantErr:           false,
			wantSuccess:       true,
		},
		{
			name: "Refresh token generation fails",
			mockVerifyTokenFn: func(token string, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(10), "type": "refresh"}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 10}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return nil, errors.New("token generation failed")
			},
			inputRefreshToken: "valid.refresh.token",
			wantErr:           true,
		},
		{
			name: "OK - everything correct",
			mockVerifyTokenFn: func(token string, tokenType string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"id": float64(10), "type": "refresh", "exp": float64(time.Now().Add(time.Hour).Unix())}, nil
			},
			mockGetOneByMapFn: func(m map[string]interface{}) (*domainUser.User, error) {
				return &domainUser.User{ID: 10}, nil
			},
			mockGenerateTokenFn: func(userID int, tokenType string) (*security.AppToken, error) {
				return &security.AppToken{Token: "new.token", TokenType: tokenType, ExpirationTime: time.Now().Add(time.Hour)}, nil
			},
			inputRefreshToken: "valid.refresh.token",
			wantErr:           false,
			wantSuccess:       true,
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

			logger := setupLogger(t)
			uc := NewAuthUseCase(userRepoMock, jwtMock, logger)

			user, authTokens, err := uc.AccessTokenByRefreshToken(tt.inputRefreshToken)
			if (err != nil) != tt.wantErr {
				t.Fatalf("[%s] got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
			}

			if tt.wantErrType != "" && err != nil {
				appErr, ok := err.(*domainErrors.AppError)
				if !ok || appErr.Type != tt.wantErrType {
					t.Errorf("[%s] expected error type = %s, got = %v", tt.name, tt.wantErrType, err)
				}
			}

			if !tt.wantErr && tt.wantSuccess {
				if authTokens.AccessToken == "" {
					t.Errorf("[%s] expected a non-empty AccessToken, got empty", tt.name)
				}
				if user == nil {
					t.Errorf("[%s] expected a non-nil user, got nil", tt.name)
				}
			}
		})
	}
}
