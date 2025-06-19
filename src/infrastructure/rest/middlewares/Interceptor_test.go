package middlewares

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// MockResponseWriter implements gin.ResponseWriter for testing
type MockResponseWriter struct {
	*httptest.ResponseRecorder
}

func (m *MockResponseWriter) CloseNotify() <-chan bool {
	return make(chan bool)
}

func (m *MockResponseWriter) Flush() {
}

func (m *MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (m *MockResponseWriter) Size() int {
	return len(m.Body.Bytes())
}

func (m *MockResponseWriter) Status() int {
	return m.Code
}

func (m *MockResponseWriter) WriteHeaderNow() {
}

func (m *MockResponseWriter) Written() bool {
	return m.Code != 0
}

func (m *MockResponseWriter) WriteString(string) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(code int) {
	m.Code = code
}

func (m *MockResponseWriter) Pusher() http.Pusher {
	return nil
}

func TestGinBodyLogMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	// Add a test route
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test response"})
	})

	// Create test request body
	requestBody := `{"test": "data"}`

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check response body
	expectedResponse := `{"message":"test response"}`
	if !strings.Contains(w.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %s, got %s", expectedResponse, w.Body.String())
	}
}

func TestGinBodyLogMiddleware_EmptyBody(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Create a test request with empty body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Ensure the request body is properly initialized
	if req.Body == nil {
		req.Body = io.NopCloser(bytes.NewBuffer([]byte("")))
	}

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestBodyLogWriter_Write(t *testing.T) {
	// Create a mock response writer
	mockWriter := &MockResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}

	// Create bodyLogWriter
	blw := &bodyLogWriter{
		ResponseWriter: mockWriter,
		body:           bytes.NewBufferString(""),
	}

	// Test data to write
	testData := []byte("test response data")

	// Write data
	bytesWritten, err := blw.Write(testData)

	// Check no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check bytes written
	if bytesWritten != len(testData) {
		t.Errorf("Expected %d bytes written, got %d", len(testData), bytesWritten)
	}

	// Check body buffer contains the data
	if blw.body.String() != string(testData) {
		t.Errorf("Expected body to contain %s, got %s", string(testData), blw.body.String())
	}

	// Check response writer also contains the data
	if mockWriter.Body.String() != string(testData) {
		t.Errorf("Expected response writer to contain %s, got %s", string(testData), mockWriter.Body.String())
	}
}

func TestGinBodyLogMiddleware_LargeBody(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	// Add a test route
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "large body test"})
	})

	// Create a large request body (larger than the 4096 buffer)
	largeBody := strings.Repeat("a", 5000)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(largeBody))

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
