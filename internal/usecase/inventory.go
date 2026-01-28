package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"
)

// inventoryUsecase adalah implementasi dari interface InventoryUsecase
type inventoryUsecase struct {
	inventoryRepo repository.InventoryRepository
}

// NewInventoryUsecase membuat instance baru dari InventoryUsecase
func NewInventoryUsecase(inventoryRepo repository.InventoryRepository) InventoryUsecase {
	return &inventoryUsecase{
		inventoryRepo: inventoryRepo,
	}
}

// CreateInventory membuat item inventory baru
func (u *inventoryUsecase) CreateInventory(ctx context.Context, req dto.InventoryRequest) (*dto.InventoryResponse, error) {
	// Validasi input
	if req.Name == "" {
		return nil, errors.New("nama item tidak boleh kosong")
	}
	if req.Quantity < 0 {
		return nil, errors.New("quantity tidak boleh negatif")
	}
	if req.MinStock < 0 {
		return nil, errors.New("minimum stock tidak boleh negatif")
	}
	if req.Unit == "" {
		return nil, errors.New("satuan tidak boleh kosong")
	}

	// Buat entity inventory
	inventory := &entity.Inventory{
		Name:     req.Name,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		MinStock: req.MinStock,
	}

	// Simpan ke database
	if err := u.inventoryRepo.Create(ctx, inventory); err != nil {
		return nil, fmt.Errorf("gagal membuat inventory: %w", err)
	}

	// Convert ke response
	return u.toInventoryResponse(inventory), nil
}

// UpdateInventory memperbarui item inventory yang sudah ada
func (u *inventoryUsecase) UpdateInventory(ctx context.Context, id int64, req dto.InventoryRequest) (*dto.InventoryResponse, error) {
	// Validasi input
	if req.Name == "" {
		return nil, errors.New("nama item tidak boleh kosong")
	}
	if req.Quantity < 0 {
		return nil, errors.New("quantity tidak boleh negatif")
	}
	if req.MinStock < 0 {
		return nil, errors.New("minimum stock tidak boleh negatif")
	}
	if req.Unit == "" {
		return nil, errors.New("satuan tidak boleh kosong")
	}

	// Cek apakah inventory ada
	inventory, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Update data
	inventory.Name = req.Name
	inventory.Quantity = req.Quantity
	inventory.Unit = req.Unit
	inventory.MinStock = req.MinStock

	// Simpan perubahan
	if err := u.inventoryRepo.Update(ctx, inventory); err != nil {
		return nil, fmt.Errorf("gagal memperbarui inventory: %w", err)
	}

	// Convert ke response
	return u.toInventoryResponse(inventory), nil
}

// DeleteInventory menghapus item inventory berdasarkan ID
func (u *inventoryUsecase) DeleteInventory(ctx context.Context, id int64) error {
	// Cek apakah inventory ada
	_, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Hapus inventory
	if err := u.inventoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("gagal menghapus inventory: %w", err)
	}

	return nil
}

// GetInventoryByID mengambil detail item inventory berdasarkan ID
func (u *inventoryUsecase) GetInventoryByID(ctx context.Context, id int64) (*dto.InventoryResponse, error) {
	// Ambil data dari repository
	inventory, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Convert ke response
	return u.toInventoryResponse(inventory), nil
}

// GetAllInventories mengambil semua item inventory dengan filter dan pagination
func (u *inventoryUsecase) GetAllInventories(ctx context.Context, filter dto.InventoryFilter) (*dto.InventoryListResponse, error) {
	// Set default pagination jika tidak ada
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	// Ambil data dari repository
	inventories, pagination, err := u.inventoryRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data inventory: %w", err)
	}

	// Convert ke response
	var responses []dto.InventoryResponse
	for _, inv := range inventories {
		responses = append(responses, *u.toInventoryResponse(&inv))
	}

	return &dto.InventoryListResponse{
		Data:       responses,
		Pagination: *pagination,
	}, nil
}

// AddStock menambah stok item inventory
func (u *inventoryUsecase) AddStock(ctx context.Context, id int64, amount int) error {
	// Validasi input
	if amount <= 0 {
		return errors.New("jumlah yang ditambahkan harus lebih dari 0")
	}

	// Ambil data inventory
	inventory, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Tambah stok
	newQuantity := inventory.Quantity + amount

	// Update stok
	if err := u.inventoryRepo.UpdateStock(ctx, id, newQuantity); err != nil {
		return fmt.Errorf("gagal menambah stok: %w", err)
	}

	return nil
}

// ReduceStock mengurangi stok item inventory
func (u *inventoryUsecase) ReduceStock(ctx context.Context, id int64, amount int) error {
	// Validasi input
	if amount <= 0 {
		return errors.New("jumlah yang dikurangi harus lebih dari 0")
	}

	// Ambil data inventory
	inventory, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Cek apakah stok cukup
	if inventory.Quantity < amount {
		return fmt.Errorf("stok tidak cukup. tersedia: %d, diminta: %d", inventory.Quantity, amount)
	}

	// Kurangi stok
	newQuantity := inventory.Quantity - amount

	// Update stok
	if err := u.inventoryRepo.UpdateStock(ctx, id, newQuantity); err != nil {
		return fmt.Errorf("gagal mengurangi stok: %w", err)
	}

	return nil
}

// GetLowStockItems mengambil semua item dengan stok rendah
func (u *inventoryUsecase) GetLowStockItems(ctx context.Context) ([]dto.InventoryResponse, error) {
	// Ambil data dari repository
	inventories, err := u.inventoryRepo.GetLowStockItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data low stock: %w", err)
	}

	// Convert ke response
	var responses []dto.InventoryResponse
	for _, inv := range inventories {
		responses = append(responses, *u.toInventoryResponse(&inv))
	}

	return responses, nil
}

// CheckStockAvailability mengecek apakah stok tersedia untuk jumlah tertentu
func (u *inventoryUsecase) CheckStockAvailability(ctx context.Context, id int64, requiredAmount int) (bool, error) {
	// Validasi input
	if requiredAmount <= 0 {
		return false, errors.New("jumlah yang dicek harus lebih dari 0")
	}

	// Ambil data inventory
	inventory, err := u.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("inventory tidak ditemukan: %w", err)
	}

	// Cek ketersediaan stok
	return inventory.Quantity >= requiredAmount, nil
}

// toInventoryResponse mengkonversi entity ke response DTO
func (u *inventoryUsecase) toInventoryResponse(inv *entity.Inventory) *dto.InventoryResponse {
	// Cek apakah stok rendah
	isLowStock := inv.Quantity < inv.MinStock

	return &dto.InventoryResponse{
		ID:         inv.ID,
		Name:       inv.Name,
		Quantity:   inv.Quantity,
		Unit:       inv.Unit,
		MinStock:   inv.MinStock,
		IsLowStock: isLowStock,
		CreatedAt:  inv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  inv.UpdatedAt.Format(time.RFC3339),
	}
}
