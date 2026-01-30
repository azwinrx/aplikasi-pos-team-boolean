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
	setupRoutes(router, adaptorInstance.AuthAdaptor, adaptorInstance.AdminAdaptor, adaptorInstance.InventoriesAdaptor, adaptorInstance.StaffAdaptor, adaptorInstance.OrderAdaptor, adaptorInstance.NotificationAdaptor, logger)

	return router
}

// setupRoutes mengatur semua routing untuk aplikasi
func setupRoutes(router *gin.Engine, authHandler *adaptor.AuthAdaptor, adminHandler *adaptor.AdminAdaptor, inventoriesHandler *adaptor.InventoriesAdaptor, staffHandler *adaptor.StaffAdaptor, orderHandler *adaptor.OrderAdaptor, notificationHandler *adaptor.NotificationAdaptor, logger *zap.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.ResponseSuccess(c.Writer, 200, "Server is running", map[string]string{
			"status": "healthy",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
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

			// 6. POST Logout (protected)
			auth.POST("/logout", middleware.AuthMiddleware(logger), adminHandler.Logout)
		}

		// User Profile routes (protected)
		profile := v1.Group("/profile")
		profile.Use(middleware.AuthMiddleware(logger))
		{
			// 1. GET user profile
			profile.GET("", adminHandler.GetUserProfile)

			// 2. PUT update user profile
			profile.PUT("", adminHandler.UpdateUserProfile)
		}

		// Admin routes (protected, superadmin only)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(logger))
		{
			// 1. GET list admins (superadmin only)
			admin.GET("/list", adminHandler.ListAdmins)

			// 2. PUT edit admin access (superadmin only)
			admin.PUT("/:id/access", adminHandler.EditAdminAccess)

			// 3. POST create admin (superadmin only)
			admin.POST("/create", adminHandler.CreateAdmin)
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

		// Notification routes
		notifications := v1.Group("/notifications")
		{
			// 1. GET all notifications with filters and pagination
			notifications.GET("", notificationHandler.ListNotifications)

			// 2. PUT Update notification status
			notifications.PUT("/:id/status", notificationHandler.UpdateNotificationStatus)

			// 3. DELETE notification
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		}
	}

	logger.Info("Routes registered successfully")
}

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

		// Notification routes
		notifications := v1.Group("/notifications")
		{
			// 1. GET all notifications with filters and pagination
			notifications.GET("", notificationHandler.ListNotifications)

			// 2. PUT Update notification status
			notifications.PUT("/:id/status", notificationHandler.UpdateNotificationStatus)

			// 3. DELETE notification
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		}
	}

	logger.Info("Routes registered successfully")
}
