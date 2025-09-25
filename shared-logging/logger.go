package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig holds configuration for logger initialization
type LoggerConfig struct {
	Level       string `yaml:"level" json:"level"`
	Format      string `yaml:"format" json:"format"`
	Output      string `yaml:"output" json:"output"`
	ServiceName string `yaml:"service_name" json:"service_name"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig(serviceName string) *LoggerConfig {
	return &LoggerConfig{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: serviceName,
	}
}

// NewLogger creates a new standardized zap logger
func NewLogger(config *LoggerConfig) (*zap.Logger, error) {
	var zapConfig zap.Config

	// Use production config as base
	zapConfig = zap.NewProductionConfig()

	// Set log level
	level := zap.InfoLevel
	switch config.Level {
	case "debug":
		level = zap.DebugLevel
		zapConfig = zap.NewDevelopmentConfig() // Use development config for debug
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Standardize encoder configuration
	zapConfig.Encoding = config.Format
	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // Uses system time in ISO8601 format
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set output paths
	zapConfig.OutputPaths = []string{config.Output}
	zapConfig.ErrorOutputPaths = []string{config.Output}

	// Build logger
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	// Add service name as a default field
	if config.ServiceName != "" {
		logger = logger.With(zap.String("service_name", config.ServiceName))
	}

	return logger, nil
}

// NewDevelopmentLogger creates a logger optimized for development
func NewDevelopmentLogger(serviceName string) (*zap.Logger, error) {
	config := &LoggerConfig{
		Level:       "debug",
		Format:      "console", // More readable in development
		Output:      "stdout",
		ServiceName: serviceName,
	}
	return NewLogger(config)
}

// NewProductionLogger creates a logger optimized for production
func NewProductionLogger(serviceName string) (*zap.Logger, error) {
	config := &LoggerConfig{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: serviceName,
	}
	return NewLogger(config)
}

// Fields provides common field helpers
type Fields struct{}

var Field Fields

// String creates a string field
func (Fields) String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an int field
func (Fields) Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Uint creates a uint field
func (Fields) Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Int64 creates an int64 field
func (Fields) Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Float64 creates a float64 field
func (Fields) Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool creates a bool field
func (Fields) Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Error creates an error field
func (Fields) Error(err error) zap.Field {
	return zap.Error(err)
}

// Duration creates a duration field
func (Fields) Duration(key string, val interface{}) zap.Field {
	switch v := val.(type) {
	case string:
		return zap.String(key, v)
	default:
		return zap.Any(key, val)
	}
}

// RequestID creates a request ID field
func (Fields) RequestID(id string) zap.Field {
	return zap.String("request_id", id)
}

// UserID creates a user ID field
func (Fields) UserID(id uint) zap.Field {
	return zap.Uint("user_id", id)
}

// TenantID creates a tenant ID field
func (Fields) TenantID(id uint) zap.Field {
	return zap.Uint("tenant_id", id)
}

// CorrelationID creates a correlation ID field
func (Fields) CorrelationID(id string) zap.Field {
	return zap.String("correlation_id", id)
}

// HTTPMethod creates an HTTP method field
func (Fields) HTTPMethod(method string) zap.Field {
	return zap.String("http_method", method)
}

// HTTPPath creates an HTTP path field
func (Fields) HTTPPath(path string) zap.Field {
	return zap.String("http_path", path)
}

// HTTPStatus creates an HTTP status code field
func (Fields) HTTPStatus(status int) zap.Field {
	return zap.Int("http_status", status)
}

// ClientIP creates a client IP field
func (Fields) ClientIP(ip string) zap.Field {
	return zap.String("client_ip", ip)
}

// UserAgent creates a user agent field
func (Fields) UserAgent(ua string) zap.Field {
	return zap.String("user_agent", ua)
}