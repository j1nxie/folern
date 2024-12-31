package models

import "github.com/shopspring/decimal"

// Represents a score entry in the database, after processing from the Kamaitachi API.
type Score struct {
	ID        uint64          `gorm:"primaryKey" json:"-"`                             // The ID of the score entry.
	ChartID   string          `gorm:"uniqueIndex:idx_chart_song_user" json:"chart_id"` // The ID of the chart the score is associated with.
	SongID    int             `gorm:"uniqueIndex:idx_chart_song_user" json:"song_id"`  // The ID of the song the chart, and by extension, the score, is associated with.
	Score     int64           `json:"score"`                                           // The score, as displayed ingame.
	Lamp      ScoreLamp       `gorm:"type:score_lamp" json:"lamp"`                     // The clear status of the score.
	OverPower decimal.Decimal `json:"over_power"`                                      // The calculated OVER POWER value of the score.
	UserID    string          `gorm:"uniqueIndex:idx_chart_song_user" json:"-"`        // The ID of the user that the score belongs to.
	Chart     Chart           `gorm:"foreignKey:ChartID" json:"chart"`
	Song      Song            `gorm:"foreignKey:SongID" json:"song"`
}

type OverPowerStatsResponse struct {
	All     decimal.Decimal            `json:"all"`
	Genre   map[string]decimal.Decimal `json:"genre"`
	Version map[string]decimal.Decimal `json:"version"`
}

type OverPowerScoresResponse struct {
	Genre   map[string][]Score `json:"genre"`
	Version map[string][]Score `json:"version"`
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
