package database

import (
	"fmt"
	"os"

	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"gorm.io/gorm"
)

func updateSeedData(db *gorm.DB, seedName string, filename string) error {
	buffer, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read seed file %s: %w", filename, err)
	}

	content := string(buffer)

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	switch seedName {
	case "songs":
		if err := tx.Raw(`
			INSERT INTO songs
				(id, title, artist, version, genre)
			SELECT
				j ->> 'id', j ->> 'title', j ->> 'artist', j ->> 'data' ->> 'displayVersion', j ->> 'data' ->> 'genre'
			FROM json_array_elements(?) as j
			ON CONFLICT
				(id)
			DO UPDATE SET
				(title, artist, version, genre) = (EXCLUDED.title, EXCLUDED.artist, EXCLUDED.version, EXCLUDED.genre)
			RETURNING id, xmax <> 0 AS updated;
		`, content).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert songs: %w", err)
		}

	case "charts":
		type ChartResult struct {
			ID      string `gorm:"column:id"`
			Updated bool   `gorm:"column:updated"`
		}

		var results []ChartResult

		if err := tx.Raw(`
			INSERT INTO
				charts (id, song_id, level)
			SELECT
				j ->> 'chartID', (j ->> 'songID')::INTEGER, (j ->> 'levelNum')::REAL
			FROM json_array_elements(?) as j
			ON CONFLICT
				(id)
			DO UPDATE SET
				(song_id, level) = (EXCLUDED.song_id, EXCLUDED.level)
			WHERE
				(charts.song_id, charts.level) IS DISTINCT FROM (EXCLUDED.song_id, EXCLUDED.level)
			RETURNING id, xmax <> 0 AS updated;
		`, content).Scan(&results).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upsert charts: %w", err)
		}

		for _, result := range results {
			if result.Updated {
				var affectedScores []models.Score

				if err := tx.Model(&models.Score{}).
					Preload("Chart").
					Where("chart_id = ?", result.ID).
					Find(&affectedScores).Error; err != nil {
					logger.Error("db.updateSeedData", err, "failed to query affected scores")
					continue
				}

				if len(affectedScores) == 0 {
					continue
				}

				logger.Operation("db.updateSeedData", fmt.Sprintf("recalculating OP for score on chart %s", result.ID))
				for _, score := range affectedScores {
					score.OverPower = utils.CalculateOverpower(score.Score, score.Chart.Level, string(score.Lamp))

					if err := tx.Save(&score).Error; err != nil {
						logger.Error("db.updateSeedData", err, "failed to save recalculated score")
						continue
					}
				}
			}
		}

	default:
		tx.Rollback()
		return fmt.Errorf("unsupported model type")
	}

	return tx.Commit().Error
}

func SeedDatabase(db *gorm.DB) error {
	if err := updateSeedData(db, "songs", "seeds/songs.json"); err != nil {
		return fmt.Errorf("failed to seed songs: %w", err)
	}

	if err := updateSeedData(db, "charts", "seeds/charts.json"); err != nil {
		return fmt.Errorf("failed to seed charts: %w", err)
	}

	return nil
}
