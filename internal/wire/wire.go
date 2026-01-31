package wire

import (
	"aplikasi-pos-team-boolean/internal/adaptor"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/middleware"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitializeApp membuat dan mengkonfigurasi aplikasi dengan semua dependencies
func InitializeApp(db *gorm.DB, logger *zap.Logger) *gin.Engine {
	// Setup Gin with default middleware
	router := gin.Default()

	// Add custom logging middleware with zap
	router.Use(middleware.LoggingMiddleware(logger))

	// Setup repositories
	repo := repository.NewRepository(db, logger)

	// Setup use cases with UseCase struct (embedding)
	uc := usecase.NewUseCase(&repo, logger, db)

	// Setup adaptor
	adaptorInstance := adaptor.NewAdaptor(uc, logger)

	// Setup routes
	setupRoutes(router, adaptorInstance.InventoriesAdaptor, adaptorInstance.StaffAdaptor, adaptorInstance.OrderAdaptor, adaptorInstance.RevenueAdaptor, adaptorInstance.ReservationsAdaptor, logger)

	return router
}

// setupRoutes mengatur semua routing untuk aplikasi
func setupRoutes(router *gin.Engine, inventoriesHandler *adaptor.InventoriesAdaptor, staffHandler *adaptor.StaffAdaptor, orderHandler *adaptor.OrderAdaptor, revenueHandler *adaptor.RevenueAdaptor, reservationHandler *adaptor.ReservationsAdaptor, logger *zap.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.ResponseSuccess(c.Writer, 200, "Server is running", map[string]string{
			"status": "healthy",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Inventories routes
		inventories := v1.Group("/inventories")
		{
			// 1. Get all inventories (tanpa filter)
			inventories.GET("", inventoriesHandler.GetAllInventories)

			// 2. Get inventories dengan filter (semua filter masuk di query params)
			// Filter: status, category, stock, unit, min_qty, max_qty, min_price, max_price
			inventories.GET("/filter", inventoriesHandler.GetInventoryByFilter)

			// 3. Create inventory
			inventories.POST("", inventoriesHandler.CreateInventory)

			// 4. Update inventory
			inventories.PUT("/:id", inventoriesHandler.UpdateInventory)

			// 5. Delete inventory
			inventories.DELETE("/:id", inventoriesHandler.DeleteInventory)
		}

		// Staff routes
		staff := v1.Group("/staff")
		{
			// 1. GET all staff (pagination only)
			staff.GET("", staffHandler.GetList)

			// 2. POST Create staff
			staff.POST("", staffHandler.Create)

			// 3. PUT Update staff
			staff.PUT("/:id", staffHandler.Update)

			// 4. GET staff by email (query param: ?email=xxx) - MUST BE BEFORE /:id
			staff.GET("/email", staffHandler.GetByEmail)

			// 5. GET staff by ID
			staff.GET("/:id", staffHandler.GetByID)

			// 6. DELETE staff
			staff.DELETE("/:id", staffHandler.Delete)
		}

		// Order routes
		order := v1.Group("/orders")
		{
			// 1. GET all orders
			order.GET("", orderHandler.GetAllOrders)

			// 2. POST Create order
			order.POST("", orderHandler.CreateOrder)

			// 3. PUT Update order
			order.PUT("/:id", orderHandler.UpdateOrder)

			// 4. DELETE order
			order.DELETE("/:id", orderHandler.DeleteOrder)

			// 5. GET all tables
			order.GET("/tables", orderHandler.GetAllTables)

			// 6. GET all payment methods
			order.GET("/payment-methods", orderHandler.GetAllPaymentMethods)

			// 7. GET available chairs
			order.GET("/available-chairs", orderHandler.GetAvailableChairs)
		}

		// Revenue Report routes
		revenue := v1.Group("/revenue")
		{
			// 1. GET total revenue dan breakdown berdasarkan status
			revenue.GET("/by-status", revenueHandler.GetRevenueByStatus)

			// 2. GET total revenue per bulan (optional query param: year)
			revenue.GET("/per-month", revenueHandler.GetRevenuePerMonth)

			// 3. GET list produk beserta detail revenue
			revenue.GET("/products", revenueHandler.GetProductRevenueList)
		}

		// Reservation routes
		reservation := v1.Group("/reservations")
		{
			// 1. GET all reservations
			reservation.GET("", reservationHandler.GetAllReservations)

			// 2. POST Create reservation
			reservation.POST("", reservationHandler.CreateReservation)

			// 3. PUT Update reservation
			reservation.PUT("/:id", reservationHandler.UpdateReservation)

			// 4. DELETE reservation
			reservation.DELETE("/:id", reservationHandler.DeleteReservation)
		}
	}

	logger.Info("Routes registered successfully")
}
