package dto

// InventoriesFilter merepresentasikan parameter filter untuk pencarian Inventories
type InventoriesFilter struct {
	Search   string  `json:"search" form:"search"`       // Pencarian berdasarkan nama
	Status   string  `json:"status" form:"status"`       // Filter berdasarkan status (active/inactive)
	Category string  `json:"category" form:"category"`   // Filter berdasarkan kategori
	Stock    string  `json:"stock" form:"stock"`         // Filter berdasarkan status stok (instock/lowstock/outofstock)
	MinPrice float64 `json:"min_price" form:"min_price"` // Filter harga minimum
	MaxPrice float64 `json:"max_price" form:"max_price"` // Filter harga maksimum
	Unit     string  `json:"unit" form:"unit"`           // Filter berdasarkan satuan (litre/pcs/kg)
	MinQty   int     `json:"min_qty" form:"min_qty"`     // Filter quantity minimum
	MaxQty   int     `json:"max_qty" form:"max_qty"`     // Filter quantity maksimum
	Page     int     `json:"page" form:"page"`           // Nomor halaman untuk pagination
	Limit    int     `json:"limit" form:"limit"`         // Jumlah item per halaman
	SortBy   string  `json:"sort_by" form:"sort_by"`     // Field untuk sorting (name/quantity/price/created_at)
	SortDir  string  `json:"sort_dir" form:"sort_dir"`   // Arah sorting (asc/desc)
}

// InventoriesRequest merepresentasikan request payload untuk create/update Inventories
type InventoriesRequest struct {
	Image       string  `json:"image"`
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=0"`
	Status      string  `json:"status" binding:"required,oneof=active inactive"`
	RetailPrice float64 `json:"retail_price" binding:"required,min=0"`
}

// InventoriesResponse merepresentasikan response payload untuk Inventories
type InventoriesResponse struct {
	ID          int64   `json:"id"`
	Image       string  `json:"image"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Quantity    int     `json:"quantity"`
	Status      string  `json:"status"`
	RetailPrice float64 `json:"retail_price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// InventoriesListResponse merepresentasikan response list dengan pagination
type InventoriesListResponse struct {
	Data       []InventoriesResponse `json:"data"`
	Pagination Pagination            `json:"pagination"`
}

// InventoriesSummary merepresentasikan summary/statistik inventory
type InventoriesSummary struct {
	TotalProducts      int `json:"total_products"`
	ActiveProducts     int `json:"active_products"`
	InactiveProducts   int `json:"inactive_products"`
	LowStockProducts   int `json:"low_stock_products"`
	OutOfStockProducts int `json:"out_of_stock_products"`
}
