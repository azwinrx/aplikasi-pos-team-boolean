package repository

import (
	"context"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDailySalesSummary(ctx context.Context) (*dto.SalesSummary, error)
	GetMonthlySalesSummary(ctx context.Context) (*dto.SalesSummary, error)
	GetTableSummary(ctx context.Context) (*dto.TableSummary, error)
	GetPopularProducts(ctx context.Context, limit int) ([]PopularProductResult, error)
	GetNewProducts(ctx context.Context, days int, limit int) ([]NewProductResult, error)
}

// PopularProductResult untuk hasil query produk populer
type PopularProductResult struct {
	ProductID    uint
	ProductImage string
	ProductName  string
	Price        float64
	TotalSold    int
	TotalRevenue float64
	Stock        int
	Availability string
}

type dashboardRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDashboardRepository(db *gorm.DB, logger *zap.Logger) DashboardRepository {
	return &dashboardRepository{db, logger}
}

func (r *dashboardRepository) GetDailySalesSummary(ctx context.Context) (*dto.SalesSummary, error) {
	r.logger.Info("Getting daily sales summary")

	// Get today's date range
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.getSalesSummary(ctx, startOfDay, endOfDay)
}

func (r *dashboardRepository) GetMonthlySalesSummary(ctx context.Context) (*dto.SalesSummary, error) {
	r.logger.Info("Getting monthly sales summary")

	// Get current month's date range
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return r.getSalesSummary(ctx, startOfMonth, endOfMonth)
}

func (r *dashboardRepository) getSalesSummary(ctx context.Context, startDate, endDate time.Time) (*dto.SalesSummary, error) {
	var summary dto.SalesSummary

	// Query for total orders, revenue, tax
	type Result struct {
		TotalOrders  int
		TotalRevenue float64
		TotalTax     float64
	}
	var result Result

	err := r.db.WithContext(ctx).Model(&entity.Order{}).
		Select("COUNT(*) as total_orders, COALESCE(SUM(total_amount), 0) as total_revenue, COALESCE(SUM(tax), 0) as total_tax").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Scan(&result).Error
	if err != nil {
		r.logger.Error("Failed to get sales summary", zap.Error(err))
		return nil, err
	}

	summary.TotalOrders = result.TotalOrders
	summary.TotalRevenue = result.TotalRevenue
	summary.TotalTax = result.TotalTax

	if result.TotalOrders > 0 {
		summary.AverageOrder = result.TotalRevenue / float64(result.TotalOrders)
	}

	// Count orders by status
	type StatusCount struct {
		Status string
		Count  int
	}
	var statusCounts []StatusCount

	err = r.db.WithContext(ctx).Model(&entity.Order{}).
		Select("status, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		r.logger.Error("Failed to get status counts", zap.Error(err))
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "paid":
			summary.PaidOrders = sc.Count
		case "pending":
			summary.PendingOrders = sc.Count
		case "cancelled":
			summary.CancelledOrders = sc.Count
		}
	}

	r.logger.Info("Successfully retrieved sales summary",
		zap.Int("total_orders", summary.TotalOrders),
		zap.Float64("total_revenue", summary.TotalRevenue))

	return &summary, nil
}

func (r *dashboardRepository) GetTableSummary(ctx context.Context) (*dto.TableSummary, error) {
	r.logger.Info("Getting table summary")

	var summary dto.TableSummary

	// Count total tables
	var totalTables int64
	if err := r.db.WithContext(ctx).Model(&entity.Table{}).Count(&totalTables).Error; err != nil {
		r.logger.Error("Failed to count total tables", zap.Error(err))
		return nil, err
	}
	summary.TotalTables = int(totalTables)

	// Count tables by status
	type StatusCount struct {
		Status string
		Count  int
	}
	var statusCounts []StatusCount

	err := r.db.WithContext(ctx).Model(&entity.Table{}).
		Select("status, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		r.logger.Error("Failed to get table status counts", zap.Error(err))
		return nil, err
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "available":
			summary.AvailableTables = sc.Count
		case "occupied":
			summary.OccupiedTables = sc.Count
		case "reserved":
			summary.ReservedTables = sc.Count
		}
	}

	r.logger.Info("Successfully retrieved table summary",
		zap.Int("total_tables", summary.TotalTables),
		zap.Int("available_tables", summary.AvailableTables))

	return &summary, nil
}

func (r *dashboardRepository) GetPopularProducts(ctx context.Context, limit int) ([]PopularProductResult, error) {
	r.logger.Info("Getting popular products", zap.Int("limit", limit))

	var results []PopularProductResult

	// Query popular products based on order_items
	err := r.db.WithContext(ctx).
		Table("order_items oi").
		Select(`
		       p.id as product_id,
		       p.product_image,
		       p.product_name,
		       p.price,
		       SUM(oi.quantity) as total_sold,
		       SUM(oi.subtotal) as total_revenue,
		       p.stock,
		       CASE WHEN p.stock > 0 THEN 'in_stock' ELSE 'out_of_stock' END as availability
	       `).
		Joins("JOIN products p ON oi.product_id = p.id").
		Where("oi.deleted_at IS NULL").
		Where("p.deleted_at IS NULL").
		Group("p.id, p.product_image, p.product_name, p.price, p.stock").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		r.logger.Error("Failed to get popular products", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully retrieved popular products", zap.Int("count", len(results)))
	return results, nil
}

type NewProductResult struct {
	ProductID    uint
	ProductImage string
	ProductName  string
	Price        float64
	Stock        int
	Availability string
	TotalSold    int
	CreatedAt    time.Time
}

func (r *dashboardRepository) GetNewProducts(ctx context.Context, days int, limit int) ([]NewProductResult, error) {
	r.logger.Info("Getting new products", zap.Int("days", days), zap.Int("limit", limit))

	cutoffDate := time.Now().AddDate(0, 0, -days)

	var products []NewProductResult
	err := r.db.WithContext(ctx).
		Table("products p").
		Select(`
		       p.id as product_id,
		       p.product_image,
		       p.product_name,
		       p.price,
		       p.stock,
		       CASE WHEN p.stock > 0 THEN 'in_stock' ELSE 'out_of_stock' END as availability,
		       COALESCE(SUM(oi.quantity), 0) as total_sold,
		       p.created_at
	       `).
		Joins("LEFT JOIN order_items oi ON oi.product_id = p.id AND oi.deleted_at IS NULL").
		Where("p.created_at >= ?", cutoffDate).
		Where("p.deleted_at IS NULL").
		Group("p.id, p.product_image, p.product_name, p.price, p.stock, p.created_at").
		Order("p.created_at DESC").
		Limit(limit).
		Scan(&products).Error

	if err != nil {
		r.logger.Error("Failed to get new products", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully retrieved new products", zap.Int("count", len(products)))
	return products, nil
}

// GetDB returns the *gorm.DB instance
func (r *dashboardRepository) GetDB() *gorm.DB {
	return r.db
}
