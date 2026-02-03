package database

import (
	"aplikasi-pos-team-boolean/internal/data/entity"
	"aplikasi-pos-team-boolean/pkg/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Helper function untuk hash password
func hashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

// AutoMigrate melakukan auto migration untuk semua entity/table
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto migration...")

	entities := []interface{}{
		&entity.User{},
		&entity.OTP{},
		&entity.Staff{},
		&entity.Inventories{},
		&entity.Table{},
		&entity.PaymentMethod{},
		&entity.Reservations{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.Category{},
		&entity.Product{},
		// Tambahkan entity lain jika ada
	}

	if err := db.AutoMigrate(entities...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database auto migration completed successfully!")
	return nil
}

// MigrateWithSeed melakukan migration dan seeding data (jika diperlukan)
func MigrateWithSeed(db *gorm.DB, withSeed bool) error {
	if err := AutoMigrate(db); err != nil {
		return err
	}
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

		hashedPassword, err := hashPassword("customer123")
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		seedUsers := []entity.User{
			{
				Email:     "john.customer@pos.com",
				Password:  hashedPassword,
				Name:      "John Customer",
				Role:      "customer",
				Status:    "active",
				IsDeleted: false,
			},
			{
				Email:     "jane.customer@pos.com",
				Password:  hashedPassword,
				Name:      "Jane Customer",
				Role:      "customer",
				Status:    "active",
				IsDeleted: false,
			},
			{
				Email:     "bob.customer@pos.com",
				Password:  hashedPassword,
				Name:      "Bob Customer",
				Role:      "customer",
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

	// Seed Payment Methods
	db.Model(&entity.PaymentMethod{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding payment_methods data...")
		paymentMethods := []entity.PaymentMethod{
			{Name: "Cash"},
			{Name: "QRIS"},
			{Name: "Debit"},
		}
		if err := db.Create(&paymentMethods).Error; err != nil {
			return fmt.Errorf("failed to seed payment_methods: %w", err)
		}
	}

	// Seed Staff
	db.Model(&entity.Staff{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding staff data...")
		staff := []entity.Staff{
			{FullName: "John Doe", Email: "john.doe@example.com", Role: "manager", PhoneNumber: "081234567890", Salary: 5000000, Address: "Jakarta, Indonesia"},
			{FullName: "Jane Smith", Email: "jane.smith@example.com", Role: "cashier", PhoneNumber: "081234567891", Salary: 3500000, Address: "Bandung, Indonesia"},
			{FullName: "Bob Wilson", Email: "bob.wilson@example.com", Role: "staff", PhoneNumber: "081234567892", Salary: 4000000, Address: "Surabaya, Indonesia"},
		}
		if err := db.Create(&staff).Error; err != nil {
			return fmt.Errorf("failed to seed staff: %w", err)
		}
	}

	// Seed Tables
	db.Model(&entity.Table{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding tables data...")
		tables := []entity.Table{
			{Number: "T01", Capacity: 4, Status: "available"},
			{Number: "T02", Capacity: 2, Status: "available"},
			{Number: "T03", Capacity: 6, Status: "available"},
			{Number: "T04", Capacity: 4, Status: "available"},
			{Number: "T05", Capacity: 2, Status: "available"},
		}
		if err := db.Create(&tables).Error; err != nil {
			return fmt.Errorf("failed to seed tables: %w", err)
		}
	}

	// Seed Inventories
	db.Model(&entity.Inventories{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding inventories data...")
		inventories := []entity.Inventories{
			{Name: "Coca Cola 1L", Category: "beverage", Quantity: 150, Unit: "litre", MinStock: 50, RetailPrice: 15.50, Status: "active"},
			{Name: "Sprite 1L", Category: "beverage", Quantity: 120, Unit: "litre", MinStock: 50, RetailPrice: 14.00, Status: "active"},
			{Name: "Pepsi 1L", Category: "beverage", Quantity: 45, Unit: "litre", MinStock: 50, RetailPrice: 15.00, Status: "active"},
			{Name: "Mineral Water 1.5L", Category: "beverage", Quantity: 200, Unit: "litre", MinStock: 100, RetailPrice: 5.00, Status: "active"},
			{Name: "Orange Juice 1L", Category: "beverage", Quantity: 30, Unit: "litre", MinStock: 40, RetailPrice: 25.00, Status: "active"},
		}
		if err := db.Create(&inventories).Error; err != nil {
			return fmt.Errorf("failed to seed inventories: %w", err)
		}
	}

	// Seed Reservations
	db.Model(&entity.Reservations{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding reservations data...")
		now := time.Now()
		reservations := []entity.Reservations{
			{CustomerName: "Jane Smith", CustomerPhone: "081234567891", TableID: 1, ReservationTime: &now, Status: "pending"},
			{CustomerName: "Bob Wilson", CustomerPhone: "081234567892", TableID: 2, ReservationTime: &now, Status: "confirmed"},
		}
		if err := db.Create(&reservations).Error; err != nil {
			return fmt.Errorf("failed to seed reservations: %w", err)
		}
	}

	// Seed PaymentMethod, Order, OrderItem jika ingin contoh order
	db.Model(&entity.Order{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding orders data...")
		orders := []entity.Order{
			{
				UserID:          1,
				TableID:         1,
				PaymentMethodID: 1,
				CustomerName:    "John Doe",
				TotalAmount:     85500,
				Tax:             5500,
				Status:          "pending",
			},
			{
				UserID:          2,
				TableID:         2,
				PaymentMethodID: 2,
				CustomerName:    "Jane Smith",
				TotalAmount:     30000,
				Tax:             3000,
				Status:          "paid",
			},
		}
		if err := db.Create(&orders).Error; err != nil {
			return fmt.Errorf("failed to seed orders: %w", err)
		}
	}

	db.Model(&entity.OrderItem{}).Count(&count)
	if count == 0 {
		log.Println("   Seeding order_items data...")
		orderItems := []entity.OrderItem{
			{
				OrderID:   1,
				ProductID: 1,
				Quantity:  2,
				Price:     25000,
				Subtotal:  50000,
			},
			{
				OrderID:   2,
				ProductID: 2,
				Quantity:  1,
				Price:     30000,
				Subtotal:  30000,
			},
		}
		if err := db.Create(&orderItems).Error; err != nil {
			return fmt.Errorf("failed to seed order_items: %w", err)
		}
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
		&entity.Reservations{},
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
