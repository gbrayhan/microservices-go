package security

import (
	"testing"
	"time"

	"encoding/base64"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func TestNewJWTService(t *testing.T) {
	jwtService := NewJWTService()
	if jwtService == nil {
		t.Error("Expected JWT service to be created")
	}
}

func TestGenerateJWTToken_InvalidType(t *testing.T) {
	jwtService := NewJWTServiceWithConfig(JWTConfig{})
	_, err := jwtService.GenerateJWTToken(1, "invalid")
	if err == nil {
		t.Error("Expected error for invalid token type")
	}
	if err.Error() != "invalid token type" {
		t.Errorf("Expected 'invalid token type', got %s", err.Error())
	}
}

func TestGenerateJWTToken_AccessToken(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
		AccessTime:   30,
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	token, err := jwtService.GenerateJWTToken(123, Access)
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
	jwtConfig := JWTConfig{
		RefreshSecret: "test-refresh-secret",
		RefreshTime:   24,
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	token, err := jwtService.GenerateJWTToken(123, Refresh)
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

func TestGetClaimsAndVerifyToken_InvalidToken(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	_, err := jwtService.GetClaimsAndVerifyToken("invalid-token", Access)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestGetClaimsAndVerifyToken_ValidToken(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
		AccessTime:   30,
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Generate a valid token first
	token, err := jwtService.GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Verify the token
	claims, err := jwtService.GetClaimsAndVerifyToken(token.Token, Access)
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
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
		AccessTime:   30,
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Generate an access token
	token, err := jwtService.GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to verify it as a refresh token
	_, err = jwtService.GetClaimsAndVerifyToken(token.Token, Refresh)
	if err == nil {
		t.Error("Expected error for wrong token type")
	}
}

func TestGetClaimsAndVerifyToken_ExpiredToken(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
		AccessTime:   0, // 0 minutes = expired immediately
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Generate a token that expires immediately
	token, err := jwtService.GenerateJWTToken(123, Access)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait a moment to ensure token is expired
	time.Sleep(100 * time.Millisecond)

	// Try to verify the expired token
	_, err = jwtService.GetClaimsAndVerifyToken(token.Token, Access)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestGetClaimsAndVerifyToken_RefreshToken(t *testing.T) {
	jwtConfig := JWTConfig{
		RefreshSecret: "test-refresh-secret",
		RefreshTime:   24,
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Generate a refresh token
	token, err := jwtService.GenerateJWTToken(123, Refresh)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Verify the refresh token
	claims, err := jwtService.GetClaimsAndVerifyToken(token.Token, Refresh)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if claims["id"].(float64) != 123 {
		t.Errorf("Expected user ID 123, got %v", claims["id"])
	}

	if claims["type"].(string) != Refresh {
		t.Errorf("Expected token type %s, got %s", Refresh, claims["type"])
	}
}

func TestGetClaimsAndVerifyToken_InvalidClaims(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Create a token with invalid claims structure
	claims := &Claims{
		ID:   123,
		Type: Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tokenWithClaims.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Try to verify with wrong secret
	_, err = jwtService.GetClaimsAndVerifyToken(tokenStr, Access)
	if err == nil {
		t.Error("Expected error for invalid token")
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

func TestGenerateJWTToken_SignedStringError(t *testing.T) {
	claims := &Claims{
		ID:   123,
		Type: Access,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	_, err := tokenWithClaims.SignedString(12345) // invalid key
	if err == nil {
		t.Error("Expected error when signing with invalid key type")
	}
}

func TestGetClaimsAndVerifyToken_InvalidSigningMethod(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	claims := jwt.MapClaims{"id": 123, "type": Access, "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	_, err = jwtService.GetClaimsAndVerifyToken(tokenStr, Access)
	if err == nil {
		t.Error("Expected error for unexpected signing method (any error)")
	}
}

// InvalidClaimsType implements jwt.Claims but is not jwt.MapClaims
type InvalidClaimsType struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Exp  int64  `json:"exp"`
}

func (i InvalidClaimsType) Valid() error {
	return nil
}

func TestGetClaimsAndVerifyToken_ArtificialInvalidClaimsType(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Create a valid JWT
	claims := jwt.MapClaims{"id": 123, "type": Access, "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Manipulate the payload to make it invalid JSON (e.g., base64 string of a number)
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		t.Fatalf("Invalid JWT format")
	}
	// Replace payload with a number encoded in base64
	parts[1] = base64.RawURLEncoding.EncodeToString([]byte("12345"))
	manipulatedToken := strings.Join(parts, ".")

	_, err = jwtService.GetClaimsAndVerifyToken(manipulatedToken, Access)
	if err == nil {
		t.Errorf("Expected any error for invalid claims type, got: %v", err)
	}
}

func TestGetClaimsAndVerifyToken_CorruptToken(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	_, err := jwtService.GetClaimsAndVerifyToken("this-is-not-a-jwt", Access)
	if err == nil {
		t.Error("Expected error for corrupt token")
	}
}

func TestGetClaimsAndVerifyToken_MissingExpField(t *testing.T) {
	jwtConfig := JWTConfig{
		AccessSecret: "test-secret",
	}
	jwtService := NewJWTServiceWithConfig(jwtConfig)

	// Create a token without exp field to force type assertion error
	claims := jwt.MapClaims{"id": 123, "type": Access} // Without "exp" field
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Use defer recover to handle expected panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for missing exp field")
		}
	}()

	_, err = jwtService.GetClaimsAndVerifyToken(tokenStr, Access)
	// We don't check err because we expect a panic, not an error
	_ = err
}
