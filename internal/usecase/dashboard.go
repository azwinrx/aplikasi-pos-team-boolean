package usecase

import (
	"context"
	"time"

	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
)

type DashboardUseCase interface {
	GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error)
	GetPopularProducts(ctx context.Context, limit int) ([]dto.PopularProductResponse, error)
	GetNewProducts(ctx context.Context, limit int) ([]dto.NewProductResponse, error)
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
			ItemID:       r.ItemID,
			CategoryName: r.CategoryName,
			Price:        r.Price,
			TotalSold:    r.TotalSold,
			TotalRevenue: r.TotalRevenue,
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
			ID:           p.ID,
			ProductImage: p.ProductImage,
			ProductName:  p.ProductName,
			ItemID:       p.ItemID,
			CategoryName: p.Category.CategoryName,
			Price:        p.Price,
			Stock:        p.Stock,
			IsAvailable:  p.IsAvailable,
			Availability: p.GetAvailabilityStatus(),
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
			DaysAgo:      daysAgo,
		})
	}

	u.logger.Info("Successfully retrieved new products", zap.Int("count", len(response)))
	return response, nil
}
