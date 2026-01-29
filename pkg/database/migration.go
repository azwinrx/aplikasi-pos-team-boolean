package database

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate melakukan auto migration untuk semua entity/table
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto migration...")

	// Daftar semua entity yang akan dimigrate
	entities := []interface{}{
		&entity.Staff{},
		&entity.Inventories{},
		&entity.Table{},
	}

	// Jalankan auto migration
	if err := db.AutoMigrate(entities...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database auto migration completed successfully!")

	// Log semua tabel yang berhasil dimigrate
	for _, e := range entities {
		tableName := db.NamingStrategy.TableName(fmt.Sprintf("%T", e))
		log.Printf("   Table migrated: %s", tableName)
	}

	return nil
}

// MigrateWithSeed melakukan migration dan seeding data (jika diperlukan)
func MigrateWithSeed(db *gorm.DB, withSeed bool) error {
	// Auto migrate
	if err := AutoMigrate(db); err != nil {
		return err
	}

	// Seeding (optional)
	if withSeed {
		log.Println("Starting database seeding...")
		if err := SeedData(db); err != nil {
			return fmt.Errorf("failed to seed data: %w", err)
		}
		log.Println("Database seeding completed!")
	}

	return nil
}

// SeedData memasukkan data awal ke database (optional)
func SeedData(db *gorm.DB) error {
	// Contoh: Seed data inventories jika table masih kosong
	var count int64
	db.Model(&entity.Inventories{}).Count(&count)

	if count == 0 {
		log.Println("   Seeding inventories data...")
		seedInventories := []entity.Inventories{
			{
				Name:        "Coca Cola 1L",
				Category:    "beverage",
				Quantity:    150,
				Unit:        "litre",
				MinStock:    50,
				RetailPrice: 15.50,
				Status:      "active",
			},
			{
				Name:        "Sprite 1L",
				Category:    "beverage",
				Quantity:    120,
				Unit:        "litre",
				MinStock:    50,
				RetailPrice: 14.00,
				Status:      "active",
			},
			{
				Name:        "Pepsi 1L",
				Category:    "beverage",
				Quantity:    45,
				Unit:        "litre",
				MinStock:    50,
				RetailPrice: 15.00,
				Status:      "active",
			},
			{
				Name:        "Mineral Water 1.5L",
				Category:    "beverage",
				Quantity:    200,
				Unit:        "litre",
				MinStock:    100,
				RetailPrice: 5.00,
				Status:      "active",
			},
			{
				Name:        "Orange Juice 1L",
				Category:    "beverage",
				Quantity:    30,
				Unit:        "litre",
				MinStock:    40,
				RetailPrice: 25.00,
				Status:      "active",
			},
		}

		if err := db.Create(&seedInventories).Error; err != nil {
			return fmt.Errorf("failed to seed inventories: %w", err)
		}
		log.Printf("   Seeded %d inventories", len(seedInventories))
	}

	// Seed staff jika masih kosong
	db.Model(&entity.Staff{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding staff data...")
		seedStaff := []entity.Staff{
			{
				FullName:    "John Doe",
				Email:       "john.doe@example.com",
				Role:        "manager",
				PhoneNumber: "081234567890",
				Salary:      5000000,
				Address:     "Jakarta, Indonesia",
			},
			{
				FullName:    "Jane Smith",
				Email:       "jane.smith@example.com",
				Role:        "cashier",
				PhoneNumber: "081234567891",
				Salary:      3500000,
				Address:     "Bandung, Indonesia",
			},
			{
				FullName:    "Bob Wilson",
				Email:       "bob.wilson@example.com",
				Role:        "staff",
				PhoneNumber: "081234567892",
				Salary:      4000000,
				Address:     "Surabaya, Indonesia",
			},
		}

		if err := db.Create(&seedStaff).Error; err != nil {
			return fmt.Errorf("failed to seed staff: %w", err)
		}
		log.Printf("   Seeded %d staff", len(seedStaff))
	}

	// Seed tables jika masih kosong
	db.Model(&entity.Table{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding tables data...")
		seedTables := []entity.Table{
			{
				Number:   "T01",
				Capacity: 4,
				Status:   "available",
			},
			{
				Number:   "T02",
				Capacity: 2,
				Status:   "available",
			},
			{
				Number:   "T03",
				Capacity: 6,
				Status:   "available",
			},
			{
				Number:   "T04",
				Capacity: 4,
				Status:   "available",
			},
			{
				Number:   "T05",
				Capacity: 2,
				Status:   "available",
			},
		}

		if err := db.Create(&seedTables).Error; err != nil {
			return fmt.Errorf("failed to seed tables: %w", err)
		}
		log.Printf("   Seeded %d tables", len(seedTables))
	}

	return nil
}

// DropAllTables menghapus semua table (HATI-HATI! Untuk development only)
func DropAllTables(db *gorm.DB) error {
	log.Println("WARNING: Dropping all tables...")

	entities := []interface{}{
		&entity.Staff{},
		&entity.Inventories{},
		&entity.Table{},
	}

	for _, e := range entities {
		if err := db.Migrator().DropTable(e); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("All tables dropped successfully!")
	return nil
}

// ResetDatabase drop semua table dan migrate ulang (DEVELOPMENT ONLY!)
func ResetDatabase(db *gorm.DB, withSeed bool) error {
	log.Println("Resetting database...")

	// Drop all tables
	if err := DropAllTables(db); err != nil {
		return err
	}

	// Migrate ulang
	if err := MigrateWithSeed(db, withSeed); err != nil {
		return err
	}

	log.Println("Database reset completed!")
	return nil
}
