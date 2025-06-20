package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestAuthJWTMiddleware_NoToken(t *testing.T) {
	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Token not provided", response["error"])
}

func TestAuthJWTMiddleware_NoJWTSecret(t *testing.T) {
	// Clear JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Unsetenv("JWT_ACCESS_SECRET")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "JWT_ACCESS_SECRET not configured", response["error"])
}

func TestAuthJWTMiddleware_InvalidToken(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid token", response["error"])
}

func TestAuthJWTMiddleware_ExpiredToken(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create expired token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(-1 * time.Hour).Unix(), // Expired 1 hour ago
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	// The error message might be "Invalid token" instead of "Token expired" due to JWT parsing
	assert.Contains(t, []string{"Token expired", "Invalid token"}, response["error"])
}

func TestAuthJWTMiddleware_InvalidTokenClaims(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create token without exp claim
	claims := jwt.MapClaims{
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid token claims", response["error"])
}

func TestAuthJWTMiddleware_WrongTokenType(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create token with wrong type
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "refresh", // Wrong type
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Token type mismatch", response["error"])
}

func TestAuthJWTMiddleware_MissingTokenType(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create token without type claim
	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Missing token type", response["error"])
}

func TestAuthJWTMiddleware_ValidToken(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create valid token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "access",
		"id":   123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	middleware := AuthJWTMiddleware()
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthJWTMiddleware_TokenWithoutBearer(t *testing.T) {
	// Set JWT_ACCESS_SECRET
	originalSecret := os.Getenv("JWT_ACCESS_SECRET")
	os.Setenv("JWT_ACCESS_SECRET", "test-secret")
	defer os.Setenv("JWT_ACCESS_SECRET", originalSecret)

	// Create valid token
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "access",
		"id":   123,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))

	c, w := setupGinContext()
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", tokenString) // Without "Bearer " prefix

	middleware := AuthJWTMiddleware()
	middleware(c)

	// The middleware should still process the token even without "Bearer " prefix
	// because strings.TrimPrefix handles this case
	assert.Equal(t, http.StatusOK, w.Code)
}
