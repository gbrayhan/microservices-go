package security

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestGenerateJWTToken_InvalidType(t *testing.T) {
	_, err := GenerateJWTToken(1, "invalid")
	if err == nil {
		t.Error("Expected error for invalid token type")
	}
	if err.Error() != "invalid token type" {
		t.Errorf("Expected 'invalid token type', got %s", err.Error())
	}
}

func TestGenerateJWTToken_MissingEnvVars(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("JWT_ACCESS_SECRET")
	os.Unsetenv("JWT_ACCESS_TIME_MINUTE")

	_, err := GenerateJWTToken(1, Access)
	if err == nil {
		t.Error("Expected error for missing environment variables")
	}
	if err.Error() != "missing token environment variables" {
		t.Errorf("Expected 'missing token environment variables', got %s", err.Error())
	}
}

func TestGenerateJWTToken_InvalidExpirationTime(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "invalid")

	_, err := GenerateJWTToken(1, Access)
	if err == nil {
		t.Error("Expected error for invalid expiration time")
	}
}

func TestGenerateJWTToken_AccessToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "30")

	token, err := GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token.TokenType != Access {
		t.Errorf("Expected token type %s, got %s", Access, token.TokenType)
	}

	if token.Token == "" {
		t.Error("Expected non-empty token")
	}

	// Verify token is not expired
	if time.Now().After(token.ExpirationTime) {
		t.Error("Token should not be expired")
	}
}

func TestGenerateJWTToken_RefreshToken(t *testing.T) {
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")
	os.Setenv("JWT_REFRESH_TIME_HOUR", "24")

	token, err := GenerateJWTToken(123, Refresh)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if token.TokenType != Refresh {
		t.Errorf("Expected token type %s, got %s", Refresh, token.TokenType)
	}

	if token.Token == "" {
		t.Error("Expected non-empty token")
	}

	// Verify token is not expired
	if time.Now().After(token.ExpirationTime) {
		t.Error("Token should not be expired")
	}
}

func TestGetClaimsAndVerifyToken_MissingEnvVar(t *testing.T) {
	os.Unsetenv("JWT_ACCESS_SECRET")

	_, err := GetClaimsAndVerifyToken("test-token", Access)
	if err == nil {
		t.Error("Expected error for missing environment variable")
	}
}

func TestGetClaimsAndVerifyToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")

	_, err := GetClaimsAndVerifyToken("invalid-token", Access)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestGetClaimsAndVerifyToken_ValidToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")

	// Generate a valid token first
	token, err := GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Verify the token
	claims, err := GetClaimsAndVerifyToken(token.Token, Access)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if claims["id"].(float64) != 123 {
		t.Errorf("Expected user ID 123, got %v", claims["id"])
	}

	if claims["type"].(string) != Access {
		t.Errorf("Expected token type %s, got %s", Access, claims["type"])
	}
}

func TestGetClaimsAndVerifyToken_WrongTokenType(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")

	// Generate an access token
	token, err := GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to verify it as a refresh token
	_, err = GetClaimsAndVerifyToken(token.Token, Refresh)
	if err == nil {
		t.Error("Expected error for wrong token type")
	}
}

func TestGetClaimsAndVerifyToken_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	os.Setenv("JWT_ACCESS_TIME_MINUTE", "0") // 0 minutes = expired immediately

	// Generate a token that expires immediately
	token, err := GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a moment to ensure token is expired
	time.Sleep(100 * time.Millisecond)

	// Try to verify the expired token
	_, err = GetClaimsAndVerifyToken(token.Token, Access)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestClaims_Structure(t *testing.T) {
	claims := &Claims{
		ID:   123,
		Type: Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	if claims.ID != 123 {
		t.Errorf("Expected ID 123, got %d", claims.ID)
	}

	if claims.Type != Access {
		t.Errorf("Expected type %s, got %s", Access, claims.Type)
	}
}

func TestAppToken_Structure(t *testing.T) {
	now := time.Now()
	token := &AppToken{
		Token:          "test-token",
		TokenType:      Access,
		ExpirationTime: now.Add(time.Hour),
	}

	if token.Token != "test-token" {
		t.Errorf("Expected token 'test-token', got %s", token.Token)
	}

	if token.TokenType != Access {
		t.Errorf("Expected token type %s, got %s", Access, token.TokenType)
	}

	if !token.ExpirationTime.After(now) {
		t.Error("Expected expiration time to be in the future")
	}
}
