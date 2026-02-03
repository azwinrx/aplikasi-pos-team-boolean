package adaptor

import (
	"net/http"
	"strconv"

	"aplikasi-pos-team-boolean/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardHandler interface {
	GetSummary(c *gin.Context)
	GetPopularProducts(c *gin.Context)
	GetNewProducts(c *gin.Context)
	ExportDashboard(c *gin.Context)
}

type dashboardHandler struct {
	dashboardUC usecase.DashboardUseCase
	logger      *zap.Logger
}

func NewDashboardHandler(dashboardUC usecase.DashboardUseCase, logger *zap.Logger) DashboardHandler {
	return &dashboardHandler{
		dashboardUC: dashboardUC,
		logger:      logger,
	}
}

// GetSummary godoc
// @Summary Get Dashboard Summary
// @Description Get dashboard summary including daily sales, monthly sales, and table summary
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.DashboardSummaryResponse
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /dashboard/summary [get]
func (h *dashboardHandler) GetSummary(c *gin.Context) {
	h.logger.Info("GetSummary request received")

	summary, err := h.dashboardUC.GetSummary(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get dashboard summary", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get dashboard summary",
			"data":    nil,
		})
		return
	}

	h.logger.Info("Dashboard summary retrieved successfully")
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Dashboard summary retrieved successfully",
		"data":    summary,
	})
}

// GetPopularProducts godoc
// @Summary Get Popular Products
// @Description Get list of popular products based on total sales
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of results (default: 10)"
// @Success 200 {object} []dto.PopularProductResponse
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /dashboard/popular-products [get]
func (h *dashboardHandler) GetPopularProducts(c *gin.Context) {
	h.logger.Info("GetPopularProducts request received")

	// Get limit from query param
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	products, err := h.dashboardUC.GetPopularProducts(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get popular products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get popular products",
			"data":    nil,
		})
		return
	}

	h.logger.Info("Popular products retrieved successfully", zap.Int("count", len(products)))
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Popular products retrieved successfully",
		"data":    products,
	})
}

// GetNewProducts godoc
// @Summary Get New Products
// @Description Get list of new products created in the last 30 days
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of results (default: 10)"
// @Success 200 {object} []dto.NewProductResponse
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /dashboard/new-products [get]
func (h *dashboardHandler) GetNewProducts(c *gin.Context) {
	h.logger.Info("GetNewProducts request received")

	// Get limit from query param
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	products, err := h.dashboardUC.GetNewProducts(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get new products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get new products",
			"data":    nil,
		})
		return
	}

	h.logger.Info("New products retrieved successfully", zap.Int("count", len(products)))
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "New products retrieved successfully",
		"data":    products,
	})
}

// ExportDashboard godoc
// @Summary Export Dashboard Data
// @Description Export monthly dashboard data (bulan, jumlah order, sales, revenue)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []dto.DashboardExportRow
// @Router /dashboard/export [get]
func (h *dashboardHandler) ExportDashboard(c *gin.Context) {
	h.logger.Info("ExportDashboard request received")

	rows, err := h.dashboardUC.ExportDashboard(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to export dashboard data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to export dashboard data",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Export dashboard data success",
		"data":    rows,
	})
}
