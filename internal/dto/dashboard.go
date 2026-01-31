package dto

// DashboardSummaryResponse untuk response ringkasan dashboard
type DashboardSummaryResponse struct {
	DailySales   SalesSummary `json:"daily_sales"`
	MonthlySales SalesSummary `json:"monthly_sales"`
	TableSummary TableSummary `json:"table_summary"`
}

// SalesSummary untuk ringkasan penjualan
type SalesSummary struct {
	TotalOrders     int     `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalTax        float64 `json:"total_tax"`
	AverageOrder    float64 `json:"average_order"`
	PaidOrders      int     `json:"paid_orders"`
	PendingOrders   int     `json:"pending_orders"`
	CancelledOrders int     `json:"cancelled_orders"`
}

// TableSummary untuk ringkasan meja
type TableSummary struct {
	TotalTables     int `json:"total_tables"`
	AvailableTables int `json:"available_tables"`
	OccupiedTables  int `json:"occupied_tables"`
	ReservedTables  int `json:"reserved_tables"`
}

// PopularProductResponse untuk daftar produk populer
type PopularProductResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	ItemID       string  `json:"item_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

// NewProductResponse untuk daftar produk baru (< 30 hari)
type NewProductResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	ItemID       string  `json:"item_id"`
	CategoryName string  `json:"category_name"`
	Stock        int     `json:"stock"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
	Availability string  `json:"availability"`
	CreatedAt    string  `json:"created_at"`
	DaysAgo      int     `json:"days_ago"`
}
