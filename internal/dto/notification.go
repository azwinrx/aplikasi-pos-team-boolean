package dto

import "time"

// NotificationListRequest untuk filter notifikasi
type NotificationListRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Status    string `json:"status"` // new, readed, all
	Type      string `json:"type"`   // order, payment, system, alert
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sort_by"`    // created_at, status
	SortOrder string `json:"sort_order"` // asc, desc
}

// NotificationResponse merepresentasikan response notifikasi
type NotificationResponse struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	ReadedAt  *time.Time `json:"readed_at,omitempty"`
	Data      string     `json:"data,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NotificationListResponse untuk list notifikasi dengan pagination
type NotificationListResponse struct {
	Data        []NotificationResponse `json:"data"`
	Total       int                    `json:"total"`
	Page        int                    `json:"page"`
	Limit       int                    `json:"limit"`
	TotalPages  int                    `json:"total_pages"`
	UnreadCount int                    `json:"unread_count"`
}

// UpdateNotificationStatusRequest untuk update status notifikasi
type UpdateNotificationStatusRequest struct {
	NotificationID uint   `json:"notification_id" binding:"required"`
	Status         string `json:"status" binding:"required,oneof=new readed"`
}

// UpdateNotificationStatusResponse untuk response update status
type UpdateNotificationStatusResponse struct {
	ID        uint       `json:"id"`
	Status    string     `json:"status"`
	ReadedAt  *time.Time `json:"readed_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	Message   string     `json:"message"`
}

// DeleteNotificationRequest untuk delete notifikasi
type DeleteNotificationRequest struct {
	NotificationID uint `json:"notification_id" binding:"required"`
}

// DeleteNotificationResponse untuk response delete
type DeleteNotificationResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

// CreateNotificationRequest untuk membuat notifikasi (internal use)
type CreateNotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required,max=255"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=order payment system alert"`
	Data    string `json:"data"`
}
