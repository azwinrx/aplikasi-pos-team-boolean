package entity

import (
	"time"

	"gorm.io/gorm"
)

// Inventories merepresentasikan tabel inventories di database
type Inventories struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Image       string         `gorm:"type:varchar(500)" json:"image"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Category    string         `gorm:"type:varchar(100);default:'uncategorized'" json:"category"`
	Quantity    int            `gorm:"type:integer;default:0;not null" json:"quantity"`
	Status      string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	RetailPrice float64        `gorm:"type:decimal(10,2);not null;default:0" json:"retail_price"`
	Unit        string         `gorm:"type:varchar(50);not null" json:"unit"`
	MinStock    int            `gorm:"type:integer;default:5;not null" json:"min_stock"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;index" json:"deleted_at"`
}

// TableName override nama tabel
func (Inventories) TableName() string {
	return "inventories"
}
