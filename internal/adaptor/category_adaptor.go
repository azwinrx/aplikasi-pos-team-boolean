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

type CategoryAdaptor struct {
	service usecase.CategoryUseCase
	logger  *zap.Logger
}

func NewCategoryAdaptor(service usecase.CategoryUseCase, logger *zap.Logger) *CategoryAdaptor {
	return &CategoryAdaptor{
		service: service,
		logger:  logger,
	}
}

// GetList menangani GET /categories
func (h *CategoryAdaptor) GetList(c *gin.Context) {
	h.logger.Debug("GetList category handler called",
		zap.String("client_ip", c.ClientIP()),
		zap.String("query_string", c.Request.URL.RawQuery),
	)

	ctx := c.Request.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	// Construct DTO
	req := dto.CategoryFilterRequest{
		Page:      page,
		Limit:     limit,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	result, pagination, err := h.service.GetListCategory(ctx, req)
	if err != nil {
		h.logger.Error("Failed to get category list",
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

	h.logger.Info("GetList category completed successfully",
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

// GetByID menangani GET /categories/{id}
func (h *CategoryAdaptor) GetByID(c *gin.Context) {
	h.logger.Debug("GetByID category handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid category ID parameter",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid category id",
			"data":    nil,
		})
		return
	}

	result, err := h.service.GetCategoryByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to get category by ID",
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

	h.logger.Info("GetByID category completed successfully",
		zap.Uint64("id", id),
		zap.String("name", result.CategoryName),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success get category detail",
		"data":    result,
	})
}

// Create menangani POST /categories
func (h *CategoryAdaptor) Create(c *gin.Context) {
	h.logger.Debug("Create category handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Parse request body
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create category",
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
		h.logger.Warn("Validation failed for create category",
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

	result, err := h.service.CreateCategory(ctx, req)
	if err != nil {
		if err.Error() == "category name already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to create category",
			zap.Error(err),
			zap.String("category_name", req.CategoryName),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Category created successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.CategoryName),
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "success create category",
		"data":    result,
	})
}

// Update menangani PUT /categories/{id}
func (h *CategoryAdaptor) Update(c *gin.Context) {
	h.logger.Debug("Update category handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid category ID parameter for update",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid category id",
			"data":    nil,
		})
		return
	}

	// Parse request body
	var req dto.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update category",
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
		h.logger.Warn("Validation failed for update category",
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

	result, err := h.service.UpdateCategory(ctx, uint(id), req)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if err.Error() == "category name already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to update category",
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

	h.logger.Info("Category updated successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.CategoryName),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success update category",
		"data":    result,
	})
}

// Delete menangani DELETE /categories/{id}
func (h *CategoryAdaptor) Delete(c *gin.Context) {
	h.logger.Debug("Delete category handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid category ID parameter for delete",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid category id",
			"data":    nil,
		})
		return
	}

	err = h.service.DeleteCategory(ctx, uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if err.Error() == "cannot delete category with existing products" {
			c.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to delete category",
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

	h.logger.Info("Category deleted successfully", zap.Uint64("id", id))

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success delete category",
		"data":    nil,
	})
}
