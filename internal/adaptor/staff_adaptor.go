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

type StaffAdaptor struct {
	service usecase.StaffUseCase
	logger  *zap.Logger
}

func NewStaffAdaptor(service usecase.StaffUseCase, logger *zap.Logger) *StaffAdaptor {
	return &StaffAdaptor{
		service: service,
		logger:  logger,
	}
}

// GetList menangani GET /staff
func (h *StaffAdaptor) GetList(c *gin.Context) {
	h.logger.Debug("GetList handler called",
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
	role := c.Query("role")

	// Construct DTO
	req := dto.StaffFilterRequest{
		Page:      page,
		Limit:     limit,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Role:      role,
	}

	result, pagination, err := h.service.GetListStaff(ctx, req)
	if err != nil {
		h.logger.Error("Failed to get staff list",
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

	h.logger.Info("GetList completed successfully",
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

// GetByID menangani GET /staff/{id}
func (h *StaffAdaptor) GetByID(c *gin.Context) {
	h.logger.Debug("GetByID handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid staff ID parameter",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid staff id",
			"data":    nil,
		})
		return
	}

	result, err := h.service.GetStaffByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			h.logger.Warn("Staff not found", zap.Uint64("id", id))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to get staff by ID",
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

	h.logger.Info("GetByID completed successfully",
		zap.Uint64("id", id),
		zap.String("name", result.FullName),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success get staff detail",
		"data":    result,
	})
}

// GetByEmail menangani GET /staff/email?email=xxx
func (h *StaffAdaptor) GetByEmail(c *gin.Context) {
	h.logger.Debug("GetByEmail handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get email from query parameter
	email := c.Query("email")
	if email == "" {
		h.logger.Warn("Empty staff email query parameter", zap.String("client_ip", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "email query parameter is required",
			"data":    nil,
		})
		return
	}

	result, err := h.service.GetStaffByEmail(ctx, email)
	if err != nil {
		if err.Error() == "staff not found" {
			h.logger.Warn("Staff not found by email", zap.String("email", email))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to get staff by email",
			zap.Error(err),
			zap.String("email", email),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("GetByEmail completed successfully",
		zap.String("email", email),
		zap.Uint("id", result.ID),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success get staff by email",
		"data":    result,
	})
}

// Create menangani POST /staff
func (h *StaffAdaptor) Create(c *gin.Context) {
	h.logger.Debug("Create handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Parse request body
	var req dto.StaffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create staff",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		h.logger.Warn("Validation error for create staff",
			zap.Any("validation_errors", err),
			zap.String("email", req.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "validation error",
			"errors":  err,
		})
		return
	}

	result, err := h.service.CreateStaff(ctx, req)
	if err != nil {
		if err.Error() == "email already exists" {
			h.logger.Warn("Email already exists",
				zap.String("email", req.Email),
			)
			c.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to create staff",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	h.logger.Info("Staff created successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.FullName),
		zap.String("email", result.Email),
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "success create staff",
		"data":    result,
	})
}

// Update menangani PUT /staff/{id}
func (h *StaffAdaptor) Update(c *gin.Context) {
	h.logger.Debug("Update handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid staff ID parameter for update",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid staff id",
			"data":    nil,
		})
		return
	}

	// Parse request body
	var req dto.StaffUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update staff",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		h.logger.Warn("Validation error for update staff",
			zap.Any("validation_errors", err),
			zap.Uint64("id", id),
			zap.String("email", req.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "validation error",
			"errors":  err,
		})
		return
	}

	result, err := h.service.UpdateStaff(ctx, uint(id), req)
	if err != nil {
		if err.Error() == "staff not found" {
			h.logger.Warn("Staff not found for update", zap.Uint64("id", id))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if err.Error() == "email already exists" {
			h.logger.Warn("Email already exists for another staff",
				zap.String("email", req.Email),
				zap.Uint64("id", id),
			)
			c.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to update staff",
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

	h.logger.Info("Staff updated successfully",
		zap.Uint("id", result.ID),
		zap.String("name", result.FullName),
		zap.String("email", result.Email),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success update staff",
		"data":    result,
	})
}

// Delete menangani DELETE /staff/{id}
func (h *StaffAdaptor) Delete(c *gin.Context) {
	h.logger.Debug("Delete handler called", zap.String("client_ip", c.ClientIP()))

	ctx := c.Request.Context()

	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid staff ID parameter for delete",
			zap.Error(err),
			zap.String("id_str", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid staff id",
			"data":    nil,
		})
		return
	}

	err = h.service.DeleteStaff(ctx, uint(id))
	if err != nil {
		if err.Error() == "staff not found" {
			h.logger.Warn("Staff not found for deletion", zap.Uint64("id", id))
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		h.logger.Error("Failed to delete staff",
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

	h.logger.Info("Staff deleted successfully", zap.Uint64("id", id))

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success delete staff",
		"data":    nil,
	})
}
