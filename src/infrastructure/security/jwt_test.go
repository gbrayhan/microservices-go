package security

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	service := NewJWTService()
	assert.NotNil(t, service)
	assert.Implements(t, (*IJWTService)(nil), service)
}

func TestNewJWTServiceWithConfig(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)
	assert.NotNil(t, service)
	assert.Implements(t, (*IJWTService)(nil), service)
}

func TestLoadJWTConfig(t *testing.T) {
	// Test with environment variables set
	os.Setenv("JWT_ACCESS_SECRET_KEY", "custom_access_secret")
	os.Setenv("JWT_REFRESH_SECRET_KEY", "custom_refresh_secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "45")
	os.Setenv("JWT_REFRESH_TIME_HOUR", "48")

	config := loadJWTConfig()
	assert.Equal(t, "custom_access_secret", config.AccessSecret)
	assert.Equal(t, "custom_refresh_secret", config.RefreshSecret)
	assert.Equal(t, int64(45), config.AccessTime)
	assert.Equal(t, int64(48), config.RefreshTime)

	// Clean up
	os.Unsetenv("JWT_ACCESS_SECRET_KEY")
	os.Unsetenv("JWT_REFRESH_SECRET_KEY")
	os.Unsetenv("JWT_ACCESS_TIME_MINUTE")
	os.Unsetenv("JWT_REFRESH_TIME_HOUR")
}

func TestGenerateJWTToken_Access(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, Access, token.TokenType)
	assert.True(t, token.ExpirationTime.After(time.Now()))
}

func TestGenerateJWTToken_Refresh(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 456
	token, err := service.GenerateJWTToken(userID, Refresh)
	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, Refresh, token.TokenType)
	assert.True(t, token.ExpirationTime.After(time.Now()))
}

func TestGenerateJWTToken_InvalidType(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, "invalid_type")
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "invalid token type")
}

func TestGenerateJWTToken_EmptySecret(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "",
		RefreshSecret: "",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	// This should still work with empty secrets (they're just empty strings)
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestGetClaimsAndVerifyToken_ValidAccessToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Access)
	require.NoError(t, err)
	assert.Equal(t, float64(userID), claims["id"])
	assert.Equal(t, Access, claims["type"])
	assert.NotNil(t, claims["exp"])
}

func TestGetClaimsAndVerifyToken_ValidRefreshToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 456
	token, err := service.GenerateJWTToken(userID, Refresh)
	require.NoError(t, err)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Refresh)
	require.NoError(t, err)
	assert.Equal(t, float64(userID), claims["id"])
	assert.Equal(t, Refresh, claims["type"])
	assert.NotNil(t, claims["exp"])
}

func TestGetClaimsAndVerifyToken_InvalidToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err := service.GetClaimsAndVerifyToken("invalid_token", Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_WrongTokenType(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	// Generate access token but try to verify as refresh token
	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Refresh)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_ExpiredToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    0, // 0 minutes = immediate expiration
		RefreshTime:   0, // 0 hours = immediate expiration
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(1 * time.Second)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_WrongSigningMethod(t *testing.T) {
	// Create a token with wrong signing method
	claims := jwt.MapClaims{
		"id":   123,
		"type": Access,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("wrong_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err = service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_EmptyToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err := service.GetClaimsAndVerifyToken("", Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_MalformedToken(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err := service.GetClaimsAndVerifyToken("header.payload.signature", Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_TokenWithoutClaims(t *testing.T) {
	// Create a token without proper claims
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err := service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_TokenWithInvalidClaims(t *testing.T) {
	// Create a token with invalid claims structure
	claims := jwt.MapClaims{
		"id":   "not_a_number", // Should be a number
		"type": Access,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err = service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_TokenWithInvalidExpiration(t *testing.T) {
	// Create a token with invalid expiration format
	claims := jwt.MapClaims{
		"id":   123,
		"type": Access,
		"exp":  "not_a_number", // Should be a number
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err = service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_TokenWithMissingType(t *testing.T) {
	// Create a token without type claim
	claims := jwt.MapClaims{
		"id":  123,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err = service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_TokenWithMissingExpiration(t *testing.T) {
	// Create a token without expiration claim
	claims := jwt.MapClaims{
		"id":   123,
		"type": Access,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	claims, err = service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_KEY", "test_value")
	result := getEnvOrDefault("TEST_KEY", "default_value")
	assert.Equal(t, "test_value", result)

	// Test with environment variable not set
	result = getEnvOrDefault("NONEXISTENT_KEY", "default_value")
	assert.Equal(t, "default_value", result)

	// Clean up
	os.Unsetenv("TEST_KEY")
}

func TestGetEnvAsInt64OrDefault(t *testing.T) {
	// Test with valid integer environment variable
	os.Setenv("TEST_INT", "123")
	result := getEnvAsInt64OrDefault("TEST_INT", 456)
	assert.Equal(t, int64(123), result)

	// Test with invalid integer environment variable
	os.Setenv("TEST_INVALID", "not_a_number")
	result = getEnvAsInt64OrDefault("TEST_INVALID", 456)
	assert.Equal(t, int64(456), result)

	// Test with environment variable not set
	result = getEnvAsInt64OrDefault("NONEXISTENT_INT", 789)
	assert.Equal(t, int64(789), result)

	// Clean up
	os.Unsetenv("TEST_INT")
	os.Unsetenv("TEST_INVALID")
}

func TestJWTService_InterfaceCompliance(t *testing.T) {
	var _ IJWTService = (*JWTService)(nil)
}

func TestGenerateJWTToken_EdgeCases(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    0, // Test with zero duration
		RefreshTime:   0,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)

	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.ExpirationTime.Equal(time.Now()) || token.ExpirationTime.Before(time.Now()))
}

func TestGenerateJWTToken_NegativeUserID(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := -123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestGenerateJWTToken_ZeroUserID(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 0
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestGenerateJWTToken_LargeUserID(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 999999999
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)
	assert.NotNil(t, token)
}

func TestGetClaimsAndVerifyToken_WithDifferentSecrets(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "access_secret_1",
		RefreshSecret: "refresh_secret_2",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	// Should work with access secret
	claims, err := service.GetClaimsAndVerifyToken(token.Token, Access)
	require.NoError(t, err)
	assert.Equal(t, float64(userID), claims["id"])

	// Should fail with refresh secret
	claims, err = service.GetClaimsAndVerifyToken(token.Token, Refresh)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaimsAndVerifyToken_WithVeryLongSecrets(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "very_long_access_secret_that_exceeds_normal_length_for_testing_purposes_and_to_ensure_proper_handling_of_long_secrets_in_the_jwt_service_implementation",
		RefreshSecret: "very_long_refresh_secret_that_exceeds_normal_length_for_testing_purposes_and_to_ensure_proper_handling_of_long_secrets_in_the_jwt_service_implementation",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Access)
	require.NoError(t, err)
	assert.Equal(t, float64(userID), claims["id"])
}

func TestGetClaimsAndVerifyToken_WithSpecialCharactersInSecrets(t *testing.T) {
	config := JWTConfig{
		AccessSecret:  "access_secret_with_special_chars_!@#$%^&*()_+-=[]{}|;':\",./<>?",
		RefreshSecret: "refresh_secret_with_special_chars_!@#$%^&*()_+-=[]{}|;':\",./<>?",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	userID := 123
	token, err := service.GenerateJWTToken(userID, Access)
	require.NoError(t, err)

	claims, err := service.GetClaimsAndVerifyToken(token.Token, Access)
	require.NoError(t, err)
	assert.Equal(t, float64(userID), claims["id"])
}

// BadClaims es un tipo inválido para forzar error en SignedString
type BadClaims struct{}

func (b BadClaims) Valid() error { return nil }

// TestGetClaimsAndVerifyToken_SignedStringError cubre el branch de error de SignedString, pero no es forzable sin cambiar la API de GenerateJWTToken.
// Por defecto, jwt-go permite serializar cualquier struct vacío y el secret raro no genera error.
// Si se desea cubrir este branch, se debe usar un mock o cambiar la firma de la función para inyectar el error.
//func TestGenerateJWTToken_SignedStringError(t *testing.T) {
//	config := JWTConfig{
//		AccessSecret:  string([]byte{0xff, 0xfe, 0xfd}), // secret raro
//		RefreshSecret: string([]byte{0xff, 0xfe, 0xfd}),
//		AccessTime:    30,
//		RefreshTime:   24,
//	}
//	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, BadClaims{})
//	_, err := tokenWithClaims.SignedString([]byte(config.AccessSecret))
//	assert.Error(t, err)
//}

func TestGetClaimsAndVerifyToken_UnexpectedSigningMethod(t *testing.T) {
	// Crea un token con un método de firma diferente
	claims := jwt.MapClaims{
		"id":   123,
		"type": Access,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	result, err := service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetClaimsAndVerifyToken_InvalidClaimsType(t *testing.T) {
	// Crea un token válido pero con claims que no son MapClaims
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("test_access_secret"))
	require.NoError(t, err)

	config := JWTConfig{
		AccessSecret:  "test_access_secret",
		RefreshSecret: "test_refresh_secret",
		AccessTime:    30,
		RefreshTime:   24,
	}
	service := NewJWTServiceWithConfig(config)

	// Forzamos el error de tipo de claims
	result, err := service.GetClaimsAndVerifyToken(tokenString, Access)
	assert.Error(t, err)
	assert.Nil(t, result)
}
