package database

import (
	"fmt"
	"os"

	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("database.init", err, "failed to connect to database")
		panic("failed to connect to database")
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Error("database.init", err, "failed to migrate database")
		panic("failed to migrate database")
	}

	DB = db
}
