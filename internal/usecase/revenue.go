package usecase

import (
	"context"
	"time"

	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
)

type RevenueUseCase interface {
	GetRevenueByStatus(ctx context.Context, status string) (*dto.RevenueByStatusResponse, error)
	GetRevenuePerMonth(ctx context.Context, year int, month int) (*dto.RevenuePerMonthResponse, error)
	GetProductRevenueList(ctx context.Context, productID int) (*dto.ProductRevenueListResponse, error)
}

type revenueUseCase struct {
	revenueRepo repository.RevenueRepository
	logger      *zap.Logger
}

func NewRevenueUseCase(revenueRepo repository.RevenueRepository, logger *zap.Logger) *revenueUseCase {
	return &revenueUseCase{
		revenueRepo: revenueRepo,
		logger:      logger,
	}
}

func (uc *revenueUseCase) GetRevenueByStatus(ctx context.Context, status string) (*dto.RevenueByStatusResponse, error) {
	uc.logger.Info("Getting revenue by status", zap.String("status", status))

	response, err := uc.revenueRepo.GetRevenueByStatus(ctx, status)
	if err != nil {
		uc.logger.Error("Failed to get revenue by status", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Successfully retrieved revenue by status",
		zap.Float64("total_revenue", response.TotalRevenue),
		zap.String("status", status),
		zap.Int("status_count", len(response.Breakdown)))

	return response, nil
}

func (uc *revenueUseCase) GetRevenuePerMonth(ctx context.Context, year int, month int) (*dto.RevenuePerMonthResponse, error) {
	uc.logger.Info("Getting revenue for specific month", zap.Int("year", year), zap.Int("month", month))

	if year == 0 {
		year = time.Now().Year()
		uc.logger.Info("Year not provided, using current year", zap.Int("year", year))
	}

	response, err := uc.revenueRepo.GetRevenuePerMonth(ctx, year, month)
	if err != nil {
		uc.logger.Error("Failed to get revenue for month", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Successfully retrieved revenue for month",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Float64("total_revenue", response.TotalRevenue),
		zap.Int("month_count", len(response.Monthly)))

	return response, nil
}

func (uc *revenueUseCase) GetProductRevenueList(ctx context.Context, productID int) (*dto.ProductRevenueListResponse, error) {
	uc.logger.Info("Getting product revenue for product", zap.Int("product_id", productID))

	response, err := uc.revenueRepo.GetProductRevenueList(ctx, productID)
	if err != nil {
		uc.logger.Error("Failed to get product revenue", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Successfully retrieved product revenue",
		zap.Int("product_id", productID),
		zap.Int("total_products", response.TotalProducts))

	return response, nil
}
