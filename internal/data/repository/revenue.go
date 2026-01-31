package repository

import (
	"context"
	"time"

	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RevenueRepository interface {
	GetRevenueByStatus(ctx context.Context, status string) (*dto.RevenueByStatusResponse, error)
	GetRevenuePerMonth(ctx context.Context, year int, month int) (*dto.RevenuePerMonthResponse, error)
	GetProductRevenueList(ctx context.Context, productID int) (*dto.ProductRevenueListResponse, error)
}

type revenueRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRevenueRepository(db *gorm.DB, logger *zap.Logger) RevenueRepository {
	return &revenueRepository{db, logger}
}

func (r *revenueRepository) GetRevenueByStatus(ctx context.Context, status string) (*dto.RevenueByStatusResponse, error) {
	r.logger.Info("Getting revenue by status", zap.String("status", status))

	var breakdown dto.RevenueStatusBreakdown

	query := `
		SELECT
			status,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as order_count
		FROM orders
		WHERE deleted_at IS NULL AND status = ?
		GROUP BY status
	`

	err := r.db.WithContext(ctx).Raw(query, status).Scan(&breakdown).Error
	if err != nil {
		r.logger.Error("Failed to get revenue by status", zap.Error(err))
		return nil, err
	}

	response := &dto.RevenueByStatusResponse{
		TotalRevenue: breakdown.TotalRevenue,
		Breakdown:    []dto.RevenueStatusBreakdown{breakdown},
	}

	r.logger.Info("Successfully got revenue by status",
		zap.Float64("total_revenue", breakdown.TotalRevenue),
		zap.String("status", status))

	return response, nil
}

func (r *revenueRepository) GetRevenuePerMonth(ctx context.Context, year int, month int) (*dto.RevenuePerMonthResponse, error) {
	r.logger.Info("Getting revenue for specific month", zap.Int("year", year), zap.Int("month", month))

	if year == 0 {
		year = time.Now().Year()
	}

	var detail dto.RevenueMonthlyDetail

	query := `
		SELECT
			EXTRACT(MONTH FROM created_at)::int as month,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as order_count
		FROM orders
		WHERE deleted_at IS NULL
			AND EXTRACT(YEAR FROM created_at) = ?
			AND EXTRACT(MONTH FROM created_at) = ?
		GROUP BY EXTRACT(MONTH FROM created_at)
	`

	err := r.db.WithContext(ctx).Raw(query, year, month).Scan(&detail).Error
	if err != nil {
		r.logger.Error("Failed to get revenue for month", zap.Error(err))
		return nil, err
	}

	monthNames := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if detail.Month > 0 && detail.Month <= 12 {
		detail.MonthName = monthNames[detail.Month]
	}

	response := &dto.RevenuePerMonthResponse{
		Year:         year,
		TotalRevenue: detail.TotalRevenue,
		Monthly:      []dto.RevenueMonthlyDetail{detail},
	}

	r.logger.Info("Successfully got revenue for month",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Float64("total_revenue", detail.TotalRevenue),
	)

	return response, nil
}

func (r *revenueRepository) GetProductRevenueList(ctx context.Context, productID int) (*dto.ProductRevenueListResponse, error) {
	r.logger.Info("Getting product revenue for product", zap.Int("product_id", productID))

	var product dto.ProductRevenueDetail

	query := `
		SELECT
			p.id as product_id,
			p.name as product_name,
			p.price as price,
			COALESCE(SUM(oi.subtotal), 0) as total_revenue,
			COALESCE(SUM(oi.quantity), 0) as total_sold,
			COUNT(DISTINCT oi.order_id) as order_count,
			MAX(o.created_at) as last_order_at
		FROM products p
		LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
		LEFT JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
		WHERE p.deleted_at IS NULL AND p.id = ?
		GROUP BY p.id, p.name, p.price
	`

	err := r.db.WithContext(ctx).Raw(query, productID).Scan(&product).Error
	if err != nil {
		r.logger.Error("Failed to get product revenue", zap.Error(err))
		return nil, err
	}

	response := &dto.ProductRevenueListResponse{
		TotalProducts: 1,
		Products:      []dto.ProductRevenueDetail{product},
	}

	r.logger.Info("Successfully got product revenue",
		zap.Int("product_id", productID))

	return response, nil
}