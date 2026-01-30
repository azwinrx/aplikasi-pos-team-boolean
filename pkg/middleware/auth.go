package middleware

import (
	"aplikasi-pos-team-boolean/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware returns a Gin middleware handler for JWT authentication
func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			utils.ResponseError(c, http.StatusUnauthorized, "Missing authorization header")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format",
				zap.String("auth_header", authHeader),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token and get claims
		claims, err := utils.ValidateToken(token)
		if err != nil {
			logger.Warn("Invalid token",
				zap.Error(err),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)

		logger.Debug("User authenticated",
			zap.Uint("user_id", claims.UserID),
			zap.String("email", claims.Email),
			zap.String("role", claims.Role),
		)

		c.Next()
	}
}
