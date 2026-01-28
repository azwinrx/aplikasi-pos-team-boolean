package main

import (
	"aplikasi-pos-team-boolean/cmd"
	"aplikasi-pos-team-boolean/internal/data/repository"
	"aplikasi-pos-team-boolean/internal/wire"
	"aplikasi-pos-team-boolean/pkg/database"
	"aplikasi-pos-team-boolean/pkg/utils"
	"log"
)

func main() {
	// Baca config dari .env
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := repository.NewRepository(db)
	router := wire.Wiring(repo)
	cmd.APiserver(router)
}
