package entity

import (
	"time"

	"gorm.io/gorm"
)

// PaymentMethod merepresentasikan tabel payment_methods di database
type PaymentMethod struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(50);not null;unique" json:"name"`
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (PaymentMethod) TableName() string {
	return "payment_methods"
}
