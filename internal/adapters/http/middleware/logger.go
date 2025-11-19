package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger creates a middleware that logs HTTP requests using zap structured logging
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Get request ID from context (set by RequestID middleware)
		requestID, _ := c.Get("requestID")

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Build log fields
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		}

		// Add request ID if available
		if requestID != nil {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("userID"); exists {
			fields = append(fields, zap.String("user_id", userID.(string)))
		}

		// Log errors at ERROR level, otherwise at INFO level
		if len(c.Errors) > 0 {
			// Log errors
			for _, err := range c.Errors {
				logger.Error("HTTP request error",
					append(fields, zap.Error(err.Err))...,
				)
			}
		} else if status >= 500 {
			// Server errors
			logger.Error("HTTP request completed with server error", fields...)
		} else if status >= 400 {
			// Client errors
			logger.Warn("HTTP request completed with client error", fields...)
		} else {
			// Success
			logger.Info("HTTP request completed", fields...)
		}
	}
}
