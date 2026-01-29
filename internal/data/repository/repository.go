package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	StaffRepo       StaffRepository
	InventoriesRepo InventoriesRepository
	OrderRepo       OrderRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) Repository {
	return Repository{
		InventoriesRepo: NewInventoriesRepository(db, logger),
		StaffRepo:       NewStaffRepository(db, logger),
		OrderRepo:       NewOrderRepository(db, logger),
	}
}
