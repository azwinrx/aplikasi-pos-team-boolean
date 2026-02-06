package dto

// CategoryFilterRequest untuk filter category list
type CategoryFilterRequest struct {
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
	Search    string `json:"search" form:"search"`         // Cari berdasarkan category_name, description
	SortBy    string `json:"sort_by" form:"sort_by"`       // category_name, created_at
	SortOrder string `json:"sort_order" form:"sort_order"` // asc, desc
}

// CategoryCreateRequest untuk create category baru
type CategoryCreateRequest struct {
	IconCategory string `json:"icon_category"`
	CategoryName string `json:"category_name" binding:"required,min=2,max=100"`
	Description  string `json:"description"`
}

// CategoryUpdateRequest untuk update category
type CategoryUpdateRequest struct {
	IconCategory string `json:"icon_category"`
	CategoryName string `json:"category_name" binding:"required,min=2,max=100"`
	Description  string `json:"description"`
}

// CategoryResponse untuk response category
type CategoryResponse struct {
	ID           uint   `json:"id"`
	IconCategory string `json:"icon_category"`
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
	ProductCount int    `json:"product_count,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CategoryListResponse untuk response list category yang sederhana
type CategoryListResponse struct {
	ID           uint   `json:"id"`
	IconCategory string `json:"icon_category"`
	CategoryName string `json:"category_name"`
	ProductCount int    `json:"product_count"`
}
