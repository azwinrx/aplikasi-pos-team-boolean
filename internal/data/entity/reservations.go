package entity

import (
	"time"

	"gorm.io/gorm"
)

// Reservations merepresentasikan tabel reservations di database
type Reservations struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerName    string         `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerPhone   string         `gorm:"type:varchar(255);not null" json:"customer_phone"`
	TableID         int64          `gorm:"not null" json:"table_id"`
	ReservationTime *time.Time     `gorm:"type:timestamp with time zone" json:"reservation_time,omitempty"`
	Status          string         `gorm:"type:varchar(255);not null;default:'pending'" json:"status"`
	CreatedAt       time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"type:timestamp with time zone;index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel untuk Reservations
func (Reservations) TableName() string {
	return "reservations"
}

// BeforeCreate hook untuk Reservations
func (r *Reservations) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

// BeforeUpdate hook untuk Reservations
func (r *Reservations) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}
