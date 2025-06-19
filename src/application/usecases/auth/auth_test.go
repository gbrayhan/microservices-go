package auth

import (
	"errors"
	"os"
	"testing"

	jwtInfrastructure "github.com/gbrayhan/microservices-go/src/infrastructure/security"

	errorsDomain "github.com/gbrayhan/microservices-go/src/domain/errors"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
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
	os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "60")
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret")
	os.Setenv("JWT_REFRESH_TIME_HOUR", "24")

	tests := []struct {
		name                   string
		mockGetOneByMapFn      func(map[string]interface{}) (*userDomain.User, error)
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
			inputLogin: LoginUser{Email: "test@example.com", Password: "123456"},
			wantErr:    true,
		},
		{
			name: "User not found (ID=0)",
			mockGetOneByMapFn: func(m map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{ID: 0}, nil
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
			uc := NewAuthUseCase(userRepoMock)

			if tt.name == "Access token generation fails" {
				os.Setenv("JWT_ACCESS_TIME_MINUTE", "")
			}

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

			if tt.name == "Access token generation fails" {
				os.Setenv("JWT_ACCESS_TIME_MINUTE", "60")
			}
		})
	}
}

func TestAuthUseCase_AccessTokenByRefreshToken(t *testing.T) {
	os.Setenv("JWT_REFRESH_SECRET", "test_refresh_secret")
	os.Setenv("JWT_REFRESH_TIME_HOUR", "24")
	os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "60")

	validRefresh, _ := jwtInfrastructure.GenerateJWTToken(123, "refresh")

	tests := []struct {
		name            string
		refreshToken    string
		mockGetOneByMap func(map[string]interface{}) (*userDomain.User, error)
		modifySecretEnv bool
		wantErr         bool
	}{
		{
			name:         "Invalid token string -> claims error",
			refreshToken: "some.invalid.token",
			mockGetOneByMap: func(map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{}, nil
			},
			wantErr: true,
		},
		{
			name:         "DB error retrieving user",
			refreshToken: validRefresh.Token,
			mockGetOneByMap: func(map[string]interface{}) (*userDomain.User, error) {
				return nil, errors.New("db error")
			},
			wantErr: true,
		},
		{
			name:         "Access token generation error",
			refreshToken: validRefresh.Token,
			mockGetOneByMap: func(map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{ID: 999}, nil
			},
			modifySecretEnv: true,
			wantErr:         true,
		},
		{
			name:         "OK - valid refresh token and user found",
			refreshToken: validRefresh.Token,
			mockGetOneByMap: func(map[string]interface{}) (*userDomain.User, error) {
				return &userDomain.User{ID: 999}, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := &mockUserService{
				getOneByMapFn: tt.mockGetOneByMap,
			}
			uc := NewAuthUseCase(userRepoMock)

			if tt.modifySecretEnv {
				os.Setenv("JWT_ACCESS_SECRET", "")
			}

			resp, err := uc.AccessTokenByRefreshToken(tt.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("[%s] got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
			}
			if !tt.wantErr && resp.Security.JWTAccessToken == "" {
				t.Errorf("[%s] expected a new AccessToken, got empty", tt.name)
			}

			if tt.modifySecretEnv {
				os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
			}
		})
	}
}
