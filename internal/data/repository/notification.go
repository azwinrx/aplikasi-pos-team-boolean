package repository

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NotificationRepository mendefinisikan interface untuk notification operations
type NotificationRepository interface {
	GetNotificationsByUserID(ctx context.Context, userID uint, status string, notificationType string, page int, limit int, sortBy string, sortOrder string) ([]entity.Notification, int64, error)
	GetNotificationByID(ctx context.Context, id uint) (*entity.Notification, error)
	CreateNotification(ctx context.Context, notification *entity.Notification) error
	UpdateNotificationStatus(ctx context.Context, id uint, status string) error
	DeleteNotification(ctx context.Context, id uint) error
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	DeleteOldNotifications(ctx context.Context, days int) error
}

// notificationRepository implementasi dari NotificationRepository interface
type notificationRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewNotificationRepository membuat instance baru dari notificationRepository
func NewNotificationRepository(db *gorm.DB, logger *zap.Logger) NotificationRepository {
	return &notificationRepository{
		db:     db,
		logger: logger,
	}
}

// GetNotificationsByUserID mengambil notifikasi berdasarkan user ID dengan filter
func (r *notificationRepository) GetNotificationsByUserID(ctx context.Context, userID uint, status string, notificationType string, page int, limit int, sortBy string, sortOrder string) ([]entity.Notification, int64, error) {
	var notifications []entity.Notification
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Filter by status
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter by type
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// Get total count
	if err := query.Model(&entity.Notification{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count notifications",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, 0, err
	}

	// Set default values
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Sort and pagination
	sortColumn := sortBy
	if sortOrder == "asc" {
		query = query.Order(sortColumn + " ASC")
	} else {
		query = query.Order(sortColumn + " DESC")
	}

	if err := query.Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("No notifications found",
				zap.Uint("user_id", userID),
			)
			return []entity.Notification{}, 0, nil
		}

		r.logger.Error("Failed to get notifications",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, 0, err
	}

	r.logger.Info("Notifications retrieved successfully",
		zap.Uint("user_id", userID),
		zap.Int("count", len(notifications)),
		zap.Int64("total", total),
	)

	return notifications, total, nil
}

// GetNotificationByID mengambil notifikasi berdasarkan ID
func (r *notificationRepository) GetNotificationByID(ctx context.Context, id uint) (*entity.Notification, error) {
	var notification entity.Notification

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("Notification not found",
				zap.Uint("id", id),
			)
			return nil, nil
		}

		r.logger.Error("Failed to get notification",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return &notification, nil
}

// CreateNotification membuat notifikasi baru
func (r *notificationRepository) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	if notification.Status == "" {
		notification.Status = "new"
	}

	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		r.logger.Error("Failed to create notification",
			zap.Uint("user_id", notification.UserID),
			zap.String("title", notification.Title),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Notification created successfully",
		zap.Uint("id", notification.ID),
		zap.Uint("user_id", notification.UserID),
		zap.String("type", notification.Type),
	)

	return nil
}

// UpdateNotificationStatus mengupdate status notifikasi
func (r *notificationRepository) UpdateNotificationStatus(ctx context.Context, id uint, status string) error {
	update := map[string]interface{}{
		"status": status,
	}

	// Jika status diubah ke "readed", set readed_at
	if status == "readed" {
		update["readed_at"] = time.Now()
	}

	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Updates(update).Error; err != nil {
		r.logger.Error("Failed to update notification status",
			zap.Uint("id", id),
			zap.String("status", status),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Notification status updated successfully",
		zap.Uint("id", id),
		zap.String("status", status),
	)

	return nil
}

// DeleteNotification menghapus notifikasi (soft delete)
func (r *notificationRepository) DeleteNotification(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Notification{}, id).Error; err != nil {
		r.logger.Error("Failed to delete notification",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Notification deleted successfully",
		zap.Uint("id", id),
	)

	return nil
}

// GetUnreadCount menghitung jumlah notifikasi yang belum dibaca
func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND status = ?", userID, "new").Count(&count).Error; err != nil {
		r.logger.Error("Failed to count unread notifications",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return 0, err
	}

	return count, nil
}

// DeleteOldNotifications menghapus notifikasi yang sudah lama (>days)
func (r *notificationRepository) DeleteOldNotifications(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)

	if err := r.db.WithContext(ctx).Where("created_at < ?", cutoffDate).Delete(&entity.Notification{}).Error; err != nil {
		r.logger.Error("Failed to delete old notifications",
			zap.Int("days", days),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Old notifications deleted successfully",
		zap.Int("days", days),
	)

	return nil
}
