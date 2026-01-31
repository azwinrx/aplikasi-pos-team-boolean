package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	StaffRepo       StaffRepository
	InventoriesRepo InventoriesRepository
	OrderRepo       OrderRepository
	RevenueRepo     RevenueRepository
	ReservationRepo ReservationsRepository
}

func NewRepository(db *gorm.DB, logger *zap.Logger) Repository {
	return Repository{
		InventoriesRepo: NewInventoriesRepository(db, logger),
		StaffRepo:       NewStaffRepository(db, logger),
		OrderRepo:       NewOrderRepository(db, logger),
		RevenueRepo:     NewRevenueRepository(db, logger),
		ReservationRepo: NewReservationsRepository(db, logger),
	}
}
