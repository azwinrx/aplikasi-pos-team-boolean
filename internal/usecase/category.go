package usecase

import (
	"context"
	"errors"
	"math"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryUseCase interface {
	GetListCategory(ctx context.Context, req dto.CategoryFilterRequest) ([]dto.CategoryListResponse, dto.Pagination, error)
	GetCategoryByID(ctx context.Context, id uint) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, req dto.CategoryCreateRequest) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req dto.CategoryUpdateRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error
}

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository, logger *zap.Logger) *categoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *categoryUseCase) GetListCategory(ctx context.Context, req dto.CategoryFilterRequest) ([]dto.CategoryListResponse, dto.Pagination, error) {
	s.logger.Debug("Fetching category list",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit),
	)

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "category_name"
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Get data from repository
	categories, total, err := s.categoryRepo.FindAll(ctx, req)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	s.logger.Info("Category list fetched successfully",
		zap.Int64("total_items", total),
		zap.Int("returned_items", len(categories)),
		zap.Int("page", req.Page),
	)

	// Convert to list response
	var categoryResponses []dto.CategoryListResponse
	for _, category := range categories {
		// Count products in this category
		productCount, _ := s.categoryRepo.CountProductsByCategory(ctx, category.ID)
		categoryResponses = append(categoryResponses, dto.CategoryListResponse{
			ID:           category.ID,
			IconCategory: category.IconCategory,
			CategoryName: category.CategoryName,
			ProductCount: int(productCount),
		})
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := dto.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		TotalItems: int(total),
	}

	return categoryResponses, pagination, nil
}

func (s *categoryUseCase) GetCategoryByID(ctx context.Context, id uint) (*dto.CategoryResponse, error) {
	s.logger.Debug("Fetching category by ID", zap.Uint("id", id))

	category, err := s.categoryRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	s.logger.Info("Category fetched successfully", zap.Uint("id", id), zap.String("category_name", category.CategoryName))

	// Count products
	productCount, _ := s.categoryRepo.CountProductsByCategory(ctx, category.ID)

	response := s.toCategoryResponse(category, int(productCount))
	return &response, nil
}

func (s *categoryUseCase) CreateCategory(ctx context.Context, req dto.CategoryCreateRequest) (*dto.CategoryResponse, error) {
	s.logger.Info("Creating new category",
		zap.String("category_name", req.CategoryName),
	)

	// Validate required fields
	if req.CategoryName == "" {
		return nil, errors.New("category name is required")
	}

	// Check if category name already exists
	existingCategory, err := s.categoryRepo.FindByName(ctx, req.CategoryName)
	if err == nil && existingCategory != nil {
		return nil, errors.New("category name already exists")
	}

	// Create category entity
	category := &entity.Category{
		IconCategory: req.IconCategory,
		CategoryName: req.CategoryName,
		Description:  req.Description,
	}

	// Save to database
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		s.logger.Error("Failed to create category",
			zap.String("category_name", req.CategoryName),
			zap.Error(err),
		)
		return nil, errors.New("failed to create category")
	}

	s.logger.Info("Category created successfully",
		zap.Uint("id", category.ID),
		zap.String("category_name", category.CategoryName),
	)

	response := s.toCategoryResponse(category, 0)
	return &response, nil
}

func (s *categoryUseCase) UpdateCategory(ctx context.Context, id uint, req dto.CategoryUpdateRequest) (*dto.CategoryResponse, error) {
	s.logger.Info("Updating category",
		zap.Uint("id", id),
		zap.String("category_name", req.CategoryName),
	)

	// Get existing category
	category, err := s.categoryRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Check if category name already exists (exclude current category)
	if req.CategoryName != category.CategoryName {
		existingCategory, err := s.categoryRepo.FindByName(ctx, req.CategoryName)
		if err == nil && existingCategory != nil && existingCategory.ID != id {
			return nil, errors.New("category name already exists")
		}
	}

	// Update category data
	category.IconCategory = req.IconCategory
	category.CategoryName = req.CategoryName
	category.Description = req.Description

	// Save to database
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		s.logger.Error("Failed to update category",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return nil, errors.New("failed to update category")
	}

	s.logger.Info("Category updated successfully",
		zap.Uint("id", category.ID),
		zap.String("category_name", category.CategoryName),
	)

	productCount, _ := s.categoryRepo.CountProductsByCategory(ctx, category.ID)
	response := s.toCategoryResponse(category, int(productCount))
	return &response, nil
}

func (s *categoryUseCase) DeleteCategory(ctx context.Context, id uint) error {
	s.logger.Info("Deleting category", zap.Uint("id", id))

	// Check if category exists
	category, err := s.categoryRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	// Check if category has products
	productCount, _ := s.categoryRepo.CountProductsByCategory(ctx, category.ID)
	if productCount > 0 {
		return errors.New("cannot delete category with existing products")
	}

	// Soft delete
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete category",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return errors.New("failed to delete category")
	}

	s.logger.Info("Category deleted successfully", zap.Uint("id", id))
	return nil
}

// toCategoryResponse converts entity to response DTO
func (s *categoryUseCase) toCategoryResponse(category *entity.Category, productCount int) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:           category.ID,
		IconCategory: category.IconCategory,
		CategoryName: category.CategoryName,
		Description:  category.Description,
		ProductCount: productCount,
		CreatedAt:    category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
