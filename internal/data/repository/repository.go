package repository

import (
	"gorm.io/gorm"
	"go.uber.org/zap"
)

type Repository struct {
	StaffRepo       StaffRepository
	InventoriesRepo InventoriesRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) Repository {
	return Repository{
		InventoriesRepo: NewInventoriesRepository(db, logger),
		StaffRepo:       NewStaffRepository(db, logger),
	}
}
