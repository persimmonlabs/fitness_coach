package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery creates a middleware that recovers from panics and logs the stack trace
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()

				// Get request ID if available
				requestID := GetRequestID(c)

				// Log the panic with full context
				logger.Error("Panic recovered",
					zap.String("request_id", requestID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.Any("error", err),
					zap.String("stack", string(stack)),
				)

				// Return error response
				errorMessage := "Internal server error"
				if gin.Mode() == gin.DebugMode {
					errorMessage = fmt.Sprintf("Panic: %v", err)
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":      errorMessage,
					"request_id": requestID,
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
