package entity

import (
	"time"

	"gorm.io/gorm"
)

// Order merepresentasikan tabel orders di database
type Order struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	TableID         uint           `gorm:"not null" json:"table_id"`
	PaymentMethodID uint           `gorm:"not null" json:"payment_method_id"`
	CustomerName    string         `gorm:"type:varchar(100);not null" json:"customer_name"`
	TotalAmount     float64        `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	Tax             float64        `gorm:"type:decimal(15,2);not null" json:"tax"`
	Status          string         `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt       time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Table           Table          `gorm:"foreignKey:TableID" json:"table,omitempty"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodID" json:"payment_method,omitempty"`
	Items           []OrderItem    `gorm:"foreignKey:OrderID;references:ID" json:"items"`
}

// OrderItem merepresentasikan tabel order_items di database
type OrderItem struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint           `gorm:"not null" json:"order_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Subtotal  float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName override nama tabel untuk Order
func (Order) TableName() string {
	return "orders"
}

// TableName override nama tabel untuk OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate hook untuk Order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook untuk Order
func (o *Order) BeforeUpdate(tx *gorm.DB) error {
	o.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook untuk OrderItem
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	oi.CreatedAt = time.Now()
	oi.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook untuk OrderItem
func (oi *OrderItem) BeforeUpdate(tx *gorm.DB) error {
	oi.UpdatedAt = time.Now()
	return nil
}
