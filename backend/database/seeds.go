package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/j1nxie/folern/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeedVersion struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChartSeed struct {
	ChartID    string          `json:"chartID"`
	Data       ChartSeedData   `json:"data"`
	Difficulty string          `json:"difficulty"`
	IsPrimary  bool            `json:"isPrimary"`
	Level      string          `json:"level"`
	LevelNum   decimal.Decimal `json:"levelNum"`
	Playtype   string          `json:"playtype"`
	SongID     int             `json:"songID"`
	Versions   []string        `json:"versions"`
}

type ChartSeedData struct {
	IngameID int `json:"inGameID"`
}

func (s ChartSeed) ToModel() models.Chart {
	return models.Chart{
		ID:     s.ChartID,
		SongID: s.SongID,
		Level:  s.LevelNum,
	}
}

type SongSeed struct {
	AltTitles   []string     `json:"altTitles"`
	Artist      string       `json:"artist"`
	Data        SongSeedData `json:"data"`
	ID          int          `json:"id"`
	SearchTerms []string     `json:"searchTerms"`
	Title       string       `json:"title"`
}

type SongSeedData struct {
	DisplayVersion string `json:"displayVersion"`
	Genre          string `json:"genre"`
}

func (s SongSeed) ToModel() models.Song {
	return models.Song{
		ID:      s.ID,
		Title:   s.Title,
		Artist:  s.Artist,
		Version: s.Data.DisplayVersion,
		Genre:   s.Data.Genre,
	}
}

func processSeedData[T any, M any](seeds []T, converter func(T) M) []M {
	models := make([]M, len(seeds))
	for i, seed := range seeds {
		models[i] = converter(seed)
	}

	return models
}

func updateSeedData[T any, M any](db *gorm.DB, seedName string, filename string, converter func(T) M) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read seed file %s: %w", filename, err)
	}

	currentHash := fmt.Sprintf("%x", sha256.Sum256(content))

	var seedVersion SeedVersion
	result := db.Where("name = ?", seedName).First(&seedVersion)

	if result.Error == nil && seedVersion.Hash == currentHash {
		return nil
	}

	var seedData []T
	if err := json.Unmarshal(content, &seedData); err != nil {
		return fmt.Errorf("failed to parse seed file %s: %w", filename, err)
	}

	data := processSeedData(seedData, converter)

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	switch m := any(data).(type) {
	case []models.Song:
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).CreateInBatches(m, 100).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert songs: %w", err)
		}

	case []models.Chart:
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).CreateInBatches(m, 100).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert songs: %w", err)
		}

	default:
		tx.Rollback()
		return fmt.Errorf("unsupported model type")
	}

	if result.Error == gorm.ErrRecordNotFound {
		if err := tx.Create(&SeedVersion{
			Name: seedName,
			Hash: currentHash,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		seedVersion.Hash = currentHash
		if err := tx.Save(&seedVersion).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func SeedDatabase(db *gorm.DB) error {
	if err := updateSeedData(db, "songs", "seeds/songs.json",
		func(s SongSeed) models.Song { return s.ToModel() }); err != nil {
		return fmt.Errorf("failed to seed songs: %w", err)
	}

	if err := updateSeedData(db, "charts", "seeds/charts.json",
		func(c ChartSeed) models.Chart { return c.ToModel() }); err != nil {
		return fmt.Errorf("failed to seed charts: %w", err)
	}

	return nil
}
