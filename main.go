package main

import (
	"fmt"
	"log"
	"aplikasi-pos-team-boolean/pkg/database"
	"aplikasi-pos-team-boolean/pkg/utils"
)

func main() {
	// Baca config dari .env
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Cek koneksi database (ignore return DB, hanya handle error)
	_, err = database.InitDB(config.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful!")
}
