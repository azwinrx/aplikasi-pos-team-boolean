package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	AuthRepo        AuthRepository
	StaffRepo       StaffRepository
	InventoriesRepo InventoriesRepository
	OrderRepo       OrderRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) Repository {
	return Repository{
		AuthRepo:        NewAuthRepository(db, logger),
		InventoriesRepo: NewInventoriesRepository(db, logger),
		StaffRepo:       NewStaffRepository(db, logger),
		OrderRepo:       NewOrderRepository(db, logger),
	}
}
