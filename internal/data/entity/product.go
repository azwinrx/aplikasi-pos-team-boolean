package entity

import (
	"time"

	"gorm.io/gorm"
)

// Product merepresentasikan tabel product di database
type Product struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductImage string         `gorm:"type:varchar(255)" json:"product_image"`
	ProductName  string         `gorm:"type:varchar(100);not null;index" json:"product_name"`
	ItemID       string         `gorm:"type:varchar(50);unique;not null;index" json:"item_id"`
	Stock        int            `gorm:"type:int;not null;default:0" json:"stock"`
	CategoryID   uint           `gorm:"not null;index" json:"category_id"`
	Price        float64        `gorm:"type:decimal(15,2);not null;default:0" json:"price"`
	IsAvailable  bool           `gorm:"type:boolean;not null;default:false" json:"is_available"`
	CreatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relation
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName override nama tabel
func (Product) TableName() string {
	return "products"
}

// BeforeCreate hook untuk set is_available berdasarkan stock
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.updateAvailability()
	return nil
}

// BeforeUpdate hook untuk update timestamp dan is_available
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	p.updateAvailability()
	return nil
}

// updateAvailability mengupdate is_available berdasarkan stock
func (p *Product) updateAvailability() {
	p.IsAvailable = p.Stock > 0
}

// GetAvailabilityStatus mengembalikan status availability sebagai string untuk response
func (p *Product) GetAvailabilityStatus() string {
	if p.Stock > 0 {
		return "in_stock"
	}
	return "out_of_stock"
}
