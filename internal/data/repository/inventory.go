package repository

import (
	"context"
	"fmt"
	"strings"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"gorm.io/gorm"
)

// inventoryRepository adalah implementasi dari interface InventoryRepository
type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository membuat instance baru dari InventoryRepository
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Create menambahkan item inventory baru ke database
func (r *inventoryRepository) Create(ctx context.Context, inventory *entity.Inventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

// Update memperbarui item inventory yang sudah ada
func (r *inventoryRepository) Update(ctx context.Context, inventory *entity.Inventory) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

// Delete menghapus item inventory berdasarkan ID (soft delete)
func (r *inventoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Inventory{}, id).Error
}

// FindByID mengambil satu item inventory berdasarkan ID
func (r *inventoryRepository) FindByID(ctx context.Context, id int64) (*entity.Inventory, error) {
	var inventory entity.Inventory
	err := r.db.WithContext(ctx).First(&inventory, id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// FindAll mengambil semua item inventory dengan filter dan pagination
func (r *inventoryRepository) FindAll(ctx context.Context, filter dto.InventoryFilter) ([]entity.Inventory, *dto.Pagination, error) {
	var inventories []entity.Inventory

	// Build query dengan filter
	query := r.buildFilterQuery(r.db.WithContext(ctx), filter)

	// Hitung total item yang sesuai dengan filter
	var totalItems int64
	if err := query.Model(&entity.Inventory{}).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	// Terapkan pagination
	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// Terapkan sorting
	sortBy := filter.SortBy
	sortDir := strings.ToUpper(filter.SortDir)
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortDir != "ASC" && sortDir != "DESC" {
		sortDir = "DESC"
	}

	// Eksekusi query dengan pagination dan sorting
	err := query.
		Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
		Limit(limit).
		Offset(offset).
		Find(&inventories).Error

	if err != nil {
		return nil, nil, err
	}

	pagination := &dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: int(totalItems),
	}

	return inventories, pagination, nil
}

// Count menghitung total item inventory yang sesuai dengan filter
func (r *inventoryRepository) Count(ctx context.Context, filter dto.InventoryFilter) (int64, error) {
	var count int64
	query := r.buildFilterQuery(r.db.WithContext(ctx), filter)
	err := query.Model(&entity.Inventory{}).Count(&count).Error
	return count, err
}

// UpdateStock memperbarui hanya quantity dari item inventory
func (r *inventoryRepository) UpdateStock(ctx context.Context, id int64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&entity.Inventory{}).
		Where("id = ?", id).
		Update("quantity", quantity).Error
}

// GetLowStockItems mengambil semua item dengan quantity di bawah min_stock
func (r *inventoryRepository) GetLowStockItems(ctx context.Context) ([]entity.Inventory, error) {
	var inventories []entity.Inventory
	err := r.db.WithContext(ctx).
		Where("quantity < min_stock").
		Find(&inventories).Error
	return inventories, err
}

// buildFilterQuery membangun GORM query dengan semua filter yang berlaku
func (r *inventoryRepository) buildFilterQuery(db *gorm.DB, filter dto.InventoryFilter) *gorm.DB {
	query := db

	// Filter pencarian - cari berdasarkan nama
	if filter.Search != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Search)+"%")
	}

	// Filter satuan
	if filter.Unit != "" {
		query = query.Where("LOWER(unit) = ?", strings.ToLower(filter.Unit))
	}

	// Filter status stok
	if filter.Stock != "" {
		switch strings.ToLower(filter.Stock) {
		case "instock":
			query = query.Where("quantity >= min_stock")
		case "lowstock":
			query = query.Where("quantity < min_stock AND quantity > 0")
		case "outofstock":
			query = query.Where("quantity = 0")
		}
	}

	// Filter rentang harga (jika ada field price di masa depan)
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	return query
}
