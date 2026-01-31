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
	CategoryRepo    CategoryRepository
	ProductRepo     ProductRepository
	DashboardRepo   DashboardRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) Repository {
	return Repository{
		AuthRepo:        NewAuthRepository(db, logger),
		InventoriesRepo: NewInventoriesRepository(db, logger),
		StaffRepo:       NewStaffRepository(db, logger),
		OrderRepo:       NewOrderRepository(db, logger),
		CategoryRepo:    NewCategoryRepository(db, logger),
		ProductRepo:     NewProductRepository(db, logger),
		DashboardRepo:   NewDashboardRepository(db, logger),
	}
}
