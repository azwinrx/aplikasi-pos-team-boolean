package adaptor

import (
	"aplikasi-pos-team-boolean/internal/usecase"

	"go.uber.org/zap"
)

// Adaptor is a struct that holds all HTTP handlers for the application
type Adaptor struct {
	AuthAdaptor        *AuthAdaptor
	InventoriesAdaptor *InventoriesAdaptor
	StaffAdaptor       *StaffAdaptor
	OrderAdaptor       *OrderAdaptor
	CategoryAdaptor    *CategoryAdaptor
	ProductAdaptor     *ProductAdaptor
	DashboardAdaptor   DashboardHandler
	RevenueAdaptor      *RevenueAdaptor
	ReservationsAdaptor *ReservationsAdaptor
}

// NewAdaptor creates a new instance of Adaptor with all handlers
func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger) *Adaptor {
	return &Adaptor{
		AuthAdaptor:        NewAuthAdaptor(uc.AuthUseCase, logger),
		InventoriesAdaptor: NewInventoriesAdaptor(uc.InventoriesUsecase, logger),
		StaffAdaptor:       NewStaffAdaptor(uc.StaffUseCase, logger),
		OrderAdaptor:       NewOrderAdaptor(uc.OrderUseCase, logger),
		CategoryAdaptor:    NewCategoryAdaptor(uc.CategoryUseCase, logger),
		ProductAdaptor:     NewProductAdaptor(uc.ProductUseCase, logger),
		DashboardAdaptor:   NewDashboardHandler(uc.DashboardUseCase, logger),
		RevenueAdaptor:      NewRevenueAdaptor(uc.RevenueUseCase, logger),
		ReservationsAdaptor: NewReservationsAdaptor(uc.ReservationsUseCase, logger),
	}
}
