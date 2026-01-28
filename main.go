package main

import (
	"flag"
	"fmt"
	"log"

	"aplikasi-pos-team-boolean/internal/wire"
	"aplikasi-pos-team-boolean/pkg/database"
	"aplikasi-pos-team-boolean/pkg/utils"

	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	migrate := flag.Bool("migrate", true, "Run database migration on startup")
	seed := flag.Bool("seed", false, "Seed database with initial data")
	reset := flag.Bool("reset", false, "Reset database (drop and recreate tables - DEVELOPMENT ONLY!)")
	flag.Parse()

	// Baca config dari .env
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Initialize logger
	logger, err := utils.InitLogger("logs", config.Debug)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Application starting",
		zap.String("environment", config.Env),
		zap.Bool("debug", config.Debug),
	)

	// Koneksi ke database
	db, err := database.InitDB(config.DB)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	logger.Info("Database connection successful")

	// Database Migration
	if *reset {
		// Reset database (HATI-HATI!)
		logger.Warn("RESET MODE: Dropping all tables and recreating...")
		if err := database.ResetDatabase(db, *seed); err != nil {
			logger.Fatal("Failed to reset database", zap.Error(err))
		}
	} else if *migrate {
		// Auto migrate
		logger.Info("Running database migration", zap.Bool("seed", *seed))
		if err := database.MigrateWithSeed(db, *seed); err != nil {
			logger.Fatal("Failed to run migration", zap.Error(err))
		}
	}

	// Initialize app dengan dependency injection
	router := wire.InitializeApp(db, logger)

	// Jalankan server
	port := config.Port
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", zap.String("port", port))
	fmt.Println("========================================")
	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	fmt.Println("========================================")

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
