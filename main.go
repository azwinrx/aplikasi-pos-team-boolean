package main

import (
	"fmt"
	"log"

	"aplikasi-pos-team-boolean/internal/wire"
	"aplikasi-pos-team-boolean/pkg/database"
	"aplikasi-pos-team-boolean/pkg/utils"
)

func main() {
	// Baca config dari .env
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Koneksi ke database
	db, err := database.InitDB(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful!")

	// Initialize app dengan dependency injection
	router := wire.InitializeApp(db)

	// Jalankan server
	port := config.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
