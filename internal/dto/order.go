package dto

import "time"

// OrderItemRequest untuk item dalam order
type OrderItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

// OrderCreateRequest untuk create order baru
type OrderCreateRequest struct {
	UserID          uint               `json:"user_id" binding:"required"`
	TableID         uint               `json:"table_id" binding:"required"`
	PaymentMethodID uint               `json:"payment_method_id" binding:"required"`
	CustomerName    string             `json:"customer_name" binding:"required,min=1,max=100"`
	Items           []OrderItemRequest `json:"items" binding:"required,min=1"`
	Tax             float64            `json:"tax" binding:"min=0"` // Pajak statis, bisa dihitung otomatis jika tidak disediakan
}

// OrderUpdateRequest untuk update order
type OrderUpdateRequest struct {
	CustomerName    string             `json:"customer_name" binding:"required,min=1,max=100"`
	PaymentMethodID uint               `json:"payment_method_id" binding:"required"`
	Items           []OrderItemRequest `json:"items" binding:"required,min=1"`
}

// OrderResponse untuk response order
type OrderResponse struct {
	ID              uint                `json:"id"`
	UserID          uint                `json:"user_id"`
	TableID         uint                `json:"table_id"`
	PaymentMethodID uint                `json:"payment_method_id"`
	CustomerName    string              `json:"customer_name"`
	TotalAmount     float64             `json:"total_amount"`
	Tax             float64             `json:"tax"`
	Status          string              `json:"status"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Items           []OrderItemResponse `json:"items"`
}

// OrderItemResponse untuk response item order
type OrderItemResponse struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// OrderListResponse untuk response list order
type OrderListResponse struct {
	ID           uint      `json:"id"`
	CustomerName string    `json:"customer_name"`
	TableNumber  string    `json:"table_number"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// OrderDetailResponse untuk response detail order
type OrderDetailResponse struct {
	ID              uint                `json:"id"`
	UserID          uint                `json:"user_id"`
	TableID         uint                `json:"table_id"`
	TableNumber     string              `json:"table_number"`
	PaymentMethodID uint                `json:"payment_method_id"`
	PaymentMethod   string              `json:"payment_method"`
	CustomerName    string              `json:"customer_name"`
	TotalAmount     float64             `json:"total_amount"`
	Tax             float64             `json:"tax"`
	Status          string              `json:"status"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Items           []OrderItemResponse `json:"items"`
}

// TableResponse untuk response table
type TableResponse struct {
	ID       uint   `json:"id"`
	Number   string `json:"number"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}

// PaymentMethodResponse untuk response payment method
type PaymentMethodResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
