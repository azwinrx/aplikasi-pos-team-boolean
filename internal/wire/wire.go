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
	// Setup Gin (without default middleware)
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add custom logging middleware with zap
	router.Use(middleware.LoggingMiddleware(logger))

	// Setup repositories
	repo := repository.NewRepository(db, logger)

	// Setup use cases with logger
	inventoriesUsecase := usecase.NewInventoriesUsecase(repo.InventoriesRepo, logger)
	staffUsecase := usecase.NewStaffUseCase(repo.StaffRepo, logger)

	// Setup handlers with logger
	inventoriesHandler := adaptor.NewInventoriesAdaptor(inventoriesUsecase, logger)
	staffHandler := adaptor.NewStaffAdaptor(staffUsecase, logger)

	// Setup routes
	setupRoutes(router, inventoriesHandler, staffHandler, logger)

	return router
}

// setupRoutes mengatur semua routing untuk aplikasi
func setupRoutes(router *gin.Engine, inventoriesHandler *adaptor.InventoriesAdaptor, staffHandler *adaptor.StaffAdaptor, logger *zap.Logger) {
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
	}

	logger.Info("Routes registered successfully")
}
