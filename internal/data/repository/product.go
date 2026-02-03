package repository

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(ctx context.Context, f dto.ProductFilterRequest) ([]entity.Product, int64, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Detail(ctx context.Context, id uint) (*entity.Product, error)
	Delete(ctx context.Context, id uint) error
	FindByItemID(ctx context.Context, itemID string) (*entity.Product, error)
}

type productRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProductRepository(db *gorm.DB, logger *zap.Logger) ProductRepository {
	return &productRepository{db, logger}
}

func (r *productRepository) FindAll(ctx context.Context, f dto.ProductFilterRequest) ([]entity.Product, int64, error) {
	r.logger.Info("Finding all products",
		zap.Int("page", f.Page),
		zap.Int("limit", f.Limit),
		zap.Uint("category_id", f.CategoryID))

	var products []entity.Product
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	// Search filter
	if f.Search != "" {
		searchPattern := "%" + f.Search + "%"
		query = query.Where("product_name ILIKE ? OR item_id ILIKE ?", searchPattern, searchPattern)
	}

	// Category filter
	if f.CategoryID > 0 {
		query = query.Where("category_id = ?", f.CategoryID)
	}

	// Availability filter
	if f.IsAvailable != nil {
		query = query.Where("is_available = ?", *f.IsAvailable)
	}

	// Price range filter
	if f.MinPrice > 0 {
		query = query.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice > 0 {
		query = query.Where("price <= ?", f.MaxPrice)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		r.logger.Error("Failed to count product items", zap.Error(err))
		return nil, 0, err
	}

	// Sorting - always by product_name
	sortColumn := "product_name"

	sortOrder := "asc"
	if f.SortOrder == "desc" {
		sortOrder = "desc"
	}

	query = query.Order(sortColumn + " " + sortOrder)

	// Pagination
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}

	offset := (f.Page - 1) * f.Limit

	// Preload Category
	err := query.Preload("Category").Limit(f.Limit).Offset(offset).Find(&products).Error
	if err != nil {
		r.logger.Error("Failed to find all products", zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Successfully found all products",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(products)))

	return products, totalItems, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	r.logger.Info("Creating new product",
		zap.String("product_name", product.ProductName),
		zap.String("item_id", product.ItemID))

	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		r.logger.Error("Failed to create product",
			zap.String("product_name", product.ProductName),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully created product",
		zap.Uint("id", product.ID),
		zap.String("product_name", product.ProductName))
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	r.logger.Info("Updating product",
		zap.Uint("id", product.ID),
		zap.String("product_name", product.ProductName))

	err := r.db.WithContext(ctx).Save(product).Error
	if err != nil {
		r.logger.Error("Failed to update product",
			zap.Uint("id", product.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully updated product", zap.Uint("id", product.ID))
	return nil
}

func (r *productRepository) Detail(ctx context.Context, id uint) (*entity.Product, error) {
	r.logger.Info("Getting product detail", zap.Uint("id", id))

	var product entity.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error
	if err != nil {
		r.logger.Error("Failed to get product detail",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully retrieved product detail",
		zap.Uint("id", id),
		zap.String("product_name", product.ProductName))
	return &product, nil
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting product", zap.Uint("id", id))

	// Soft delete
	err := r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete product",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted product", zap.Uint("id", id))
	return nil
}

func (r *productRepository) FindByItemID(ctx context.Context, itemID string) (*entity.Product, error) {
	r.logger.Info("Finding product by item_id", zap.String("item_id", itemID))

	var product entity.Product
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&product).Error
	if err != nil {
		r.logger.Error("Failed to find product by item_id",
			zap.String("item_id", itemID),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found product by item_id",
		zap.Uint("id", product.ID),
		zap.String("item_id", product.ItemID))
	return &product, nil
}
