package dto

// ProductFilterRequest untuk filter product list
type ProductFilterRequest struct {
	Page        int     `json:"page" form:"page"`
	Limit       int     `json:"limit" form:"limit"`
	Search      string  `json:"search" form:"search"`             // Cari berdasarkan product_name, item_id
	CategoryID  uint    `json:"category_id" form:"category_id"`   // Filter by category
	IsAvailable *bool   `json:"is_available" form:"is_available"` // Filter by availability: true/false
	SortOrder   string  `json:"sort_order" form:"sort_order"`     // asc, desc
	MinPrice    float64 `json:"min_price" form:"min_price"`
	MaxPrice    float64 `json:"max_price" form:"max_price"`
}

// ProductCreateRequest untuk create product baru
type ProductCreateRequest struct {
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name" binding:"required,min=2,max=100"`
	Stock        int     `json:"stock" binding:"min=0"`
	CategoryID   uint    `json:"category_id" binding:"required"`
	Price        float64 `json:"price" binding:"required,min=0"`
}

// ProductUpdateRequest untuk update product
type ProductUpdateRequest struct {
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name" binding:"required,min=2,max=100"`
	Stock        int     `json:"stock" binding:"min=0"`
	CategoryID   uint    `json:"category_id" binding:"required"`
	Price        float64 `json:"price" binding:"required,min=0"`
}

// ProductResponse untuk response product detail
type ProductResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	ItemID       string  `json:"item_id"`
	Stock        int     `json:"stock"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
	Availability string  `json:"availability"` // "in_stock" atau "out_of_stock"
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// ProductListResponse untuk response list product (sesuai UI)
type ProductListResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	ItemID       string  `json:"item_id"`
	Stock        int     `json:"stock"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
	Availability string  `json:"availability"` // "in_stock" atau "out_of_stock"
}
