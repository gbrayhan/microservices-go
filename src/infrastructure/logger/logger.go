package infrastructure

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	Log *zap.Logger
}

func NewLogger() (*Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core)

	return &Logger{Log: logger}, nil
}

// NewDevelopmentLogger creates a logger for development with more debug information
func NewDevelopmentLogger() (*Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{Log: logger}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Log.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.Log.Panic(msg, fields...)
}
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
}

// SetupGinWithZapLogger configures Gin to use the Zap logger
func (l *Logger) SetupGinWithZapLogger() {
	// Configure Gin to use release mode by default
	gin.SetMode(gin.ReleaseMode)

	// Create a custom writer that uses Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// SetupGinWithZapLoggerInDevelopment configures Gin to use the Zap logger in development mode
func (l *Logger) SetupGinWithZapLoggerInDevelopment() {
	// Configure Gin to use debug mode in development
	gin.SetMode(gin.DebugMode)

	// Create a custom writer that uses Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// SetupGinWithZapLoggerWithMode configures Gin to use the Zap logger with a specific mode
func (l *Logger) SetupGinWithZapLoggerWithMode(mode string) {
	// Configure Gin to use the specified mode
	gin.SetMode(mode)

	// Create a custom writer that uses Zap
	gin.DefaultWriter = &ZapWriter{logger: l.Log}
	gin.DefaultErrorWriter = &ZapErrorWriter{logger: l.Log}
}

// ZapWriter implements io.Writer for use with Gin
type ZapWriter struct {
	logger *zap.Logger
}

func (w *ZapWriter) Write(p []byte) (n int, err error) {
	w.logger.Info("Gin-log", zap.String("message", string(p)))
	return len(p), nil
}

// ZapErrorWriter implements io.Writer for Gin errors
type ZapErrorWriter struct {
	logger *zap.Logger
}

func (w *ZapErrorWriter) Write(p []byte) (n int, err error) {
	w.logger.Error("Gin-error", zap.String("error", string(p)))
	return len(p), nil
}

func (l *Logger) GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		l.Log.Info("HTTP request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()), zap.Duration("latency", latency), zap.String("client_ip", c.ClientIP()))
	}
}

type GormZapLogger struct {
	zap    *zap.SugaredLogger
	config gormlogger.Config
}

func NewGormLogger(base *zap.Logger) *GormZapLogger {
	sugar := base.Sugar()
	return &GormZapLogger{
		zap: sugar,
		config: gormlogger.Config{
			SlowThreshold:             time.Second, // threshold to highlight slow queries
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true, // do not log "record not found"
			Colorful:                  false,
		},
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newCfg := l.config
	newCfg.LogLevel = level
	return &GormZapLogger{zap: l.zap, config: newCfg}
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		l.zap.Infof(msg, data...)
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		l.zap.Warnf(msg, data...)
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error &&
		(!l.config.IgnoreRecordNotFoundError || msg != gormlogger.ErrRecordNotFound.Error()) {
		l.zap.Errorf(msg, data...)
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	if err != nil {
		if l.config.IgnoreRecordNotFoundError && errors.Is(err, gormlogger.ErrRecordNotFound) {
			return
		}
		if l.config.LogLevel >= gormlogger.Error {
			sql, rows := fc()
			l.zap.Errorf("Error: %v | %.3fms | rows:%d | %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		return
	}

	if elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.zap.Warnf("SLOW ≥ %s | %.3fms | rows:%d | %s", l.config.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
