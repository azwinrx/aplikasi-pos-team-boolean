package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
)

// InventoriesUsecase mendefinisikan interface untuk business logic inventories
type InventoriesUsecase interface {
	CreateInventory(ctx context.Context, req dto.InventoriesRequest) (*dto.InventoriesResponse, error)
	UpdateInventory(ctx context.Context, id int64, req dto.InventoriesRequest) (*dto.InventoriesResponse, error)
	DeleteInventory(ctx context.Context, id int64) error
	GetInventoryByFilter(ctx context.Context, filter dto.InventoriesFilter) (*dto.InventoriesListResponse, error)
	GetAllInventories(ctx context.Context, filter dto.InventoriesFilter) (*dto.InventoriesListResponse, error)
}

// inventoriesUsecase adalah implementasi dari interface InventoriesUsecase
type inventoriesUsecase struct {
	inventoriesRepo repository.InventoriesRepository
	logger          *zap.Logger
}

// NewInventoriesUsecase membuat instance baru dari InventoriesUsecase
func NewInventoriesUsecase(inventoriesRepo repository.InventoriesRepository, logger *zap.Logger) InventoriesUsecase {
	return &inventoriesUsecase{
		inventoriesRepo: inventoriesRepo,
		logger:          logger,
	}
}

// CreateInventory membuat item inventory baru
func (u *inventoriesUsecase) CreateInventory(ctx context.Context, req dto.InventoriesRequest) (*dto.InventoriesResponse, error) {
	u.logger.Info("Creating new inventory item",
		zap.String("name", req.Name),
		zap.String("category", req.Category),
		zap.Int("quantity", req.Quantity),
		zap.String("status", req.Status),
		zap.Float64("retail_price", req.RetailPrice),
	)

	// Validasi input
	if req.Name == "" {
		u.logger.Warn("Validation failed: nama item tidak boleh kosong")
		return nil, errors.New("nama item tidak boleh kosong")
	}
	if req.Category == "" {
		u.logger.Warn("Validation failed: kategori tidak boleh kosong")
		return nil, errors.New("kategori tidak boleh kosong")
	}
	if req.Quantity < 0 {
		u.logger.Warn("Validation failed: quantity tidak boleh negatif", zap.Int("quantity", req.Quantity))
		return nil, errors.New("quantity tidak boleh negatif")
	}
	if req.RetailPrice < 0 {
		u.logger.Warn("Validation failed: harga tidak boleh negatif", zap.Float64("retail_price", req.RetailPrice))
		return nil, errors.New("harga tidak boleh negatif")
	}
	if req.Status != "active" && req.Status != "inactive" {
		u.logger.Warn("Validation failed: status invalid", zap.String("status", req.Status))
		return nil, errors.New("status harus active atau inactive")
	}

	// Buat entity inventory
	inventory := &entity.Inventories{
		Image:       req.Image,
		Name:        req.Name,
		Category:    req.Category,
		Quantity:    req.Quantity,
		Status:      req.Status,
		RetailPrice: req.RetailPrice,
	}

	// Simpan ke database
	if err := u.inventoriesRepo.Create(ctx, inventory); err != nil {
		u.logger.Error("Failed to create inventory in database",
			zap.Error(err),
			zap.String("name", req.Name),
		)
		return nil, fmt.Errorf("gagal membuat inventory: %w", err)
	}

	u.logger.Info("Inventory item created successfully",
		zap.Int64("id", inventory.ID),
		zap.String("name", inventory.Name),
	)

	// Convert ke response
	return u.toInventoryResponse(inventory), nil
}

// GetInventoryByFilter mengambil item inventory berdasarkan parameter filter query
func (u *inventoriesUsecase) GetInventoryByFilter(ctx context.Context, filter dto.InventoriesFilter) (*dto.InventoriesListResponse, error) {
	u.logger.Debug("Fetching inventory with filters",
		zap.String("search", filter.Search),
		zap.String("status", filter.Status),
		zap.String("category", filter.Category),
		zap.String("stock", filter.Stock),
		zap.Int("page", filter.Page),
		zap.Int("limit", filter.Limit),
	)

	// Set default pagination jika tidak ada
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	// Ambil data dari repository
	inventories, totalItems, err := u.inventoriesRepo.FindByFilter(ctx, filter)
	if err != nil {
		u.logger.Error("Failed to fetch inventory from database",
			zap.Error(err),
			zap.String("search", filter.Search),
		)
		return nil, fmt.Errorf("gagal mengambil data inventory: %w", err)
	}

	u.logger.Info("Inventory fetched successfully",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(inventories)),
		zap.Int("page", filter.Page),
	)

	// Convert ke response
	var responses []dto.InventoriesResponse
	for _, inv := range inventories {
		responses = append(responses, *u.toInventoryResponse(&inv))
	}

	// Calculate total pages
	totalPages := int(totalItems) / filter.Limit
	if int(totalItems)%filter.Limit != 0 {
		totalPages++
	}

	pagination := dto.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
		TotalItems: int(totalItems),
	}

	return &dto.InventoriesListResponse{
		Data:       responses,
		Pagination: pagination,
	}, nil
}

// UpdateInventory memperbarui item inventory yang sudah ada
func (u *inventoriesUsecase) UpdateInventory(ctx context.Context, id int64, req dto.InventoriesRequest) (*dto.InventoriesResponse, error) {
	u.logger.Info("Updating inventory item",
		zap.Int64("id", id),
		zap.String("name", req.Name),
		zap.String("category", req.Category),
	)

	// Validasi input
	if req.Name == "" {
		u.logger.Warn("Validation failed: nama item tidak boleh kosong")
		return nil, errors.New("nama item tidak boleh kosong")
	}
	if req.Category == "" {
		u.logger.Warn("Validation failed: kategori tidak boleh kosong")
		return nil, errors.New("kategori tidak boleh kosong")
	}
	if req.Quantity < 0 {
		u.logger.Warn("Validation failed: quantity tidak boleh negatif", zap.Int("quantity", req.Quantity))
		return nil, errors.New("quantity tidak boleh negatif")
	}
	if req.RetailPrice < 0 {
		u.logger.Warn("Validation failed: harga tidak boleh negatif", zap.Float64("retail_price", req.RetailPrice))
		return nil, errors.New("harga tidak boleh negatif")
	}
	if req.Status != "active" && req.Status != "inactive" {
		u.logger.Warn("Validation failed: status invalid", zap.String("status", req.Status))
		return nil, errors.New("status harus active atau inactive")
	}

	// Update entity inventory
	inventory := &entity.Inventories{
		ID:          id,
		Image:       req.Image,
		Name:        req.Name,
		Category:    req.Category,
		Quantity:    req.Quantity,
		Status:      req.Status,
		RetailPrice: req.RetailPrice,
	}

	// Update ke database
	if err := u.inventoriesRepo.Update(ctx, inventory); err != nil {
		u.logger.Error("Failed to update inventory in database",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return nil, fmt.Errorf("gagal memperbarui inventory: %w", err)
	}

	u.logger.Info("Inventory item updated successfully",
		zap.Int64("id", inventory.ID),
		zap.String("name", inventory.Name),
	)

	// Convert ke response
	return u.toInventoryResponse(inventory), nil
}

// DeleteInventory menghapus item inventory berdasarkan ID
func (u *inventoriesUsecase) DeleteInventory(ctx context.Context, id int64) error {
	u.logger.Info("Deleting inventory item",
		zap.Int64("id", id),
	)

	// Validasi ID
	if id <= 0 {
		u.logger.Warn("Validation failed: ID tidak valid", zap.Int64("id", id))
		return errors.New("ID tidak valid")
	}

	// Hapus dari database
	if err := u.inventoriesRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Failed to delete inventory from database",
			zap.Error(err),
			zap.Int64("id", id),
		)
		return fmt.Errorf("gagal menghapus inventory: %w", err)
	}

	u.logger.Info("Inventory item deleted successfully",
		zap.Int64("id", id),
	)

	return nil
}

// GetAllInventories mengambil semua item inventory dengan filter dan pagination
func (u *inventoriesUsecase) GetAllInventories(ctx context.Context, filter dto.InventoriesFilter) (*dto.InventoriesListResponse, error) {
	u.logger.Debug("Fetching all inventories with filters",
		zap.String("search", filter.Search),
		zap.String("unit", filter.Unit),
		zap.String("stock", filter.Stock),
		zap.Int("page", filter.Page),
		zap.Int("limit", filter.Limit),
	)

	// Set default pagination jika tidak ada
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	// Ambil data dari repository
	inventories, totalItems, err := u.inventoriesRepo.FindAll(ctx, filter)
	if err != nil {
		u.logger.Error("Failed to fetch all inventories from database",
			zap.Error(err),
			zap.String("search", filter.Search),
		)
		return nil, fmt.Errorf("gagal mengambil semua data inventory: %w", err)
	}

	u.logger.Info("All inventories fetched successfully",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(inventories)),
		zap.Int("page", filter.Page),
	)

	// Convert ke response
	var responses []dto.InventoriesResponse
	for _, inv := range inventories {
		responses = append(responses, *u.toInventoryResponse(&inv))
	}

	// Calculate total pages
	totalPages := int(totalItems) / filter.Limit
	if int(totalItems)%filter.Limit != 0 {
		totalPages++
	}

	pagination := dto.Pagination{
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
		TotalItems: int(totalItems),
	}

	return &dto.InventoriesListResponse{
		Data:       responses,
		Pagination: pagination,
	}, nil
}

// toInventoryResponse mengkonversi entity ke response DTO
func (u *inventoriesUsecase) toInventoryResponse(inv *entity.Inventories) *dto.InventoriesResponse {
	return &dto.InventoriesResponse{
		ID:          inv.ID,
		Image:       inv.Image,
		Name:        inv.Name,
		Category:    inv.Category,
		Quantity:    inv.Quantity,
		Status:      inv.Status,
		RetailPrice: inv.RetailPrice,
		CreatedAt:   inv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   inv.UpdatedAt.Format(time.RFC3339),
	}
}
