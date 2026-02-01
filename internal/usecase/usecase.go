package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

package usecase

import (
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/pkg/utils"

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
	AuthUseCase        AuthUseCas
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	emailService := utils.NewEmailService(logger)

	return &UseCase{
		log:  logger,
		repo: *repo,

		AuthUseCase:        NewAuthUseCase(repo.AuthRepo, logger, emailService),
		OrderUseCase:       NewOrderUseCase(repo.OrderRepo, logger),
		InventoriesUsecase: NewInventoriesUsecase(repo.InventoriesRepo, logger),
		StaffUseCase:       NewStaffUseCase(repo.StaffRepo, logger),
		RevenueUseCase:     NewRevenueUseCase(repo.RevenueRepo, logger),
		ReservationUseCase: NewReservationUseCase(repo.ReservationRepo, logger),
	}
}
