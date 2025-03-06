package database

import (
	"log"
	"os"
	"time"

	"github.com/Aditya-PS-05/NeetChamp/user-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, _ := database.DB()
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	database.AutoMigrate(&models.Mentee{}, &models.Admin{})

	DB = database
	log.Println("✅ Database connected successfully!")
}
