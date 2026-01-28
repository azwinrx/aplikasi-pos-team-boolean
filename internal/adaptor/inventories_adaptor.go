package adaptor

import (
	"fmt"
	"net/http"

	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InventoriesAdaptor menangani semua request HTTP untuk inventories
type InventoriesAdaptor struct {
	inventoriesUsecase usecase.InventoriesUsecase
	logger             *zap.Logger
}

// NewInventoriesAdaptor membuat instance baru dari InventoriesAdaptor
func NewInventoriesAdaptor(inventoriesUsecase usecase.InventoriesUsecase, logger *zap.Logger) *InventoriesAdaptor {
	return &InventoriesAdaptor{
		inventoriesUsecase: inventoriesUsecase,
		logger:             logger,
	}
}

// GetAllInventories menangani request untuk mengambil semua inventory tanpa filter
func (h *InventoriesAdaptor) GetAllInventories(c *gin.Context) {
	h.logger.Debug("GetAllInventories handler called")

	var filter dto.InventoriesFilter

	// Default pagination
	filter.Page = 1
	filter.Limit = 10

	// Call usecase
	response, err := h.inventoriesUsecase.GetInventoryByFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get all inventories",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data inventory: "+err.Error())
		return
	}

	h.logger.Info("GetAllInventories completed successfully",
		zap.Int("total_items", response.Pagination.TotalItems),
		zap.Int("returned_items", len(response.Data)),
	)

	// Return response with pagination
	utils.ResponsePagination(c.Writer, http.StatusOK, "Data inventory berhasil diambil", response.Data, response.Pagination)
}

// GetInventoryByFilter menangani request untuk mengambil inventory dengan filter
// Support filter: search, status, category, stock, unit, min_qty, max_qty, min_price, max_price
func (h *InventoriesAdaptor) GetInventoryByFilter(c *gin.Context) {
	h.logger.Debug("GetInventoryByFilter handler called",
		zap.String("query_string", c.Request.URL.RawQuery),
	)

	var filter dto.InventoriesFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("Invalid query parameters",
			zap.Error(err),
			zap.String("query_string", c.Request.URL.RawQuery),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	h.logger.Debug("Filter parameters bound successfully",
		zap.String("search", filter.Search),
		zap.String("status", filter.Status),
		zap.String("category", filter.Category),
		zap.Int("page", filter.Page),
		zap.Int("limit", filter.Limit),
	)

	// Call usecase
	response, err := h.inventoriesUsecase.GetInventoryByFilter(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get inventory by filter",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data inventory: "+err.Error())
		return
	}

	h.logger.Info("GetInventoryByFilter completed successfully",
		zap.Int("total_items", response.Pagination.TotalItems),
		zap.Int("returned_items", len(response.Data)),
		zap.String("search", filter.Search),
	)

	// Return response with pagination
	utils.ResponsePagination(c.Writer, http.StatusOK, "Data inventory berhasil diambil", response.Data, response.Pagination)
}

// CreateInventory menangani request untuk membuat inventory baru
// Body: image, name, category, quantity, status, price
func (h *InventoriesAdaptor) CreateInventory(c *gin.Context) {
	h.logger.Debug("CreateInventory handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.InventoriesRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create inventory",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	h.logger.Debug("Request body bound successfully",
		zap.String("name", req.Name),
		zap.String("category", req.Category),
		zap.Int("quantity", req.Quantity),
		zap.String("status", req.Status),
	)

	// Call usecase
	response, err := h.inventoriesUsecase.CreateInventory(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create inventory",
			zap.Error(err),
			zap.String("name", req.Name),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal membuat inventory: "+err.Error())
		return
	}

	h.logger.Info("Inventory created successfully",
		zap.Int64("id", response.ID),
		zap.String("name", response.Name),
		zap.String("category", response.Category),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusCreated, "Inventory berhasil dibuat", response)
}

// UpdateInventory menangani request untuk update inventory yang sudah ada
func (h *InventoriesAdaptor) UpdateInventory(c *gin.Context) {
	h.logger.Debug("UpdateInventory handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from path parameter
	id := c.Param("id")
	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	var req dto.InventoriesRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update inventory",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	h.logger.Debug("Request body bound successfully",
		zap.Int64("id", idInt),
		zap.String("name", req.Name),
		zap.String("category", req.Category),
	)

	// Call usecase
	response, err := h.inventoriesUsecase.UpdateInventory(c.Request.Context(), idInt, req)
	if err != nil {
		h.logger.Error("Failed to update inventory",
			zap.Error(err),
			zap.Int64("id", idInt),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal memperbarui inventory: "+err.Error())
		return
	}

	h.logger.Info("Inventory updated successfully",
		zap.Int64("id", response.ID),
		zap.String("name", response.Name),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Inventory berhasil diperbarui", response)
}

// DeleteInventory menangani request untuk menghapus inventory
func (h *InventoriesAdaptor) DeleteInventory(c *gin.Context) {
	h.logger.Debug("DeleteInventory handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from path parameter
	id := c.Param("id")
	var idInt int64
	if _, err := fmt.Sscanf(id, "%d", &idInt); err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	h.logger.Debug("ID parameter bound successfully", zap.Int64("id", idInt))

	// Call usecase
	err := h.inventoriesUsecase.DeleteInventory(c.Request.Context(), idInt)
	if err != nil {
		h.logger.Error("Failed to delete inventory",
			zap.Error(err),
			zap.Int64("id", idInt),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal menghapus inventory: "+err.Error())
		return
	}

	h.logger.Info("Inventory deleted successfully", zap.Int64("id", idInt))

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Inventory berhasil dihapus", nil)
}
