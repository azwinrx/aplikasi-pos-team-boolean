package repository

import (
	"context"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(ctx context.Context, f dto.CategoryFilterRequest) ([]entity.Category, int64, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Detail(ctx context.Context, id uint) (*entity.Category, error)
	Delete(ctx context.Context, id uint) error
	FindByName(ctx context.Context, name string) (*entity.Category, error)
	CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error)
}

type categoryRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, logger *zap.Logger) CategoryRepository {
	return &categoryRepository{db, logger}
}

func (r *categoryRepository) FindAll(ctx context.Context, f dto.CategoryFilterRequest) ([]entity.Category, int64, error) {
	r.logger.Info("Finding all categories",
		zap.Int("page", f.Page),
		zap.Int("limit", f.Limit))

	var categories []entity.Category
	var totalItems int64

	query := r.db.WithContext(ctx).Model(&entity.Category{})

	// Search filter
	if f.Search != "" {
		searchPattern := "%" + f.Search + "%"
		query = query.Where("category_name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Count total items
	if err := query.Count(&totalItems).Error; err != nil {
		r.logger.Error("Failed to count category items", zap.Error(err))
		return nil, 0, err
	}

	// Sorting - default by category_name
	sortColumn := "category_name"
	if f.SortBy != "" {
		switch f.SortBy {
		case "created_at":
			sortColumn = "created_at"
		default:
			sortColumn = "category_name"
		}
	}

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

	err := query.Limit(f.Limit).Offset(offset).Find(&categories).Error
	if err != nil {
		r.logger.Error("Failed to find all categories", zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Successfully found all categories",
		zap.Int64("total_items", totalItems),
		zap.Int("returned_items", len(categories)))

	return categories, totalItems, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	r.logger.Info("Creating new category",
		zap.String("category_name", category.CategoryName))

	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		r.logger.Error("Failed to create category",
			zap.String("category_name", category.CategoryName),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully created category",
		zap.Uint("id", category.ID),
		zap.String("category_name", category.CategoryName))
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	r.logger.Info("Updating category",
		zap.Uint("id", category.ID),
		zap.String("category_name", category.CategoryName))

	err := r.db.WithContext(ctx).Save(category).Error
	if err != nil {
		r.logger.Error("Failed to update category",
			zap.Uint("id", category.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully updated category", zap.Uint("id", category.ID))
	return nil
}

func (r *categoryRepository) Detail(ctx context.Context, id uint) (*entity.Category, error) {
	r.logger.Info("Getting category detail", zap.Uint("id", id))

	var category entity.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		r.logger.Error("Failed to get category detail",
			zap.Uint("id", id),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully retrieved category detail",
		zap.Uint("id", id),
		zap.String("category_name", category.CategoryName))
	return &category, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting category", zap.Uint("id", id))

	// Soft delete
	err := r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete category",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted category", zap.Uint("id", id))
	return nil
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	r.logger.Info("Finding category by name", zap.String("name", name))

	var category entity.Category
	err := r.db.WithContext(ctx).Where("category_name = ?", name).First(&category).Error
	if err != nil {
		r.logger.Error("Failed to find category by name",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Successfully found category by name",
		zap.Uint("id", category.ID),
		zap.String("name", category.CategoryName))
	return &category, nil
}

func (r *categoryRepository) CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		r.logger.Error("Failed to count products by category",
			zap.Uint("category_id", categoryID),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}
