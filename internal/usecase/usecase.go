package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UseCase struct {
	log  *zap.Logger
	repo repository.Repository

	OrderUseCase       *orderUseCase
	InventoriesUsecase *inventoriesUsecase
	StaffUseCase       *staffUseCase
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{
		log:  logger,
		repo: *repo,

		OrderUseCase:       NewOrderUseCase(repo.OrderRepo, logger),
		InventoriesUsecase: NewInventoriesUsecase(repo.InventoriesRepo, logger),
		StaffUseCase:       NewStaffUseCase(repo.StaffRepo, logger),
	}
}
