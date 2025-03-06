package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, _ := database.DB()

	// ✅ Optimized Connection Pooling
	sqlDB.SetMaxOpenConns(1000)               // Maximum concurrent connections
	sqlDB.SetMaxIdleConns(100)                // Keep idle connections for reuse
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // No timeout for persistent connections

	// ✅ Ensure email column is indexed
	database.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)")

	// ✅ AutoMigrate with new fields
	database.AutoMigrate(&models.User{})

	DB = database
	fmt.Println("✅ Database connected successfully!")
}
