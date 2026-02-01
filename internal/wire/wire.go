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
	setupRoutes(router, adaptorInstance.AuthAdaptor, adaptorInstance.InventoriesAdaptor, adaptorInstance.StaffAdaptor, adaptorInstance.OrderAdaptor, adaptorInstance.CategoryAdaptor, adaptorInstance.ProductAdaptor, adaptorInstance.DashboardAdaptor, uc.DashboardUseCase, logger)

	return router
}

// setupRoutes mengatur semua routing untuk aplikasi
func setupRoutes(router *gin.Engine, authHandler *adaptor.AuthAdaptor, inventoriesHandler *adaptor.InventoriesAdaptor, staffHandler *adaptor.StaffAdaptor, orderHandler *adaptor.OrderAdaptor, categoryHandler *adaptor.CategoryAdaptor, productHandler *adaptor.ProductAdaptor, dashboardHandler adaptor.DashboardHandler, dashboardUC usecase.DashboardUseCase, logger *zap.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.ResponseSuccess(c.Writer, 200, "Server is running", map[string]string{
			"status": "healthy",
		})
	})

	// Inisialisasi websocket dashboard handler
	dashboardWsHandler := adaptor.NewDashboardWebsocketHandler(dashboardUC, logger)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			// 1. POST Login
			auth.POST("/login", authHandler.Login)

			// 2. POST Check Email
			auth.POST("/check-email", authHandler.CheckEmail)

			// 3. POST Send OTP
			auth.POST("/send-otp", authHandler.SendOTP)

			// 4. POST Validate OTP
			auth.POST("/validate-otp", authHandler.ValidateOTP)

			// 5. POST Reset Password
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

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

		// Category routes (Menu)
		categories := v1.Group("/categories")
		{
			// 1. GET all categories
			categories.GET("", categoryHandler.GetList)

			// 2. POST Create category
			categories.POST("", categoryHandler.Create)

			// 3. PUT Update category
			categories.PUT("/:id", categoryHandler.Update)

			// 4. GET category by ID
			categories.GET("/:id", categoryHandler.GetByID)

			// 5. DELETE category
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// Product routes (Menu)
		products := v1.Group("/products")
		{
			// 1. GET all products (with filter by category_id, is_available, price range)
			products.GET("", productHandler.GetList)

			// 2. GET products by category
			products.GET("/category/:category_id", productHandler.GetByCategory)

			// 3. POST Create product
			products.POST("", productHandler.Create)

			// 4. PUT Update product
			products.PUT("/:id", productHandler.Update)

			// 5. GET product by ID
			products.GET("/:id", productHandler.GetByID)

			// 6. DELETE product
			products.DELETE("/:id", productHandler.Delete)
		}

		// Dashboard routes
		dashboard := v1.Group("/dashboard")
		{
			// 1. GET dashboard summary (daily sales, monthly sales, table summary)
			dashboard.GET("/summary", dashboardHandler.GetSummary)

			// 2. GET popular products
			dashboard.GET("/popular-products", dashboardHandler.GetPopularProducts)

			// 3. GET new products (< 30 days)
			dashboard.GET("/new-products", dashboardHandler.GetNewProducts)

			// 4. Export dashboard data (bulan, jumlah order, sales, revenue)
			dashboard.GET("/export", dashboardHandler.ExportDashboard)

			// 5. Websocket realtime dashboard (revenue & sales)
			dashboard.GET("/ws", dashboardWsHandler.ServeWs)
		}
	}

	logger.Info("Routes registered successfully")
}
