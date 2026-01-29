package adaptor

import (
	"aplikasi-pos-team-boolean/internal/usecase"

	"go.uber.org/zap"
)

// Adaptor is a struct that holds all HTTP handlers for the application
type Adaptor struct {
	InventoriesAdaptor *InventoriesAdaptor
	StaffAdaptor       *StaffAdaptor
	OrderAdaptor       *OrderAdaptor
}

// NewAdaptor creates a new instance of Adaptor with all handlers
func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger) *Adaptor {
	return &Adaptor{
		InventoriesAdaptor: NewInventoriesAdaptor(uc.InventoriesUsecase, logger),
		StaffAdaptor:       NewStaffAdaptor(uc.StaffUseCase, logger),
		OrderAdaptor:       NewOrderAdaptor(uc.OrderUseCase, logger),
	}
}
