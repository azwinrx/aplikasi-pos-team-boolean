package entity

import (
	"time"

	"gorm.io/gorm"
)

// Inventory merepresentasikan tabel inventories di database
type Inventory struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar;not null" json:"name"`
	Quantity  int            `gorm:"type:integer;default:0;not null" json:"quantity"`
	Unit      string         `gorm:"type:varchar;not null" json:"unit"`
	MinStock  int            `gorm:"type:integer;default:5;not null" json:"min_stock"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp with time zone;index" json:"deleted_at,omitempty"`
}

// TableName mengoverride nama tabel yang digunakan oleh Inventory menjadi `inventories`
func (Inventory) TableName() string {
	return "inventories"
}
