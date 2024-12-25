package models

// Representation of a response from the Kamaitachi API.
type KamaitachiResponse[T any] struct {
	Success     bool   `json:"success"`
	Body        T      `json:"body,omitempty"`
	Description string `json:"description"`
}

type KamaitachiPBResponse struct {
	PBs []KamaitachiPB `json:"pbs"`
}

// A simplified representation of the PB document from the Kamaitachi API.
type KamaitachiPB struct {
	ChartID   string              `json:"chartID"`   // The ID of the chart the PB is associated with.
	ScoreData KamaitachiScoreData `json:"scoreData"` // The score data of the PB.
	SongID    int                 `json:"songID"`    // The ID of the song the chart is associated with. This is also the song's ingame ID.
}

// A simplified representation of score data from the Kamaitachi API.
type KamaitachiScoreData struct {
	Score int       `json:"score"` // The score, as displayed ingame.
	Lamp  ScoreLamp `json:"lamp"`  // The clear status of the score.
}

type KamaitachiAPITokenResponse struct {
	UserID        int             `json:"userID"`
	Token         string          `json:"token"`
	Identifier    string          `json:"identifier"`
	Permissions   map[string]bool `json:"permissions"`
	FromAPIClient string          `json:"fromAPIClient"`
}
