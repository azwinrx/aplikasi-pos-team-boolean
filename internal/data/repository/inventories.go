package repository

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"
	"context"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoriesRepository interface {
	Create(ctx context.Context, inventories *entity.Inventories) error
	Update(ctx context.Context, inventories *entity.Inventories) error
	Delete(ctx context.Context, id int64) error
	FindByFilter(ctx context.Context, filter dto.InventoriesFilter) ([]entity.Inventories, int64, error)
	FindAll(ctx context.Context, filter dto.InventoriesFilter) ([]entity.Inventories, int64, error)
}

// inventoriesRepository adalah implementasi dari interface InventoriesRepository
type inventoriesRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewInventoriesRepository membuat instance baru dari InventoriesRepository
func NewInventoriesRepository(db *gorm.DB, logger *zap.Logger) InventoriesRepository {
	return &inventoriesRepository{db: db, logger: logger}
}

// Create menambahkan item Inventories baru ke database
func (r *inventoriesRepository) Create(ctx context.Context, inventories *entity.Inventories) error {
	r.logger.Info("Creating new inventory item",
		zap.String("name", inventories.Name),
		zap.String("category", inventories.Category),
		zap.Int("quantity", inventories.Quantity))

	err := r.db.WithContext(ctx).Create(inventories).Error
	if err != nil {
		r.logger.Error("Failed to create inventory item",
			zap.String("name", inventories.Name),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully created inventory item",
		zap.Int64("id", inventories.ID),
		zap.String("name", inventories.Name))
	return err
}

// Update memperbarui item Inventories yang sudah ada
func (r *inventoriesRepository) Update(ctx context.Context, inventories *entity.Inventories) error {
	r.logger.Info("Updating inventory item",
		zap.Int64("id", inventories.ID),
		zap.String("name", inventories.Name))

	err := r.db.WithContext(ctx).Save(inventories).Error
	if err != nil {
		r.logger.Error("Failed to update inventory item",
			zap.Int64("id", inventories.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully updated inventory item",
		zap.Int64("id", inventories.ID))
	return err
}

// Delete menghapus item Inventories berdasarkan ID (soft delete)
func (r *inventoriesRepository) Delete(ctx context.Context, id int64) error {
	r.logger.Info("Deleting inventory item",
		zap.Int64("id", id))

	err := r.db.WithContext(ctx).Delete(&entity.Inventories{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete inventory item",
			zap.Int64("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted inventory item",
		zap.Int64("id", id))
	return err
}

// FindByFilter mengambil item Inventories berdasarkan parameter filter query
func (r *inventoriesRepository) FindByFilter(ctx context.Context, filter dto.InventoriesFilter) ([]entity.Inventories, int64, error) {
	r.logger.Info("Finding inventories by filter",
		zap.String("search", filter.Search),
		zap.String("status", filter.Status),
		zap.String("category", filter.Category),
		zap.Int("page", filter.Page),
		zap.Int("limit", filter.Limit))

	var inventories []entity.Inventories
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Inventories{})

	// Search filter - pencarian berdasarkan nama
	if filter.Search != "" {
		searchPattern := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchPattern)
	}

	// Status filter - filter berdasarkan status (active/inactive)
	if filter.Status != "" {
		query = query.Where("LOWER(status) = ?", strings.ToLower(filter.Status))
	}

	// Category filter - filter berdasarkan kategori
	if filter.Category != "" {
		query = query.Where("LOWER(category) = ?", strings.ToLower(filter.Category))
	}

	// Unit filter - filter berdasarkan satuan (Value)
	if filter.Unit != "" {
		query = query.Where("LOWER(unit) = ?", strings.ToLower(filter.Unit))
	}

	// Stock status filter - filter berdasarkan status stok
	if filter.Stock != "" {
		switch strings.ToLower(filter.Stock) {
		case "instock":
			query = query.Where("quantity >= min_stock AND quantity > 0")
		case "lowstock":
			query = query.Where("quantity < min_stock AND quantity > 0")
		case "outofstock":
			query = query.Where("quantity = 0")
		}
	}

	// Quantity range filter - filter berdasarkan range quantity (Piece/Item/Quantity)
	if filter.MinQty > 0 {
		query = query.Where("quantity >= ?", filter.MinQty)
	}
	if filter.MaxQty > 0 {
		query = query.Where("quantity <= ?", filter.MaxQty)
	}

	// Price range filter - filter berdasarkan range harga
	if filter.MinPrice > 0 {
		query = query.Where("retail_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("retail_price <= ?", filter.MaxPrice)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		r.logger.Error("Failed to count inventory items",
			zap.Error(err))
		return nil, 0, err
	}

	// Sorting
	sortColumn := "created_at"
	validSortColumns := map[string]bool{
		"name":         true,
		"quantity":     true,
		"retail_price": true,
		"created_at":   true,
		"category":     true,
		"status":       true,
	}
	if validSortColumns[filter.SortBy] {
		sortColumn = filter.SortBy
	}

	sortOrder := "desc"
	if filter.SortDir == "asc" {
		sortOrder = "asc"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Limit(filter.Limit).Offset(offset).Find(&inventories).Error

	if err != nil {
		r.logger.Error("Failed to find inventories by filter",
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Successfully found inventories by filter",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(inventories)))

	return inventories, totalItems, err
}

// FindAll mengambil semua item Inventories dengan filter dan pagination
func (r *inventoriesRepository) FindAll(ctx context.Context, filter dto.InventoriesFilter) ([]entity.Inventories, int64, error) {
	r.logger.Info("Finding all inventories",
		zap.String("search", filter.Search),
		zap.Int("page", filter.Page),
		zap.Int("limit", filter.Limit))

	var inventories []entity.Inventories
	var totalItems int64

	query := r.db.Model(&entity.Inventories{})

	// Search filter
	if filter.Search != "" {
		searchPattern := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchPattern)
	}

	// Unit filter
	if filter.Unit != "" {
		query = query.Where("LOWER(unit) = ?", strings.ToLower(filter.Unit))
	}

	// Stock status filter
	if filter.Stock != "" {
		switch strings.ToLower(filter.Stock) {
		case "instock":
			query = query.Where("quantity >= min_stock AND quantity > 0")
		case "lowstock":
			query = query.Where("quantity < min_stock AND quantity > 0")
		case "outofstock":
			query = query.Where("quantity = 0")
		}
	}

	// Quantity range filter
	if filter.MinQty > 0 {
		query = query.Where("quantity >= ?", filter.MinQty)
	}
	if filter.MaxQty > 0 {
		query = query.Where("quantity <= ?", filter.MaxQty)
	}

	// Price range filter
	if filter.MinPrice > 0 {
		query = query.Where("retail_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("retail_price <= ?", filter.MaxPrice)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		r.logger.Error("Failed to count all inventory items",
			zap.Error(err))
		return nil, 0, err
	}

	// Sorting
	sortColumn := "created_at"
	validSortColumns := map[string]bool{
		"name":         true,
		"quantity":     true,
		"retail_price": true,
		"created_at":   true,
		"category":     true,
		"status":       true,
	}
	if validSortColumns[filter.SortBy] {
		sortColumn = filter.SortBy
	}

	sortOrder := "desc"
	if filter.SortDir == "asc" {
		sortOrder = "asc"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	err := query.Limit(filter.Limit).Offset(offset).Find(&inventories).Error

	if err != nil {
		r.logger.Error("Failed to find all inventories",
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Successfully found all inventories",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(inventories)))

	return inventories, totalItems, err
}

// GetSummary mengambil summary/statistik inventory
func (r *inventoriesRepository) GetSummary(ctx context.Context) (*dto.InventoriesSummary, error) {
	r.logger.Info("Getting inventory summary")

	var summary dto.InventoriesSummary
	var count int64

	// Total products
	if err := r.db.WithContext(ctx).Model(&entity.Inventories{}).Count(&count).Error; err != nil {
		r.logger.Error("Failed to count total products", zap.Error(err))
		return nil, err
	}
	summary.TotalProducts = int(count)

	// Active products
	if err := r.db.WithContext(ctx).Model(&entity.Inventories{}).Where("status = ?", "active").Count(&count).Error; err != nil {
		r.logger.Error("Failed to count active products", zap.Error(err))
		return nil, err
	}
	summary.ActiveProducts = int(count)

	// Inactive products
	if err := r.db.WithContext(ctx).Model(&entity.Inventories{}).Where("status = ?", "inactive").Count(&count).Error; err != nil {
		r.logger.Error("Failed to count inactive products", zap.Error(err))
		return nil, err
	}
	summary.InactiveProducts = int(count)

	// Low stock products
	if err := r.db.WithContext(ctx).Model(&entity.Inventories{}).Where("quantity < min_stock AND quantity > 0").Count(&count).Error; err != nil {
		r.logger.Error("Failed to count low stock products", zap.Error(err))
		return nil, err
	}
	summary.LowStockProducts = int(count)

	// Out of stock products
	if err := r.db.WithContext(ctx).Model(&entity.Inventories{}).Where("quantity = 0").Count(&count).Error; err != nil {
		r.logger.Error("Failed to count out of stock products", zap.Error(err))
		return nil, err
	}
	summary.OutOfStockProducts = int(count)

	r.logger.Info("Successfully retrieved inventory summary",
		zap.Int("total_products", summary.TotalProducts),
		zap.Int("active_products", summary.ActiveProducts),
		zap.Int("low_stock_products", summary.LowStockProducts),
		zap.Int("out_of_stock_products", summary.OutOfStockProducts))

	return &summary, nil
}
