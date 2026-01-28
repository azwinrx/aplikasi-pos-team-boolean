package adaptor

import (
	"net/http"
	"strconv"

	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
)

// InventoryHandler menangani semua request HTTP untuk inventory
type InventoryHandler struct {
	inventoryUsecase usecase.InventoryUsecase
}

// NewInventoryHandler membuat instance baru dari InventoryHandler
func NewInventoryHandler(inventoryUsecase usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUsecase: inventoryUsecase,
	}
}

// CreateInventory menangani request untuk membuat inventory baru
// @Summary Create inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param request body dto.InventoryRequest true "Inventory Request"
// @Success 201 {object} dto.InventoryResponse
// @Router /inventory [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.InventoryRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Call usecase
	response, err := h.inventoryUsecase.CreateInventory(c.Request.Context(), req)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal membuat inventory: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusCreated, "Inventory berhasil dibuat", response)
}

// GetAllInventories menangani request untuk mengambil semua inventory dengan filter
// @Summary Get all inventories
// @Tags Inventory
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param stock query string false "Filter by stock (instock/lowstock/outofstock)"
// @Param unit query string false "Filter by unit"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_dir query string false "Sort direction (asc/desc)" default(desc)
// @Success 200 {object} dto.InventoryListResponse
// @Router /inventory [get]
func (h *InventoryHandler) GetAllInventories(c *gin.Context) {
	var filter dto.InventoryFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	// Call usecase
	response, err := h.inventoryUsecase.GetAllInventories(c.Request.Context(), filter)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data inventory: "+err.Error())
		return
	}

	// Return response with pagination
	utils.ResponsePagination(c.Writer, http.StatusOK, "Data inventory berhasil diambil", response.Data, response.Pagination)
}

// GetInventoryByID menangani request untuk mengambil detail inventory berdasarkan ID
// @Summary Get inventory by ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Success 200 {object} dto.InventoryResponse
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetInventoryByID(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	// Call usecase
	response, err := h.inventoryUsecase.GetInventoryByID(c.Request.Context(), id)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusNotFound, "Inventory tidak ditemukan: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data inventory berhasil diambil", response)
}

// UpdateInventory menangani request untuk update inventory
// @Summary Update inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Param request body dto.InventoryRequest true "Inventory Request"
// @Success 200 {object} dto.InventoryResponse
// @Router /inventory/{id} [put]
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	var req dto.InventoryRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Call usecase
	response, err := h.inventoryUsecase.UpdateInventory(c.Request.Context(), id, req)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal memperbarui inventory: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Inventory berhasil diperbarui", response)
}

// DeleteInventory menangani request untuk menghapus inventory
// @Summary Delete inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Success 200
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	// Call usecase
	if err := h.inventoryUsecase.DeleteInventory(c.Request.Context(), id); err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal menghapus inventory: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Inventory berhasil dihapus", nil)
}

// AddStock menangani request untuk menambah stok
// @Summary Add stock
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Param amount body int true "Amount to add"
// @Success 200
// @Router /inventory/{id}/add-stock [post]
func (h *InventoryHandler) AddStock(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	// Parse request body
	var req struct {
		Amount int `json:"amount" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Call usecase
	if err := h.inventoryUsecase.AddStock(c.Request.Context(), id, req.Amount); err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal menambah stok: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Stok berhasil ditambahkan", nil)
}

// ReduceStock menangani request untuk mengurangi stok
// @Summary Reduce stock
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Param amount body int true "Amount to reduce"
// @Success 200
// @Router /inventory/{id}/reduce-stock [post]
func (h *InventoryHandler) ReduceStock(c *gin.Context) {
	// Parse ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	// Parse request body
	var req struct {
		Amount int `json:"amount" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Call usecase
	if err := h.inventoryUsecase.ReduceStock(c.Request.Context(), id, req.Amount); err != nil {
		utils.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Stok berhasil dikurangi", nil)
}

// GetLowStockItems menangani request untuk mengambil semua item dengan stok rendah
// @Summary Get low stock items
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {array} dto.InventoryResponse
// @Router /inventory/low-stock [get]
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	// Call usecase
	response, err := h.inventoryUsecase.GetLowStockItems(c.Request.Context())
	if err != nil {
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data low stock: "+err.Error())
		return
	}

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data low stock berhasil diambil", response)
}
