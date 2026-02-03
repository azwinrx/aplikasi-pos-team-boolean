package usecase

import (
	"context"
	"time"

	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"gorm.io/gorm"

	"go.uber.org/zap"
)

type DashboardUseCase interface {
	GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetPopularProducts(ctx context.Context, limit int) ([]dto.PopularProductResponse, error)
	GetNewProducts(ctx context.Context, limit int) ([]dto.NewProductResponse, error)
	ExportDashboard(ctx context.Context) ([]dto.DashboardExportRow, error)
}

// ExportDashboard: export bulanan (bulan, jumlah order, sales, revenue)
func (u *dashboardUseCase) ExportDashboard(ctx context.Context) ([]dto.DashboardExportRow, error) {
	u.logger.Info("Exporting dashboard data (monthly)")

	// Query monthly summary for last 12 months
	type Result struct {
		Month       string
		TotalOrders int
		Sales       int
		Revenue     float64
	}
	var results []Result

	db := getGormDB(u.dashboardRepo)
	if db == nil {
		u.logger.Error("Failed to get gorm.DB from dashboardRepo")
		return nil, nil
	}

	err := db.WithContext(ctx).
		Raw(`
		       SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') as month,
			      COUNT(*) as total_orders,
			      SUM(total_amount) as revenue,
			      SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as sales
		       FROM orders
		       WHERE deleted_at IS NULL
		       GROUP BY month
		       ORDER BY month DESC
		       LIMIT 12
	       `).Scan(&results).Error
	if err != nil {
		u.logger.Error("Failed to export dashboard data", zap.Error(err))
		return nil, err
	}

	var exportRows []dto.DashboardExportRow
	for _, r := range results {
		exportRows = append(exportRows, dto.DashboardExportRow{
			Month:       r.Month,
			TotalOrders: r.TotalOrders,
			Sales:       r.Sales,
			Revenue:     r.Revenue,
		})
	}
	return exportRows, nil
}

type dashboardUseCase struct {
	dashboardRepo repository.DashboardRepository
	logger        *zap.Logger
}

func NewDashboardUseCase(dashboardRepo repository.DashboardRepository, logger *zap.Logger) DashboardUseCase {
	return &dashboardUseCase{
		dashboardRepo: dashboardRepo,
		logger:        logger,
	}
}

func (u *dashboardUseCase) GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	u.logger.Info("Getting dashboard summary")

	// Get daily sales summary
	dailySales, err := u.dashboardRepo.GetDailySalesSummary(ctx)
	if err != nil {
		u.logger.Error("Failed to get daily sales summary", zap.Error(err))
		return nil, err
	}

	// Get monthly sales summary
	monthlySales, err := u.dashboardRepo.GetMonthlySalesSummary(ctx)
	if err != nil {
		u.logger.Error("Failed to get monthly sales summary", zap.Error(err))
		return nil, err
	}

	// Get table summary
	tableSummary, err := u.dashboardRepo.GetTableSummary(ctx)
	if err != nil {
		u.logger.Error("Failed to get table summary", zap.Error(err))
		return nil, err
	}

	response := &dto.DashboardSummaryResponse{
		DailySales:   *dailySales,
		MonthlySales: *monthlySales,
		TableSummary: *tableSummary,
	}

	u.logger.Info("Successfully retrieved dashboard summary")
	return response, nil
}

func (u *dashboardUseCase) GetPopularProducts(ctx context.Context, limit int) ([]dto.PopularProductResponse, error) {
	u.logger.Info("Getting popular products", zap.Int("limit", limit))

	if limit <= 0 {
		limit = 10 // Default limit
	}

	results, err := u.dashboardRepo.GetPopularProducts(ctx, limit)
	if err != nil {
		u.logger.Error("Failed to get popular products", zap.Error(err))
		return nil, err
	}

	// Map to response DTO
	var response []dto.PopularProductResponse
	for _, r := range results {
		response = append(response, dto.PopularProductResponse{
			ID:           r.ProductID,
			ProductImage: r.ProductImage,
			ProductName:  r.ProductName,
			Price:        r.Price,
			TotalSold:    r.TotalSold,
			TotalRevenue: r.TotalRevenue,
			Stock:        r.Stock,
			Availability: r.Availability,
		})
	}

	u.logger.Info("Successfully retrieved popular products", zap.Int("count", len(response)))
	return response, nil
}

func (u *dashboardUseCase) GetNewProducts(ctx context.Context, limit int) ([]dto.NewProductResponse, error) {
	u.logger.Info("Getting new products", zap.Int("limit", limit))

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get products created in last 30 days
	products, err := u.dashboardRepo.GetNewProducts(ctx, 30, limit)
	if err != nil {
		u.logger.Error("Failed to get new products", zap.Error(err))
		return nil, err
	}

	// Map to response DTO
	var response []dto.NewProductResponse
	now := time.Now()
	for _, p := range products {
		daysAgo := int(now.Sub(p.CreatedAt).Hours() / 24)

		response = append(response, dto.NewProductResponse{
			ID:           p.ProductID,
			ProductImage: p.ProductImage,
			ProductName:  p.ProductName,
			Price:        p.Price,
			Stock:        p.Stock,
			Availability: p.Availability,
			TotalSold:    p.TotalSold,
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
			DaysAgo:      daysAgo,
		})
	}

	u.logger.Info("Successfully retrieved new products", zap.Int("count", len(response)))
	return response, nil
}

// getGormDB tries to extract *gorm.DB from dashboardRepo
func getGormDB(repo repository.DashboardRepository) *gorm.DB {
	type withDB interface{ GetDB() *gorm.DB }
	if v, ok := repo.(withDB); ok {
		return v.GetDB()
	}
	return nil
}
