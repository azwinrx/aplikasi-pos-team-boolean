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
	RevenueUseCase     *revenueUseCase
	ReservationUseCase ReservationsUseCase
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	return &UseCase{
		log:  logger,
		repo: *repo,

		OrderUseCase:       NewOrderUseCase(repo.OrderRepo, logger),
		InventoriesUsecase: NewInventoriesUsecase(repo.InventoriesRepo, logger),
		StaffUseCase:       NewStaffUseCase(repo.StaffRepo, logger),
		RevenueUseCase:     NewRevenueUseCase(repo.RevenueRepo, logger),
		ReservationUseCase: NewReservationUseCase(repo.ReservationRepo, logger),
	}
}
