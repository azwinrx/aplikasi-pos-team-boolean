package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aplikasi-pos-team-boolean/internal/wire"
	"aplikasi-pos-team-boolean/pkg/database"
	"aplikasi-pos-team-boolean/pkg/utils"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	migrate := pflag.Bool("migrate", true, "Run database migration on startup")
	seed := pflag.Bool("seed", true, "Seed database with initial data")
	reset := pflag.Bool("reset", false, "Reset database (drop and recreate tables - DEVELOPMENT ONLY!)")
	pflag.Parse()

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

	// Setup HTTP Server
	port := config.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Channel untuk signal interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan server di goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", port))
		fmt.Println("========================================")
		fmt.Printf("Server running on http://localhost:%s\n", port)
		fmt.Println("Press Ctrl+C to shutdown gracefully")
		fmt.Println("========================================")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Tunggu signal interrupt (Ctrl+C)
	<-quit
	logger.Info("Shutting down server...")
	fmt.Println("\nGraceful shutdown initiated...")

	// Timeout 5 detik untuk graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server dengan graceful
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
	fmt.Println("Server stopped gracefully")
}
