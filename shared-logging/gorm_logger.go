package logging

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLoggerConfig holds configuration for GORM logger
type GormLoggerConfig struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// DefaultGormConfig returns default GORM logger configuration
func DefaultGormConfig() *GormLoggerConfig {
	return &GormLoggerConfig{
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// GormLogger implements gorm logger interface with standardized logging
type GormLogger struct {
	logger *zap.Logger
	config *GormLoggerConfig
}

// NewGormLogger creates a new GORM logger with standardized format
func NewGormLogger(logger *zap.Logger, config *GormLoggerConfig) gormlogger.Interface {
	if config == nil {
		config = DefaultGormConfig()
	}

	return &GormLogger{
		logger: logger,
		config: config,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.config = &GormLoggerConfig{
		LogLevel:                  level,
		SlowThreshold:             l.config.SlowThreshold,
		IgnoreRecordNotFoundError: l.config.IgnoreRecordNotFoundError,
	}
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		fields := l.extractContextFields(ctx)
		fields = append(fields, zap.Any("data", data))
		l.logger.Info("GORM info: "+msg, fields...)
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		fields := l.extractContextFields(ctx)
		fields = append(fields, zap.Any("data", data))
		l.logger.Warn("GORM warning: "+msg, fields...)
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error {
		fields := l.extractContextFields(ctx)
		fields = append(fields, zap.Any("data", data))
		l.logger.Error("GORM error: "+msg, fields...)
	}
}

// Trace logs SQL queries with standardized format
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.config.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Base fields
	fields := l.extractContextFields(ctx)
	fields = append(fields,
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows_affected", rows),
		zap.String("sql", sql),
	)

	switch {
	case err != nil && l.config.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.config.IgnoreRecordNotFoundError):
		// Log database errors
		fields = append(fields, Field.Error(err))
		l.logger.Error("Database query error", fields...)

	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= gormlogger.Warn:
		// Log slow queries
		fields = append(fields,
			zap.Duration("slow_threshold", l.config.SlowThreshold),
			zap.Bool("is_slow_query", true),
		)
		l.logger.Warn("Slow database query detected", fields...)

	case l.config.LogLevel >= gormlogger.Info:
		// Log all queries in info level
		l.logger.Debug("Database query executed", fields...)
	}
}

// extractContextFields extracts standard fields from context
func (l *GormLogger) extractContextFields(ctx context.Context) []zap.Field {
	var fields []zap.Field

	// Try to extract request ID from context
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields = append(fields, Field.RequestID(id))
		}
	}

	// Try to extract user ID from context
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(uint); ok {
			fields = append(fields, Field.UserID(id))
		}
	}

	// Try to extract tenant ID from context
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		if id, ok := tenantID.(uint); ok {
			fields = append(fields, Field.TenantID(id))
		}
	}

	// Try to extract correlation ID from context
	if correlationID := ctx.Value("correlation_id"); correlationID != nil {
		if id, ok := correlationID.(string); ok {
			fields = append(fields, Field.CorrelationID(id))
		}
	}

	return fields
}