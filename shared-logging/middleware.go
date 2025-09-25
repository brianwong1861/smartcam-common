package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPLogger creates a Gin middleware that logs HTTP requests with standardized format
func HTTPLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response data
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// Standard log fields
		fields := []zapcore.Field{
			Field.HTTPStatus(status),
			Field.HTTPMethod(method),
			Field.HTTPPath(path),
			Field.ClientIP(clientIP),
			zap.Duration("latency", latency),
			Field.UserAgent(userAgent),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Add request ID if available
		if requestID, exists := c.Get("request_id"); exists {
			if id, ok := requestID.(string); ok {
				fields = append(fields, Field.RequestID(id))
			}
		}

		// Add user context if available
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(uint); ok {
				fields = append(fields, Field.UserID(id))
			}
		}

		// Add tenant context if available
		if tenantID, exists := c.Get("tenant_id"); exists {
			if id, ok := tenantID.(uint); ok {
				fields = append(fields, Field.TenantID(id))
			}
		}

		// Add correlation ID if available
		if correlationID := c.GetHeader("X-Correlation-ID"); correlationID != "" {
			fields = append(fields, Field.CorrelationID(correlationID))
		}

		// Log based on status code
		if status >= 500 {
			if len(c.Errors) > 0 {
				// Log server errors with error details
				for _, err := range c.Errors {
					logger.Error("HTTP request server error", append(fields, Field.Error(err))...)
				}
			} else {
				logger.Error("HTTP request server error", fields...)
			}
		} else if status >= 400 {
			if len(c.Errors) > 0 {
				// Log client errors with error details
				for _, err := range c.Errors {
					logger.Warn("HTTP request client error", append(fields, Field.Error(err))...)
				}
			} else {
				logger.Warn("HTTP request client error", fields...)
			}
		} else {
			logger.Info("HTTP request completed", fields...)
		}
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Recovery middleware creates a standardized panic recovery handler
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		fields := []zapcore.Field{
			zap.Any("panic", recovered),
			Field.HTTPPath(c.Request.URL.Path),
			Field.HTTPMethod(c.Request.Method),
			Field.ClientIP(c.ClientIP()),
		}

		// Add request ID if available
		if requestID, exists := c.Get("request_id"); exists {
			if id, ok := requestID.(string); ok {
				fields = append(fields, Field.RequestID(id))
			}
		}

		logger.Error("HTTP request panic recovered", fields...)

		c.JSON(500, gin.H{
			"error":      "Internal server error",
			"request_id": c.GetString("request_id"),
			"timestamp":  time.Now().Format(time.RFC3339),
		})
	})
}

// CORS middleware for handling cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID, X-Correlation-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}