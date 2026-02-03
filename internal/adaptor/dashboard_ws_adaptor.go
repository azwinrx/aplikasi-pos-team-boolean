package adaptor

import (
	"net/http"
	"time"

	"aplikasi-pos-team-boolean/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type DashboardWebsocketHandler struct {
	dashboardUC usecase.DashboardUseCase
	logger      *zap.Logger
}

func NewDashboardWebsocketHandler(dashboardUC usecase.DashboardUseCase, logger *zap.Logger) *DashboardWebsocketHandler {
	return &DashboardWebsocketHandler{
		dashboardUC: dashboardUC,
		logger:      logger,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Websocket endpoint: /api/v1/dashboard/ws
func (h *DashboardWebsocketHandler) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade websocket", zap.Error(err))
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Ambil summary terbaru
			summary, err := h.dashboardUC.GetSummary(c.Request.Context())
			if err != nil {
				h.logger.Error("Failed to get dashboard summary", zap.Error(err))
				continue
			}
			// Kirim summary (hanya revenue & sales)
			msg := map[string]interface{}{
				"daily_sales":    summary.DailySales.TotalRevenue,
				"monthly_sales":  summary.MonthlySales.TotalRevenue,
				"daily_orders":   summary.DailySales.TotalOrders,
				"monthly_orders": summary.MonthlySales.TotalOrders,
			}
			if err := conn.WriteJSON(msg); err != nil {
				h.logger.Error("Failed to write websocket message", zap.Error(err))
				return
			}
		}
	}
}
