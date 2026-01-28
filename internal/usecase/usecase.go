package usecase

import (
	"context"

	"aplikasi-pos-team-boolean/internal/dto"
)

// InventoryUsecase mendefinisikan interface untuk business logic inventory
type InventoryUsecase interface {
	// CreateInventory membuat item inventory baru
	CreateInventory(ctx context.Context, req dto.InventoryRequest) (*dto.InventoryResponse, error)

	// UpdateInventory memperbarui item inventory yang sudah ada
	UpdateInventory(ctx context.Context, id int64, req dto.InventoryRequest) (*dto.InventoryResponse, error)

	// DeleteInventory menghapus item inventory berdasarkan ID
	DeleteInventory(ctx context.Context, id int64) error

	// GetInventoryByID mengambil detail item inventory berdasarkan ID
	GetInventoryByID(ctx context.Context, id int64) (*dto.InventoryResponse, error)

	// GetAllInventories mengambil semua item inventory dengan filter dan pagination
	GetAllInventories(ctx context.Context, filter dto.InventoryFilter) (*dto.InventoryListResponse, error)

	// AddStock menambah stok item inventory
	AddStock(ctx context.Context, id int64, amount int) error

	// ReduceStock mengurangi stok item inventory
	ReduceStock(ctx context.Context, id int64, amount int) error

	// GetLowStockItems mengambil semua item dengan stok rendah
	GetLowStockItems(ctx context.Context) ([]dto.InventoryResponse, error)

	// CheckStockAvailability mengecek apakah stok tersedia untuk jumlah tertentu
	CheckStockAvailability(ctx context.Context, id int64, requiredAmount int) (bool, error)
}
