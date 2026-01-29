package entity

import (
	"time"

	"gorm.io/gorm"
)

// Table merepresentasikan tabel tables di database
type Table struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Number    string         `gorm:"type:varchar(10);not null;unique" json:"number"`
	Capacity  int            `gorm:"not null" json:"capacity"`
	Status    string         `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel
func (Table) TableName() string {
	return "tables"
}
