package adaptor

import (
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthAdaptor menangani semua request HTTP untuk authentication
type AuthAdaptor struct {
	authUsecase usecase.AuthUseCase
	logger      *zap.Logger
}

// NewAuthAdaptor membuat instance baru dari AuthAdaptor
func NewAuthAdaptor(authUsecase usecase.AuthUseCase, logger *zap.Logger) *AuthAdaptor {
	return &AuthAdaptor{
		authUsecase: authUsecase,
		logger:      logger,
	}
}

// Register menangani request POST /auth/register
func (h *AuthAdaptor) Register(c *gin.Context) {
	h.logger.Debug("Register handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.RegisterRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for register",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("Registration failed",
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		statusCode := http.StatusBadRequest
		if err.Error() == "database error" {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Registration successful",
		zap.String("email", req.Email),
		zap.Uint("user_id", response.ID),
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}

// Login menangani request POST /auth/login
func (h *AuthAdaptor) Login(c *gin.Context) {
	h.logger.Debug("Login handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for login",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("Login failed",
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Login successful",
		zap.String("email", req.Email),
		zap.Uint("user_id", response.ID),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login successful",
		"data":    response,
	})
}

// GetUserByID menangani request GET /auth/user/:id
func (h *AuthAdaptor) GetUserByID(c *gin.Context) {
	h.logger.Debug("GetUserByID handler called", zap.String("client_ip", c.ClientIP()))

	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		h.logger.Warn("Missing user ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User ID is required",
			"data":    nil,
		})
		return
	}

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid user ID format",
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Warn("Get user failed",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		statusCode := http.StatusNotFound
		if err.Error() == "database error" {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Get user successful",
		zap.Uint("user_id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User retrieved successfully",
		"data":    response,
	})
}

// DeleteUser menangani request DELETE /auth/user/:id
func (h *AuthAdaptor) DeleteUser(c *gin.Context) {
	h.logger.Debug("DeleteUser handler called", zap.String("client_ip", c.ClientIP()))

	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		h.logger.Warn("Missing user ID parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User ID is required",
			"data":    nil,
		})
		return
	}

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		h.logger.Warn("Invalid user ID format",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid user ID format",
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		h.logger.Warn("Delete user failed",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		statusCode := http.StatusBadRequest
		if err.Error() == "database error" {
			statusCode = http.StatusInternalServerError
		} else if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Delete user successful",
		zap.Uint("user_id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}

// CheckEmail menangani request POST /auth/check-email
func (h *AuthAdaptor) CheckEmail(c *gin.Context) {
	h.logger.Debug("CheckEmail handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.CheckEmailRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for check email",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.CheckEmail(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("CheckEmail failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("CheckEmail completed",
		zap.String("email", req.Email),
		zap.Bool("exists", response.Exists),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}

// SendOTP menangani request POST /auth/send-otp
func (h *AuthAdaptor) SendOTP(c *gin.Context) {
	h.logger.Debug("SendOTP handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.SendOTPRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for send OTP",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.SendOTP(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("SendOTP failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("OTP sent successfully",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}

// ValidateOTP menangani request POST /auth/validate-otp
func (h *AuthAdaptor) ValidateOTP(c *gin.Context) {
	h.logger.Debug("ValidateOTP handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.ValidateOTPRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for validate OTP",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.ValidateOTP(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("ValidateOTP failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("OTP validated successfully",
		zap.String("email", req.Email),
		zap.String("purpose", req.Purpose),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}

// ResetPassword menangani request POST /auth/reset-password
func (h *AuthAdaptor) ResetPassword(c *gin.Context) {
	h.logger.Debug("ResetPassword handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.ResetPasswordRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for reset password",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Call usecase
	response, err := h.authUsecase.ResetPassword(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("ResetPassword failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		statusCode := http.StatusBadRequest
		if err.Error() == "database error" {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Password reset successfully",
		zap.String("email", req.Email),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": response.Message,
		"data":    response,
	})
}
