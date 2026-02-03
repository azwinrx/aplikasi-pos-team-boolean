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
// Fields: product name, order count (total_sold), photo, stock status, price
type PopularProductResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	Stock        int     `json:"stock"`
	Availability string  `json:"availability"` // "in_stock" or "out_of_stock"
}

// NewProductResponse untuk daftar produk baru (< 30 hari)
// Fields: product name, order count (total_sold), photo, stock status, price
type NewProductResponse struct {
	ID           uint    `json:"id"`
	ProductImage string  `json:"product_image"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	Availability string  `json:"availability"` // "in_stock" or "out_of_stock"
	TotalSold    int     `json:"total_sold"`
	CreatedAt    string  `json:"created_at"`
	DaysAgo      int     `json:"days_ago"`
}

// DashboardExportRow untuk export bulanan
// Template: bulan, jumlah order, sales, revenue
type DashboardExportRow struct {
	Month       string  `json:"month"`        // Format: YYYY-MM
	TotalOrders int     `json:"total_orders"` // Jumlah order
	Sales       int     `json:"sales"`        // Jumlah order paid/completed
	Revenue     float64 `json:"revenue"`      // Total revenue
}

// DashboardWebsocketMessage untuk realtime websocket data
type DashboardWebsocketMessage struct {
	DailySales    float64 `json:"daily_sales"`    // Today's revenue
	MonthlySales  float64 `json:"monthly_sales"`  // This month's revenue
	DailyOrders   int     `json:"daily_orders"`   // Today's order count
	MonthlyOrders int     `json:"monthly_orders"` // This month's order count
}
