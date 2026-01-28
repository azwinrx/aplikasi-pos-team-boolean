package dto

// InventoryFilter merepresentasikan parameter filter untuk pencarian inventory
type InventoryFilter struct {
	Search   string `json:"search" form:"search"`       // Pencarian berdasarkan nama
	Status   string `json:"status" form:"status"`       // Filter berdasarkan status (active/inactive)
	Category string `json:"category" form:"category"`   // Filter berdasarkan kategori
	Stock    string `json:"stock" form:"stock"`         // Filter berdasarkan status stok (instock/lowstock/outofstock)
	MinPrice int    `json:"min_price" form:"min_price"` // Filter harga minimum
	MaxPrice int    `json:"max_price" form:"max_price"` // Filter harga maksimum
	Unit     string `json:"unit" form:"unit"`           // Filter berdasarkan satuan (litre/pcs/kg)
	Page     int    `json:"page" form:"page"`           // Nomor halaman untuk pagination
	Limit    int    `json:"limit" form:"limit"`         // Jumlah item per halaman
	SortBy   string `json:"sort_by" form:"sort_by"`     // Field untuk sorting (name/quantity/price/created_at)
	SortDir  string `json:"sort_dir" form:"sort_dir"`   // Arah sorting (asc/desc)
}

// InventoryRequest merepresentasikan request payload untuk create/update inventory
type InventoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=0"`
	Unit     string `json:"unit" binding:"required"`
	MinStock int    `json:"min_stock" binding:"required,min=0"`
}

// InventoryResponse merepresentasikan response payload untuk inventory
type InventoryResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Unit       string `json:"unit"`
	MinStock   int    `json:"min_stock"`
	IsLowStock bool   `json:"is_low_stock"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// InventoryListResponse merepresentasikan response list dengan pagination
type InventoryListResponse struct {
	Data       []InventoryResponse `json:"data"`
	Pagination Pagination          `json:"pagination"`
}
