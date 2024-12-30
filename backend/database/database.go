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

	db.Exec(`
		DO $$ BEGIN
			CREATE TYPE score_lamp AS ENUM ('FAILED', 'CLEAR', 'FULL COMBO', 'ALL JUSTICE', 'ALL JUSTICE CRITICAL');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		`,
	)

	if err := db.AutoMigrate(&models.User{}, &models.Score{}, &models.Song{}, &models.Chart{}, &models.UserAPIKey{}, &SeedVersion{}); err != nil {
		logger.Error("database.init", err, "failed to migrate database")
		panic("failed to migrate database")
	}

	if err := SeedDatabase(db); err != nil {
		logger.Error("database.seed", err, "failed to seed database")
		panic("failed to seed database")
	}

	DB = db
}
