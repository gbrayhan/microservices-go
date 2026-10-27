package infrastructure

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

func TestNewLogger(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.NotNil(t, l.Log)

	// Test methods don't panic
	l.Info("test info message", zap.String("key", "value"))
	l.Warn("test warn message")
	l.Error("test error message")
	l.Debug("test debug message")
}

func TestNewDevelopmentLogger(t *testing.T) {
	l, err := NewDevelopmentLogger()
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.NotNil(t, l.Log)

	l.Debug("test dev debug message")
	l.Info("test dev info message")
}

func TestSetupGinWithZapLogger(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)

	l.SetupGinWithZapLogger()
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
	assert.NotNil(t, gin.DefaultWriter)
	assert.NotNil(t, gin.DefaultErrorWriter)

	l.SetupGinWithZapLoggerInDevelopment()
	assert.Equal(t, gin.DebugMode, gin.Mode())

	l.SetupGinWithZapLoggerWithMode(gin.TestMode)
	assert.Equal(t, gin.TestMode, gin.Mode())
}

func TestZapWriters(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)

	writer := &ZapWriter{logger: l.Log}
	n, err := writer.Write([]byte("test log line"))
	assert.NoError(t, err)
	assert.Equal(t, len("test log line"), n)

	errWriter := &ZapErrorWriter{logger: l.Log}
	n, err = errWriter.Write([]byte("test error log line"))
	assert.NoError(t, err)
	assert.Equal(t, len("test error log line"), n)
}

func TestGinZapLogger_Middleware(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(l.GinZapLogger())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestGormZapLogger(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)

	gormLog := NewGormLogger(l.Log)
	assert.NotNil(t, gormLog)
	assert.Equal(t, gormlogger.Error, gormLog.config.LogLevel)

	newLog := gormLog.LogMode(gormlogger.Info)
	assert.NotNil(t, newLog)

	ctx := context.Background()

	// Info and Warn under default Error level (should not log)
	gormLog.Info(ctx, "info message: %s", "detail")
	gormLog.Warn(ctx, "warn message: %s", "detail")
	gormLog.Error(ctx, "error message: %s", "detail")

	// Test with Info level
	infoLogger := gormLog.LogMode(gormlogger.Info)
	infoGorm, ok := infoLogger.(*GormZapLogger)
	assert.True(t, ok)
	infoGorm.Info(ctx, "info message: %s", "detail")
	infoGorm.Warn(ctx, "warn message: %s", "detail")
	infoGorm.Error(ctx, "error message: %s", "detail")
}

func TestGormZapLogger_Trace(t *testing.T) {
	l, err := NewLogger()
	assert.NoError(t, err)

	gormLog := NewGormLogger(l.Log)
	ctx := context.Background()

	fc := func() (string, int64) {
		return "SELECT 1;", 1
	}

	// Normal query without error, fast
	gormLog.Trace(ctx, time.Now(), fc, nil)

	// Record not found error with IgnoreRecordNotFoundError = true (default)
	gormLog.Trace(ctx, time.Now(), fc, gormlogger.ErrRecordNotFound)

	// Real error
	realErr := errors.New("db connection lost")
	gormLog.Trace(ctx, time.Now(), fc, realErr)

	// Slow query simulation with warn level
	slowLog := gormLog.LogMode(gormlogger.Warn).(*GormZapLogger)
	slowLog.config.SlowThreshold = 10 * time.Millisecond
	slowStart := time.Now().Add(-50 * time.Millisecond)
	slowLog.Trace(ctx, slowStart, fc, nil)
}
