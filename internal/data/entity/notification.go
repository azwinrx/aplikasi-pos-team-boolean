package entity

import (
	"time"

	"gorm.io/gorm"
)

// Notification merepresentasikan tabel notifications di database
type Notification struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Type      string         `gorm:"type:varchar(50);not null;index" json:"type"` // order, payment, system, alert
	Status    string         `gorm:"type:varchar(20);not null;default:'new';index" json:"status"` // new, readed
	ReadedAt  *time.Time     `gorm:"type:timestamp;nullable" json:"readed_at,omitempty"`
	Data      string         `gorm:"type:jsonb;nullable" json:"data,omitempty"` // Extra data in JSON format
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook untuk set default values
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.Status == "" {
		n.Status = "new"
	}
	return nil
}
