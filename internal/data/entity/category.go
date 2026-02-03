package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category merepresentasikan tabel category di database
type Category struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	IconCategory string         `gorm:"type:varchar(255)" json:"icon_category"`
	CategoryName string         `gorm:"type:varchar(100);not null;unique;index" json:"category_name"`
	Description  string         `gorm:"type:text" json:"description"`
	CreatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relation
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName override nama tabel
func (Category) TableName() string {
	return "categories"
}

// BeforeUpdate hook untuk update timestamp
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
