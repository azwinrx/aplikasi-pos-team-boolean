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

	AuthUseCase        AuthUseCase
	AdminUseCase       AdminUseCase
	OrderUseCase       OrderUseCase
	InventoriesUsecase InventoriesUsecase
	StaffUseCase       StaffUseCase
	NotificationUseCase NotificationUseCase
	CategoryUseCase    CategoryUseCase
	ProductUseCase     ProductUseCase
	DashboardUseCase   DashboardUseCase
	ReservationsUseCase ReservationsUseCase
	RevenueUseCase      RevenueUseCase
}

func NewUseCase(repo *repository.Repository, logger *zap.Logger, tx *gorm.DB) *UseCase {
	emailService := utils.NewEmailService(logger, utils.Config.SMTP)

	return &UseCase{
		log:  logger,
		repo: *repo,

		AuthUseCase:        NewAuthUseCase(repo.AuthRepo, logger, emailService),
		AdminUseCase:       NewAdminUseCase(repo.AuthRepo, emailService, logger),
		OrderUseCase:       NewOrderUseCase(repo.OrderRepo, logger),
		InventoriesUsecase: NewInventoriesUsecase(repo.InventoriesRepo, logger),
		StaffUseCase:       NewStaffUseCase(repo.StaffRepo, logger),
		NotificationUseCase: NewNotificationUseCase(repo.NotificationRepo, logger),
		CategoryUseCase:    NewCategoryUseCase(repo.CategoryRepo, logger),
		ProductUseCase:     NewProductUseCase(repo.ProductRepo, repo.CategoryRepo, logger),
		DashboardUseCase:   NewDashboardUseCase(repo.DashboardRepo, logger),
		ReservationsUseCase: NewReservationUseCase(repo.ReservationRepo, logger),
		RevenueUseCase:      NewRevenueUseCase(repo.RevenueRepo, logger),
	}
}
