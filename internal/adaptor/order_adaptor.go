package adaptor

import (
	"net/http"
	"strconv"

	"aplikasi-pos-team-boolean/internal/dto"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderAdaptor menangani semua request HTTP untuk orders
type OrderAdaptor struct {
	orderUsecase usecase.OrderUseCase
	logger       *zap.Logger
}

// NewOrderAdaptor membuat instance baru dari OrderAdaptor
func NewOrderAdaptor(orderUsecase usecase.OrderUseCase, logger *zap.Logger) *OrderAdaptor {
	return &OrderAdaptor{
		orderUsecase: orderUsecase,
		logger:       logger,
	}
}

// GetAllOrders menangani request untuk mengambil semua order
func (h *OrderAdaptor) GetAllOrders(c *gin.Context) {
	h.logger.Debug("GetAllOrders handler called")

	// Call usecase
	response, err := h.orderUsecase.GetAllOrders(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get all orders",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data order: "+err.Error())
		return
	}

	h.logger.Info("GetAllOrders completed successfully",
		zap.Int("total_orders", len(response)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data order berhasil diambil", response)
}

// CreateOrder menangani request untuk membuat order baru
func (h *OrderAdaptor) CreateOrder(c *gin.Context) {
	h.logger.Debug("CreateOrder handler called", zap.String("client_ip", c.ClientIP()))

	var req dto.OrderCreateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create order",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	h.logger.Debug("Request body bound successfully",
		zap.String("customer_name", req.CustomerName),
		zap.Uint("table_id", req.TableID),
		zap.Int("items_count", len(req.Items)),
	)

	// Call usecase
	response, err := h.orderUsecase.CreateOrder(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create order",
			zap.Error(err),
			zap.String("customer_name", req.CustomerName),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal membuat order: "+err.Error())
		return
	}

	h.logger.Info("Order created successfully",
		zap.Uint("id", response.ID),
		zap.String("customer_name", response.CustomerName),
		zap.Float64("total_amount", response.TotalAmount),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusCreated, "Order berhasil dibuat", response)
}

// UpdateOrder menangani request untuk update order yang sudah ada
func (h *OrderAdaptor) UpdateOrder(c *gin.Context) {
	h.logger.Debug("UpdateOrder handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	var req dto.OrderUpdateRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update order",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	h.logger.Debug("Request body bound successfully",
		zap.Uint("id", uint(id)),
		zap.String("customer_name", req.CustomerName),
		zap.Int("items_count", len(req.Items)),
	)

	// Call usecase
	err = h.orderUsecase.UpdateOrder(c.Request.Context(), uint(id), req)
	if err != nil {
		h.logger.Error("Failed to update order",
			zap.Error(err),
			zap.Uint("id", uint(id)),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal memperbarui order: "+err.Error())
		return
	}

	h.logger.Info("Order updated successfully",
		zap.Uint("id", uint(id)),
		zap.String("customer_name", req.CustomerName),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Order berhasil diperbarui", nil)
}

// DeleteOrder menangani request untuk menghapus order
func (h *OrderAdaptor) DeleteOrder(c *gin.Context) {
	h.logger.Debug("DeleteOrder handler called", zap.String("client_ip", c.ClientIP()))

	// Get ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			zap.Error(err),
			zap.String("id", idStr),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid ID parameter")
		return
	}

	h.logger.Debug("ID parameter bound successfully", zap.Uint("id", uint(id)))

	// Call usecase
	err = h.orderUsecase.DeleteOrder(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete order",
			zap.Error(err),
			zap.Uint("id", uint(id)),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal menghapus order: "+err.Error())
		return
	}

	h.logger.Info("Order deleted successfully", zap.Uint("id", uint(id)))

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Order berhasil dihapus", nil)
}

// GetAllTables menangani request untuk mengambil semua meja
func (h *OrderAdaptor) GetAllTables(c *gin.Context) {
	h.logger.Debug("GetAllTables handler called")

	// Call usecase
	response, err := h.orderUsecase.GetAllTables(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get all tables",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data meja: "+err.Error())
		return
	}

	h.logger.Info("GetAllTables completed successfully",
		zap.Int("total_tables", len(response)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data meja berhasil diambil", response)
}

// GetAllPaymentMethods menangani request untuk mengambil semua metode pembayaran
func (h *OrderAdaptor) GetAllPaymentMethods(c *gin.Context) {
	h.logger.Debug("GetAllPaymentMethods handler called")

	// Call usecase
	response, err := h.orderUsecase.GetAllPaymentMethods(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get all payment methods",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data metode pembayaran: "+err.Error())
		return
	}

	h.logger.Info("GetAllPaymentMethods completed successfully",
		zap.Int("total_payment_methods", len(response)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data metode pembayaran berhasil diambil", response)
}

// GetAvailableChairs menangani request untuk mengambil kursi yang tersedia
func (h *OrderAdaptor) GetAvailableChairs(c *gin.Context) {
	h.logger.Debug("GetAvailableChairs handler called")

	// Call usecase
	response, err := h.orderUsecase.GetAvailableChairs(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get available chairs",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil data kursi tersedia: "+err.Error())
		return
	}

	h.logger.Info("GetAvailableChairs completed successfully",
		zap.Int("total_available_chairs", len(response)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Data kursi tersedia berhasil diambil", response)
}
