package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductUseCase interface {
	GetListProduct(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductListResponse, dto.Pagination, error)
	GetProductByID(ctx context.Context, id uint) (*dto.ProductResponse, error)
	CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uint, req dto.ProductUpdateRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error
}

type productUseCase struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	logger       *zap.Logger
}

func NewProductUseCase(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, logger *zap.Logger) *productUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *productUseCase) GetListProduct(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductListResponse, dto.Pagination, error) {
	s.logger.Debug("Fetching product list",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit),
		zap.Uint("category_id", req.CategoryID),
	)

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Get data from repository
	products, total, err := s.productRepo.FindAll(ctx, req)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	s.logger.Info("Product list fetched successfully",
		zap.Int64("total_items", total),
		zap.Int("returned_items", len(products)),
		zap.Int("page", req.Page),
	)

	// Convert to list response
	var productResponses []dto.ProductListResponse
	for _, product := range products {
		productResponses = append(productResponses, s.toProductListResponse(&product))
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := dto.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		TotalItems: int(total),
	}

	return productResponses, pagination, nil
}

func (s *productUseCase) GetProductByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	s.logger.Debug("Fetching product by ID", zap.Uint("id", id))

	product, err := s.productRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	s.logger.Info("Product fetched successfully", zap.Uint("id", id), zap.String("product_name", product.ProductName))

	response := s.toProductResponse(product)
	return &response, nil
}

func (s *productUseCase) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (*dto.ProductResponse, error) {
	s.logger.Info("Creating new product",
		zap.String("product_name", req.ProductName),
		zap.Uint("category_id", req.CategoryID),
	)

	// Validate required fields
	if req.ProductName == "" {
		return nil, errors.New("product name is required")
	}
	if req.CategoryID == 0 {
		return nil, errors.New("category is required")
	}

	// Verify category exists
	category, err := s.categoryRepo.Detail(ctx, req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	// Generate unique Item ID
	itemID := s.generateItemID()

	// Create product entity
	product := &entity.Product{
		ProductImage: req.ProductImage,
		ProductName:  req.ProductName,
		ItemID:       itemID,
		Stock:        req.Stock,
		CategoryID:   req.CategoryID,
		Price:        req.Price,
	}

	// Save to database
	if err := s.productRepo.Create(ctx, product); err != nil {
		s.logger.Error("Failed to create product",
			zap.String("product_name", req.ProductName),
			zap.Error(err),
		)
		return nil, errors.New("failed to create product")
	}

	// Set category for response
	product.Category = *category

	s.logger.Info("Product created successfully",
		zap.Uint("id", product.ID),
		zap.String("product_name", product.ProductName),
		zap.String("item_id", product.ItemID),
	)

	response := s.toProductResponse(product)
	return &response, nil
}

func (s *productUseCase) UpdateProduct(ctx context.Context, id uint, req dto.ProductUpdateRequest) (*dto.ProductResponse, error) {
	s.logger.Info("Updating product",
		zap.Uint("id", id),
		zap.String("product_name", req.ProductName),
	)

	// Get existing product
	product, err := s.productRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Verify category exists if changed
	if req.CategoryID != product.CategoryID {
		category, err := s.categoryRepo.Detail(ctx, req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, err
		}
		product.Category = *category
	}

	// Update product data
	product.ProductImage = req.ProductImage
	product.ProductName = req.ProductName
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.Price = req.Price

	// Save to database
	if err := s.productRepo.Update(ctx, product); err != nil {
		s.logger.Error("Failed to update product",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return nil, errors.New("failed to update product")
	}

	// Reload to get updated category
	product, _ = s.productRepo.Detail(ctx, id)

	s.logger.Info("Product updated successfully",
		zap.Uint("id", product.ID),
		zap.String("product_name", product.ProductName),
	)

	response := s.toProductResponse(product)
	return &response, nil
}

func (s *productUseCase) DeleteProduct(ctx context.Context, id uint) error {
	s.logger.Info("Deleting product", zap.Uint("id", id))

	// Check if product exists
	_, err := s.productRepo.Detail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Soft delete
	if err := s.productRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete product",
			zap.Uint("id", id),
			zap.Error(err),
		)
		return errors.New("failed to delete product")
	}

	s.logger.Info("Product deleted successfully", zap.Uint("id", id))
	return nil
}

// generateItemID generates a unique item ID like #22314644
func (s *productUseCase) generateItemID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("#%08d", rand.Intn(100000000))
}

// toProductListResponse converts entity to list response DTO
func (s *productUseCase) toProductListResponse(product *entity.Product) dto.ProductListResponse {
	return dto.ProductListResponse{
		ID:           product.ID,
		ProductImage: product.ProductImage,
		ProductName:  product.ProductName,
		ItemID:       product.ItemID,
		Stock:        product.Stock,
		CategoryName: product.Category.CategoryName,
		Price:        product.Price,
		IsAvailable:  product.IsAvailable,
		Availability: product.GetAvailabilityStatus(),
	}
}

// toProductResponse converts entity to response DTO
func (s *productUseCase) toProductResponse(product *entity.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:           product.ID,
		ProductImage: product.ProductImage,
		ProductName:  product.ProductName,
		ItemID:       product.ItemID,
		Stock:        product.Stock,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.CategoryName,
		Price:        product.Price,
		IsAvailable:  product.IsAvailable,
		Availability: product.GetAvailabilityStatus(),
		CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
