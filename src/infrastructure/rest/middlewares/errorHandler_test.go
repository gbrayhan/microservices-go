package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	"github.com/gin-gonic/gin"
)

func TestErrorHandler_NoErrors(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a NotFound error
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"record not found"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_ValidationError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a ValidationError
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"validation error"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_RepositoryError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a RepositoryError
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.RepositoryError)
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"error in repository operation"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_NotAuthenticatedError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a NotAuthenticated error
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotAuthenticated)
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"not Authenticated"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_NotAuthorizedError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a NotAuthorized error
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType(domainErrors.NotAuthorized)
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"not authorized"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_UnknownErrorType(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates an unknown error type
	router.GET("/test", func(c *gin.Context) {
		appErr := domainErrors.NewAppErrorWithType("UnknownErrorType")
		_ = c.Error(appErr)
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"Internal Server Error"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestErrorHandler_NonAppError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(ErrorHandler())

	// Add a test route that generates a regular error
	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("regular error"))
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// Check response body
	expectedBody := `{"error":"Internal Server Error"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}
