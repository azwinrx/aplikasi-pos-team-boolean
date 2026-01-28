package wire

import (
	"log"

	"aplikasi-pos-team-boolean/internal/adaptor"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/usecase"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeApp membuat dan mengkonfigurasi aplikasi dengan semua dependencies
func InitializeApp(db *gorm.DB) *gin.Engine {
	// Setup Gin
	router := gin.Default()

	// Setup dependencies
	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo)
	inventoryHandler := adaptor.NewInventoryHandler(inventoryUsecase)

	// Setup routes
	setupRoutes(router, inventoryHandler)

	return router
}

// setupRoutes mengatur semua routing untuk aplikasi
func setupRoutes(router *gin.Engine, inventoryHandler *adaptor.InventoryHandler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.ResponseSuccess(c.Writer, 200, "Server is running", map[string]string{
			"status": "healthy",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Inventory routes
		inventory := v1.Group("/inventory")
		{
			// 1. List produk inventaris (dengan filter & pagination)
			inventory.GET("", inventoryHandler.GetAllInventories)

			// 2. Tambah data inventaris
			inventory.POST("", inventoryHandler.CreateInventory)

			// 3. Edit data inventaris
			inventory.PUT("/:id", inventoryHandler.UpdateInventory)

			// 4. Hapus data inventaris
			inventory.DELETE("/:id", inventoryHandler.DeleteInventory)

			// 5. Filter search by ID product
			inventory.GET("/:id", inventoryHandler.GetInventoryByID)
		}
	}

	log.Println("Routes registered successfully")
}
