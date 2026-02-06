package adaptor

import (
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NotificationAdaptor menangani HTTP requests untuk notification
type NotificationAdaptor struct {
	usecase usecase.NotificationUseCase
	logger  *zap.Logger
}

// NewNotificationAdaptor membuat instance baru dari NotificationAdaptor
func NewNotificationAdaptor(uc usecase.NotificationUseCase, logger *zap.Logger) *NotificationAdaptor {
	return &NotificationAdaptor{
		usecase: uc,
		logger:  logger,
	}
}

// ListNotifications menghandle GET /notifications
func (a *NotificationAdaptor) ListNotifications(c *gin.Context) {
	// Get user ID dari context
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("User ID not found in context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	// Convert userID to uint
	uid, ok := userID.(uint)
	if !ok {
		// Try to convert from string atau other types
		if uidStr, ok := userID.(string); ok {
			uidInt, err := strconv.ParseUint(uidStr, 10, 32)
			if err != nil {
				a.logger.Error("Failed to parse user ID",
					zap.Any("user_id", userID),
					zap.Error(err),
				)
				utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
				return
			}
			uid = uint(uidInt)
		} else {
			a.logger.Error("Invalid user ID type",
				zap.Any("user_id", userID),
			)
			utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
			return
		}
	}

	// Parse query parameters
	var req dto.NotificationListRequest
	req.UserID = uid

	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		req.Page = page
	}

	if limitStr := c.DefaultQuery("limit", "10"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}
		req.Limit = limit
	}

	// Parse filters
	req.Status = c.DefaultQuery("status", "")
	req.Type = c.DefaultQuery("type", "")
	req.SortBy = c.DefaultQuery("sort_by", "created_at")
	req.SortOrder = c.DefaultQuery("sort_order", "desc")

	// Call usecase
	response, err := a.usecase.ListNotifications(c.Request.Context(), uid, &req)
	if err != nil {
		a.logger.Error("Failed to list notifications",
			zap.Uint("user_id", uid),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil daftar notifikasi: "+err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Daftar notifikasi berhasil diambil", response)
}

// UpdateNotificationStatus menghandle PUT /notifications/:id/status
func (a *NotificationAdaptor) UpdateNotificationStatus(c *gin.Context) {
	// Get user ID dari context
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("User ID not found in context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	// Convert userID to uint
	uid, ok := userID.(uint)
	if !ok {
		// Try to convert from string atau other types
		if uidStr, ok := userID.(string); ok {
			uidInt, err := strconv.ParseUint(uidStr, 10, 32)
			if err != nil {
				a.logger.Error("Failed to parse user ID",
					zap.Any("user_id", userID),
					zap.Error(err),
				)
				utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
				return
			}
			uid = uint(uidInt)
		} else {
			a.logger.Error("Invalid user ID type",
				zap.Any("user_id", userID),
			)
			utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
			return
		}
	}

	// Get notification ID dari parameter
	idStr := c.Param("id")
	notificationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		a.logger.Error("Invalid notification ID",
			zap.String("id", idStr),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "ID notifikasi tidak valid")
		return
	}

	// Parse request body
	var req dto.UpdateNotificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("Failed to bind request body",
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Request tidak valid: "+err.Error())
		return
	}

	// Call usecase
	response, err := a.usecase.UpdateNotificationStatus(
		c.Request.Context(),
		uid,
		uint(notificationID),
		&req,
	)
	if err != nil {
		a.logger.Error("Failed to update notification status",
			zap.Uint("user_id", uid),
			zap.Uint("notification_id", uint(notificationID)),
			zap.Error(err),
		)

		// Return different status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "notifikasi tidak ditemukan" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "anda tidak memiliki akses ke notifikasi ini" {
			statusCode = http.StatusForbidden
		}

		utils.ResponseError(c.Writer, statusCode, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Status notifikasi berhasil diubah", response)
}

// DeleteNotification menghandle DELETE /notifications/:id
func (a *NotificationAdaptor) DeleteNotification(c *gin.Context) {
	// Get user ID dari context
	userID, exists := c.Get("user_id")
	if !exists {
		a.logger.Warn("User ID not found in context")
		utils.ResponseError(c.Writer, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}

	// Convert userID to uint
	uid, ok := userID.(uint)
	if !ok {
		// Try to convert from string atau other types
		if uidStr, ok := userID.(string); ok {
			uidInt, err := strconv.ParseUint(uidStr, 10, 32)
			if err != nil {
				a.logger.Error("Failed to parse user ID",
					zap.Any("user_id", userID),
					zap.Error(err),
				)
				utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
				return
			}
			uid = uint(uidInt)
		} else {
			a.logger.Error("Invalid user ID type",
				zap.Any("user_id", userID),
			)
			utils.ResponseError(c.Writer, http.StatusUnauthorized, "User ID tidak valid")
			return
		}
	}

	// Get notification ID dari parameter
	idStr := c.Param("id")
	notificationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		a.logger.Error("Invalid notification ID",
			zap.String("id", idStr),
			zap.Error(err),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "ID notifikasi tidak valid")
		return
	}

	// Call usecase
	response, err := a.usecase.DeleteNotification(
		c.Request.Context(),
		uid,
		uint(notificationID),
	)
	if err != nil {
		a.logger.Error("Failed to delete notification",
			zap.Uint("user_id", uid),
			zap.Uint("notification_id", uint(notificationID)),
			zap.Error(err),
		)

		// Return different status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "notifikasi tidak ditemukan" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "anda tidak memiliki akses untuk menghapus notifikasi ini" {
			statusCode = http.StatusForbidden
		}

		utils.ResponseError(c.Writer, statusCode, err.Error())
		return
	}

	utils.ResponseSuccess(c.Writer, http.StatusOK, "Notifikasi berhasil dihapus", response)
}
