package adaptor

import (
	"aplikasi-pos-team-boolean/internal/usecase"

	"go.uber.org/zap"
)

// Adaptor is a struct that holds all HTTP handlers for the application
type Adaptor struct {
	AuthAdaptor        *AuthAdaptor
	AdminAdaptor       *AdminAdaptor
	InventoriesAdaptor *InventoriesAdaptor
	StaffAdaptor       *StaffAdaptor
	OrderAdaptor       *OrderAdaptor
	NotificationAdaptor *NotificationAdaptor
}

// NewAdaptor creates a new instance of Adaptor with all handlers
func NewAdaptor(uc *usecase.UseCase, logger *zap.Logger) *Adaptor {
	return &Adaptor{
		AuthAdaptor:        NewAuthAdaptor(uc.AuthUseCase, logger),
		AdminAdaptor:       NewAdminAdaptor(uc.AdminUseCase, logger),
		InventoriesAdaptor: NewInventoriesAdaptor(uc.InventoriesUsecase, logger),
		StaffAdaptor:       NewStaffAdaptor(uc.StaffUseCase, logger),
		OrderAdaptor:       NewOrderAdaptor(uc.OrderUseCase, logger),
		NotificationAdaptor: NewNotificationAdaptor(uc.NotificationUseCase, logger),
	}
}
