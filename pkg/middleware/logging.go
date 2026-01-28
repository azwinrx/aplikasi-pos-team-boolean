package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware returns a Gin middleware handler for logging with Zap
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body for logging (if present)
		var requestBody []byte
		var bodySize int64
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			bodySize = int64(len(requestBody))
			// Restore body untuk handler berikutnya
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Log incoming request dengan detail lengkap
		logger.Info("incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
			zap.Int64("body_size", bodySize),
			zap.String("content_type", c.Request.Header.Get("Content-Type")),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		responseSize := int64(c.Writer.Size())

		// Determine log level based on status code
		if statusCode >= 500 {
			// Server error
			logger.Error("request completed with server error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.Int64("response_size", responseSize),
				zap.String("remote_addr", c.ClientIP()),
			)
		} else if statusCode >= 400 {
			// Client error
			logger.Warn("request completed with client error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.Int64("response_size", responseSize),
				zap.String("remote_addr", c.ClientIP()),
			)
		} else {
			// Success
			logger.Info("request completed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
				zap.Int64("response_size", responseSize),
			)
		}
	}
}
