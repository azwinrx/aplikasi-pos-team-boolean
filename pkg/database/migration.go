package database

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/pkg/utils"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Helper function untuk hash password
func hashPassword(password string) string {
	return utils.HashPassword(password)
}

// AutoMigrate melakukan auto migration untuk semua entity/table
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto migration...")

	// Daftar semua entity yang akan dimigrate
	entities := []interface{}{
		&entity.User{},
		&entity.OTP{},
		&entity.Staff{},
		&entity.Inventories{},
		&entity.Table{},
		&entity.PaymentMethod{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Category{},
		&entity.Product{},
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
	// Import password hashing utility
	var count int64

	// Seed users jika masih kosong
	db.Model(&entity.User{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding users data...")
		seedUsers := []entity.User{
			{
				Email:     "admin@pos.com",
				Password:  hashPassword("admin123"),
				Name:      "Admin User",
				Role:      "admin",
				Status:    "active",
				IsDeleted: false,
			},
			{
				Email:     "manager@pos.com",
				Password:  hashPassword("manager123"),
				Name:      "Manager User",
				Role:      "manager",
				Status:    "active",
				IsDeleted: false,
			},
			{
				Email:     "staff@pos.com",
				Password:  hashPassword("staff123"),
				Name:      "Staff User",
				Role:      "staff",
				Status:    "active",
				IsDeleted: false,
			},
		}

		if err := db.Create(&seedUsers).Error; err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
		log.Printf("   Seeded %d users", len(seedUsers))
	}

	// Seed payment methods jika masih kosong
	db.Model(&entity.PaymentMethod{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding payment methods data...")
		seedPaymentMethods := []entity.PaymentMethod{
			{Name: "Cash"},
			{Name: "Credit Card"},
			{Name: "Debit Card"},
			{Name: "E-Wallet"},
			{Name: "Bank Transfer"},
		}

		if err := db.Create(&seedPaymentMethods).Error; err != nil {
			return fmt.Errorf("failed to seed payment methods: %w", err)
		}
		log.Printf("   Seeded %d payment methods", len(seedPaymentMethods))
	}

	// Contoh: Seed data inventories jika table masih kosong
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

	// Seed categories jika masih kosong
	db.Model(&entity.Category{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding categories data...")
		seedCategories := []entity.Category{
			{
				IconCategory: "🍕",
				CategoryName: "Pizza",
				Description:  "Delicious pizza varieties",
			},
			{
				IconCategory: "🍔",
				CategoryName: "Burger",
				Description:  "Juicy burgers and sandwiches",
			},
			{
				IconCategory: "🍗",
				CategoryName: "Chicken",
				Description:  "Crispy fried chicken",
			},
			{
				IconCategory: "🥐",
				CategoryName: "Bakery",
				Description:  "Fresh baked goods",
			},
			{
				IconCategory: "🥤",
				CategoryName: "Beverage",
				Description:  "Refreshing drinks",
			},
			{
				IconCategory: "🦐",
				CategoryName: "Seafood",
				Description:  "Fresh seafood dishes",
			},
		}

		if err := db.Create(&seedCategories).Error; err != nil {
			return fmt.Errorf("failed to seed categories: %w", err)
		}
		log.Printf("   Seeded %d categories", len(seedCategories))
	}

	// Seed products jika masih kosong
	db.Model(&entity.Product{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding products data...")

		// Get category IDs
		var pizzaCategory, burgerCategory, chickenCategory, beverageCategory entity.Category
		db.Where("category_name = ?", "Pizza").First(&pizzaCategory)
		db.Where("category_name = ?", "Burger").First(&burgerCategory)
		db.Where("category_name = ?", "Chicken").First(&chickenCategory)
		db.Where("category_name = ?", "Beverage").First(&beverageCategory)

		seedProducts := []entity.Product{
			{
				ProductImage: "/images/chicken-parmesan.jpg",
				ProductName:  "Chicken Parmesan",
				ItemID:       "#22314644",
				Stock:        119,
				CategoryID:   chickenCategory.ID,
				Price:        55.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/margherita-pizza.jpg",
				ProductName:  "Margherita Pizza",
				ItemID:       "#22314645",
				Stock:        85,
				CategoryID:   pizzaCategory.ID,
				Price:        45.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/pepperoni-pizza.jpg",
				ProductName:  "Pepperoni Pizza",
				ItemID:       "#22314646",
				Stock:        72,
				CategoryID:   pizzaCategory.ID,
				Price:        50.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/classic-burger.jpg",
				ProductName:  "Classic Burger",
				ItemID:       "#22314647",
				Stock:        95,
				CategoryID:   burgerCategory.ID,
				Price:        35.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/cheese-burger.jpg",
				ProductName:  "Cheese Burger",
				ItemID:       "#22314648",
				Stock:        8,
				CategoryID:   burgerCategory.ID,
				Price:        40.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/cola.jpg",
				ProductName:  "Cola",
				ItemID:       "#22314649",
				Stock:        200,
				CategoryID:   beverageCategory.ID,
				Price:        5.00,
				IsAvailable:  true,
			},
			{
				ProductImage: "/images/orange-juice.jpg",
				ProductName:  "Orange Juice",
				ItemID:       "#22314650",
				Stock:        0,
				CategoryID:   beverageCategory.ID,
				Price:        8.00,
				IsAvailable:  false,
			},
		}

		if err := db.Create(&seedProducts).Error; err != nil {
			return fmt.Errorf("failed to seed products: %w", err)
		}
		log.Printf("   Seeded %d products", len(seedProducts))
	}

	return nil
}

// DropAllTables menghapus semua table (HATI-HATI! Untuk development only)
func DropAllTables(db *gorm.DB) error {
	log.Println("WARNING: Dropping all tables...")

	entities := []interface{}{
		&entity.Product{},
		&entity.Category{},
		&entity.OrderItem{},
		&entity.Order{},
		&entity.PaymentMethod{},
		&entity.Table{},
		&entity.Inventories{},
		&entity.Staff{},
		&entity.OTP{},
		&entity.User{},
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
