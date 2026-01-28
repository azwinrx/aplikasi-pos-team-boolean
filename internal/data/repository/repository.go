package repository

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"
)

// InventoryRepository mendefinisikan interface untuk operasi data inventory
type InventoryRepository interface {
	// Create menambahkan item inventory baru ke database
	Create(ctx context.Context, inventory *entity.Inventory) error

	// Update memperbarui item inventory yang sudah ada
	Update(ctx context.Context, inventory *entity.Inventory) error

	// Delete menghapus item inventory berdasarkan ID (soft delete)
	Delete(ctx context.Context, id int64) error

	// FindByID mengambil satu item inventory berdasarkan ID
	FindByID(ctx context.Context, id int64) (*entity.Inventory, error)

	// FindAll mengambil semua item inventory dengan filter dan pagination
	FindAll(ctx context.Context, filter dto.InventoryFilter) ([]entity.Inventory, *dto.Pagination, error)

	// Count menghitung total item inventory yang sesuai dengan filter
	Count(ctx context.Context, filter dto.InventoryFilter) (int64, error)

	// UpdateStock memperbarui hanya quantity dari item inventory
	UpdateStock(ctx context.Context, id int64, quantity int) error

	// GetLowStockItems mengambil semua item dengan quantity di bawah min_stock
	GetLowStockItems(ctx context.Context) ([]entity.Inventory, error)
}
