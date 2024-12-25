package models

import "github.com/shopspring/decimal"

// Represents a score entry in the database, after processing from the Kamaitachi API.
type Score struct {
	ID        uint64          `gorm:"primaryKey" json:"id"`        // The ID of the score entry.
	ChartID   string          `json:"chart_id"`                    // The ID of the chart the score is associated with.
	SongID    string          `json:"song_id"`                     // The ID of the song the chart, and by extension, the score, is associated with.
	Score     int64           `json:"score"`                       // The score, as displayed ingame.
	Lamp      ScoreLamp       `gorm:"type:score_lamp" json:"lamp"` // The clear status of the score.
	OverPower decimal.Decimal `json:"over_power"`                  // The calculated OVER POWER value of the score.
	UserID    string          `json:"user_id"`                     // The ID of the user that the score belongs to.
}

// Represents a chart in the database.
type Chart struct {
	ID     string          `gorm:"primaryKey" json:"id"` // The chart's ID.
	SongID int             `json:"song_id"`              // The song's ID, as well as its ingame ID.
	Level  decimal.Decimal `json:"level"`                // The chart's internal level.
}

// Represents a song in the database.
type Song struct {
	ID      int    `gorm:"primaryKey;autoincrement:false" json:"id"` // The song's ID, as well as its ingame ID.
	Title   string `json:"title"`                                    // The song's title.
	Artist  string `json:"artist"`                                   // The song's artist.
	Version string `json:"version"`                                  // The game version the song was introduced in.
	Genre   string `json:"genre"`                                    // The ingame genre the song is categorized in.
}

// The clear status of a song.
type ScoreLamp string

const (
	Failed             ScoreLamp = "FAILED"
	Clear              ScoreLamp = "CLEAR"
	FullCombo          ScoreLamp = "FULL COMBO"
	AllJustice         ScoreLamp = "ALL JUSTICE"
	AllJusticeCritical ScoreLamp = "ALL JUSTICE CRITICAL"
)
