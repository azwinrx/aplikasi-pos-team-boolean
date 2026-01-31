package adaptor

import (
	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProductAdaptor struct {
	service usecase.ProductUseCase
	logger  *zap.Logger
}

func NewProductAdaptor(service usecase.ProductUseCase, logger *zap.Logger) *ProductAdaptor {
	return &ProductAdaptor{
		service: service,
		logger:  logger,
	}
}

// GetList menangani GET /products
func (h *ProductAdaptor) GetList(c *gin.Context) {
	h.logger.Debug("GetList product handler called",
		zap.String("client_ip", c.ClientIP()),
		zap.String("query_string", c.Request.URL.RawQuery),
	)

	ctx := c.Request.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	search := c.Query("search")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	status := c.Query("is_available")
	sortOrder := c.Query("sort_order")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	// Parse is_available filter
	var isAvailable *bool
	if status != "" {
		val := status == "true" || status == "1"
		isAvailable = &val
	}

	// Construct DTO
	req := dto.ProductFilterRequest{
		Page:        page,
		Limit:       limit,
		Search:      search,
		CategoryID:  uint(categoryID),
		IsAvailable: isAvailable,
		SortOrder:   sortOrder,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
	}

	result, pagination, err := h.service.GetListProduct(ctx, req)
	if err != nil {
		h.logger.Error("Failed to get product list",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("GetList product completed successfully",
		zap.Int("total_items", pagination.TotalItems),
		zap.Int("returned_items", len(result)),
		zap.Int("page", page),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"message":    "success get data",
		"data":       result,
		"pagination": pagination,
	})
}

// GetByCategory menangani GET /products/category/:category_id
func (h *ProductAdaptor) GetByCategory(c *gin.Context) {
	h.logger.Debug("GetByCategory product handler called",
		zap.String("client_ip", c.ClientIP()),
	)

	ctx := c.Request.Context()

	// Get category_id from URL parameter
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid category ID parameter",
			zap.Error(err),
			zap.String("category_id_str", categoryIDStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid category id",
			"data":    nil,
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	// Construct DTO with category filter
	req := dto.ProductFilterRequest{
		Page:       page,
		Limit:      limit,
		CategoryID: uint(categoryID),
		SortOrder:  "asc",
	}

	result, pagination, err := h.service.GetListProduct(ctx, req)
	if err != nil {
		h.logger.Error("Failed to get products by category",
			zap.Error(err),
			zap.Uint64("category_id", categoryID),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("GetByCategory product completed successfully",
		zap.Uint64("category_id", categoryID),
		zap.Int("total_items", pagination.TotalItems),
		zap.Int("returned_items", len(result)),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"message":    "success get products by category",
		"data":       result,
		"pagination": pagination,
	})
}

// GetByID menangani GET /products/{id}
func (h *ProductAdaptor) GetByID(c *gin.Context) {
	h.logger.Debug("GetByID product handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID parameter",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid product id",
			"data":    nil,
		})
		return
	}

	result, err := h.service.GetProductByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to get product by ID",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("GetByID product completed successfully",
		zap.Uint64("id", id),
		zap.String("name", result.ProductName),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success get product detail",
		"data":    result,
	})
}

// Create menangani POST /products
func (h *ProductAdaptor) Create(c *gin.Context) {
	h.logger.Debug("Create product handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Parse request body
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create product",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		h.logger.Warn("Validation failed for create product",
			zap.Any("validation_errors", err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "validation error",
			"errors":  err,
		})
		return
	}

	result, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to create product",
			zap.Error(err),
			zap.String("product_name", req.ProductName),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Product created successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.ProductName),
		zap.String("item_id", result.ItemID),
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "success create product",
		"data":    result,
	})
}

// Update menangani PUT /products/{id}
func (h *ProductAdaptor) Update(c *gin.Context) {
	h.logger.Debug("Update product handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID parameter for update",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid product id",
			"data":    nil,
		})
		return
	}

	// Parse request body
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update product",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request body: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		h.logger.Warn("Validation failed for update product",
			zap.Any("validation_errors", err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "validation error",
			"errors":  err,
		})
		return
	}

	result, err := h.service.UpdateProduct(ctx, uint(id), req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if err.Error() == "category not found" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to update product",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Product updated successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.ProductName),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success update product",
		"data":    result,
	})
}

// Delete menangani DELETE /products/{id}
func (h *ProductAdaptor) Delete(c *gin.Context) {
	h.logger.Debug("Delete product handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID parameter for delete",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid product id",
			"data":    nil,
		})
		return
	}

	err = h.service.DeleteProduct(ctx, uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to delete product",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Product deleted successfully", zap.Uint64("id", id))

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success delete product",
		"data":    nil,
	})
}
