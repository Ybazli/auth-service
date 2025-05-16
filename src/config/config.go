package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ybazli/auth-service/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = db
	fmt.Println("Database connected")
	DB.AutoMigrate(&models.User{})
}
