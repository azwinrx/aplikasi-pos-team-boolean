package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// NotificationUseCase mendefinisikan interface untuk notification business logic
type NotificationUseCase interface {
	ListNotifications(ctx context.Context, userID uint, req *dto.NotificationListRequest) (*dto.NotificationListResponse, error)
	UpdateNotificationStatus(ctx context.Context, userID uint, notificationID uint, req *dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error)
	DeleteNotification(ctx context.Context, userID uint, notificationID uint) (*dto.DeleteNotificationResponse, error)
	CreateNotification(ctx context.Context, notification *entity.Notification) error
}

// notificationUseCase implementasi dari NotificationUseCase interface
type notificationUseCase struct {
	repo   repository.NotificationRepository
	logger *zap.Logger
}

// NewNotificationUseCase membuat instance baru dari notificationUseCase
func NewNotificationUseCase(repo repository.NotificationRepository, logger *zap.Logger) NotificationUseCase {
	return &notificationUseCase{
		repo:   repo,
		logger: logger,
	}
}

// ListNotifications mengambil daftar notifikasi dengan filter dan pagination
func (u *notificationUseCase) ListNotifications(ctx context.Context, userID uint, req *dto.NotificationListRequest) (*dto.NotificationListResponse, error) {
	// Validasi input
	if userID == 0 {
		u.logger.Warn("Invalid user ID",
			zap.Uint("user_id", userID),
		)
		return nil, errors.New("user_id tidak valid")
	}

	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit 100
	}

	// Get notifications
	notifications, total, err := u.repo.GetNotificationsByUserID(
		ctx,
		userID,
		req.Status,
		req.Type,
		req.Page,
		req.Limit,
		req.SortBy,
		req.SortOrder,
	)
	if err != nil {
		u.logger.Error("Failed to get notifications",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	// Get unread count
	unreadCount, err := u.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		u.logger.Error("Failed to get unread count",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	// Convert to DTO
	notificationResponses := make([]dto.NotificationResponse, len(notifications))
	for i, n := range notifications {
		notificationResponses[i] = dto.NotificationResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			Title:     n.Title,
			Message:   n.Message,
			Type:      n.Type,
			Status:    n.Status,
			ReadedAt:  n.ReadedAt,
			Data:      n.Data,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := (int(total) + req.Limit - 1) / req.Limit
	if totalPages < 1 {
		totalPages = 1
	}

	response := &dto.NotificationListResponse{
		Data:       notificationResponses,
		Total:      int(total),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		UnreadCount: int(unreadCount),
	}

	u.logger.Info("Notifications listed successfully",
		zap.Uint("user_id", userID),
		zap.Int("count", len(notifications)),
		zap.Int64("unread", unreadCount),
	)

	return response, nil
}

// UpdateNotificationStatus mengupdate status notifikasi
func (u *notificationUseCase) UpdateNotificationStatus(ctx context.Context, userID uint, notificationID uint, req *dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error) {
	// Validasi input
	if userID == 0 || notificationID == 0 {
		u.logger.Warn("Invalid IDs",
			zap.Uint("user_id", userID),
			zap.Uint("notification_id", notificationID),
		)
		return nil, errors.New("user_id atau notification_id tidak valid")
	}

	// Validasi status
	validStatuses := map[string]bool{"new": true, "readed": true}
	if !validStatuses[req.Status] {
		u.logger.Warn("Invalid status",
			zap.String("status", req.Status),
		)
		return nil, fmt.Errorf("status '%s' tidak valid. gunakan 'new' atau 'readed'", req.Status)
	}

	// Get notification untuk verifikasi ownership
	notification, err := u.repo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		u.logger.Error("Failed to get notification",
			zap.Uint("notification_id", notificationID),
			zap.Error(err),
		)
		return nil, err
	}

	if notification == nil {
		u.logger.Warn("Notification not found",
			zap.Uint("notification_id", notificationID),
		)
		return nil, errors.New("notifikasi tidak ditemukan")
	}

	// Verifikasi ownership
	if notification.UserID != userID {
		u.logger.Warn("Unauthorized access",
			zap.Uint("user_id", userID),
			zap.Uint("notification_user_id", notification.UserID),
		)
		return nil, errors.New("anda tidak memiliki akses ke notifikasi ini")
	}

	// Update status
	if err := u.repo.UpdateNotificationStatus(ctx, notificationID, req.Status); err != nil {
		u.logger.Error("Failed to update notification status",
			zap.Uint("notification_id", notificationID),
			zap.String("status", req.Status),
			zap.Error(err),
		)
		return nil, err
	}

	response := &dto.UpdateNotificationStatusResponse{
		ID:        notificationID,
		Status:    req.Status,
		UpdatedAt: notification.UpdatedAt,
		Message:   fmt.Sprintf("Status notifikasi berhasil diubah menjadi %s", req.Status),
	}

	u.logger.Info("Notification status updated successfully",
		zap.Uint("notification_id", notificationID),
		zap.String("status", req.Status),
	)

	return response, nil
}

// DeleteNotification menghapus notifikasi
func (u *notificationUseCase) DeleteNotification(ctx context.Context, userID uint, notificationID uint) (*dto.DeleteNotificationResponse, error) {
	// Validasi input
	if userID == 0 || notificationID == 0 {
		u.logger.Warn("Invalid IDs",
			zap.Uint("user_id", userID),
			zap.Uint("notification_id", notificationID),
		)
		return nil, errors.New("user_id atau notification_id tidak valid")
	}

	// Get notification untuk verifikasi ownership
	notification, err := u.repo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		u.logger.Error("Failed to get notification",
			zap.Uint("notification_id", notificationID),
			zap.Error(err),
		)
		return nil, err
	}

	if notification == nil {
		u.logger.Warn("Notification not found",
			zap.Uint("notification_id", notificationID),
		)
		return nil, errors.New("notifikasi tidak ditemukan")
	}

	// Verifikasi ownership
	if notification.UserID != userID {
		u.logger.Warn("Unauthorized deletion attempt",
			zap.Uint("user_id", userID),
			zap.Uint("notification_user_id", notification.UserID),
		)
		return nil, errors.New("anda tidak memiliki akses untuk menghapus notifikasi ini")
	}

	// Delete notification
	if err := u.repo.DeleteNotification(ctx, notificationID); err != nil {
		u.logger.Error("Failed to delete notification",
			zap.Uint("notification_id", notificationID),
			zap.Error(err),
		)
		return nil, err
	}

	response := &dto.DeleteNotificationResponse{
		ID:      notificationID,
		Message: "Notifikasi berhasil dihapus",
	}

	u.logger.Info("Notification deleted successfully",
		zap.Uint("notification_id", notificationID),
		zap.Uint("user_id", userID),
	)

	return response, nil
}

// CreateNotification membuat notifikasi baru
func (u *notificationUseCase) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	// Validasi input
	if notification.UserID == 0 {
		u.logger.Warn("Invalid user ID for notification creation",
			zap.Uint("user_id", notification.UserID),
		)
		return errors.New("user_id tidak valid")
	}

	if notification.Title == "" {
		u.logger.Warn("Notification title is required")
		return errors.New("title notifikasi tidak boleh kosong")
	}

	// Validasi type
	validTypes := map[string]bool{"order": true, "payment": true, "system": true, "alert": true}
	if !validTypes[notification.Type] {
		u.logger.Warn("Invalid notification type",
			zap.String("type", notification.Type),
		)
		return errors.New("tipe notifikasi tidak valid. gunakan: order, payment, system, atau alert")
	}

	// Create notification
	if err := u.repo.CreateNotification(ctx, notification); err != nil {
		u.logger.Error("Failed to create notification",
			zap.Uint("user_id", notification.UserID),
			zap.Error(err),
		)
		return err
	}

	u.logger.Info("Notification created successfully",
		zap.Uint("id", notification.ID),
		zap.Uint("user_id", notification.UserID),
	)

	return nil
}
