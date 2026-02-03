package dto

import "time"

// RevenueByStatusResponse untuk response total revenue dan breakdown berdasarkan status
type RevenueByStatusResponse struct {
	TotalRevenue float64                  `json:"total_revenue"`
	Breakdown    []RevenueStatusBreakdown `json:"breakdown"`
}

// RevenueStatusBreakdown untuk breakdown revenue per status
type RevenueStatusBreakdown struct {
	Status       string  `json:"status"`
	TotalRevenue float64 `json:"total_revenue"`
	OrderCount   int     `json:"order_count"`
}

// RevenuePerMonthResponse untuk response total revenue per bulan
type RevenuePerMonthResponse struct {
	Year         int                    `json:"year"`
	TotalRevenue float64                `json:"total_revenue"`
	Monthly      []RevenueMonthlyDetail `json:"monthly"`
}

// RevenueMonthlyDetail untuk detail revenue per bulan
type RevenueMonthlyDetail struct {
	Month        int     `json:"month"`
	MonthName    string  `json:"month_name"`
	TotalRevenue float64 `json:"total_revenue"`
	OrderCount   int     `json:"order_count"`
}

// ProductRevenueListResponse untuk response list produk dengan detail revenue
type ProductRevenueListResponse struct {
	TotalProducts int                    `json:"total_products"`
	Products      []ProductRevenueDetail `json:"products"`
}

// ProductRevenueDetail untuk detail revenue per produk
type ProductRevenueDetail struct {
	ProductID    uint      `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Price        float64   `json:"price"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalSold    int       `json:"total_sold"`
	OrderCount   int       `json:"order_count"`
	LastOrderAt  time.Time `json:"last_order_at"`
}
