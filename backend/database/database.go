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
	`)

	// dropping this column at every migration because the ALTER TABLE query dies if it's ALTER COLUMN and not ADD COLUMN
	db.Exec(`
		ALTER TABLE
			charts
		DROP COLUMN
			max_over_power;
	`)

	logger.Operation("db.init", "starting db migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Score{},
		&models.Chart{},
		&models.Song{},
		&models.UserAPIKey{},
		&models.TotalOverPower{},
	); err != nil {
		logger.Error("database.init", err, "failed to migrate database")
		panic("failed to migrate database")
	}
	logger.Operation("db.init", "finished db migrations!")

	logger.Operation("db.seeds", "seeding database...")
	if err := SeedDatabase(db); err != nil {
		logger.Error("db.seeds", err, "failed to seed database")
		panic("failed to seed database")
	}
	logger.Operation("db.seeds", "finished seeding database!")

	DB = db
}
