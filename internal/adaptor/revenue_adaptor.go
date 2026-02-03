package adaptor

import (
	"net/http"
	"strconv"

	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RevenueAdaptor menangani semua request HTTP untuk revenue reports
type RevenueAdaptor struct {
	revenueUsecase usecase.RevenueUseCase
	logger         *zap.Logger
}

// NewRevenueAdaptor membuat instance baru dari RevenueAdaptor
func NewRevenueAdaptor(revenueUsecase usecase.RevenueUseCase, logger *zap.Logger) *RevenueAdaptor {
	return &RevenueAdaptor{
		revenueUsecase: revenueUsecase,
		logger:         logger,
	}
}

// GetRevenueByStatus menangani request untuk mengambil total revenue dan breakdown berdasarkan status
func (h *RevenueAdaptor) GetRevenueByStatus(c *gin.Context) {
	h.logger.Debug("GetRevenueByStatus handler called")

	// Ambil status dari query param
	status := c.Query("status")
	if status == "" {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Parameter 'status' wajib diisi (pending/paid)")
		return
	}

	// Call usecase
	response, err := h.revenueUsecase.GetRevenueByStatus(c.Request.Context(), status)
	if err != nil {
		h.logger.Error("Failed to get revenue by status",
			zap.Error(err),
			zap.String("status", status),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil revenue berdasarkan status: "+err.Error())
		return
	}

	h.logger.Info("GetRevenueByStatus completed successfully",
		zap.String("status", status),
		zap.Float64("total_revenue", response.TotalRevenue),
		zap.Int("status_count", len(response.Breakdown)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Revenue berdasarkan status berhasil diambil", response)
}

// GetRevenuePerMonth menangani request untuk mengambil total revenue per bulan tertentu
func (h *RevenueAdaptor) GetRevenuePerMonth(c *gin.Context) {
	h.logger.Debug("GetRevenuePerMonth handler called")

	// Get year from query parameter (optional)
	yearStr := c.Query("year")
	year := 0
	if yearStr != "" {
		parsedYear, err := strconv.Atoi(yearStr)
		if err != nil {
			h.logger.Warn("Invalid year parameter",
				zap.Error(err),
				zap.String("year", yearStr),
				zap.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c.Writer, http.StatusBadRequest, "Invalid year parameter")
			return
		}
		year = parsedYear
	}

	// Get month from query parameter (required)
	monthStr := c.Query("month")
	if monthStr == "" {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Parameter 'month' wajib diisi (1-12)")
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Parameter 'month' harus berupa angka 1-12")
		return
	}

	h.logger.Debug("Year and month parameter processed", zap.Int("year", year), zap.Int("month", month))

	// Call usecase
	response, err := h.revenueUsecase.GetRevenuePerMonth(c.Request.Context(), year, month)
	if err != nil {
		h.logger.Error("Failed to get revenue for month",
			zap.Error(err),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil revenue per bulan: "+err.Error())
		return
	}

	h.logger.Info("GetRevenuePerMonth completed successfully",
		zap.Int("year", response.Year),
		zap.Int("month", month),
		zap.Float64("total_revenue", response.TotalRevenue),
		zap.Int("month_count", len(response.Monthly)),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Revenue per bulan berhasil diambil", response)
}

// GetProductRevenueList menangani request untuk mengambil detail revenue produk tertentu
func (h *RevenueAdaptor) GetProductRevenueList(c *gin.Context) {
	h.logger.Debug("GetProductRevenueList handler called")

	// Ambil productID dari query param
	productIDStr := c.Query("productID")
	if productIDStr == "" {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Parameter 'productID' wajib diisi")
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		utils.ResponseError(c.Writer, http.StatusBadRequest, "Parameter 'productID' harus berupa angka lebih dari 0")
		return
	}

	// Call usecase
	response, err := h.revenueUsecase.GetProductRevenueList(c.Request.Context(), productID)
	if err != nil {
		h.logger.Error("Failed to get product revenue",
			zap.Error(err),
			zap.Int("product_id", productID),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.ResponseError(c.Writer, http.StatusInternalServerError, "Gagal mengambil detail revenue produk: "+err.Error())
		return
	}

	h.logger.Info("GetProductRevenueList completed successfully",
		zap.Int("product_id", productID),
		zap.Int("total_products", response.TotalProducts),
	)

	// Return response
	utils.ResponseSuccess(c.Writer, http.StatusOK, "Detail revenue produk berhasil diambil", response)
}
