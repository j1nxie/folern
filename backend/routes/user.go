package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/middleware"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequireAuth)
	r.Get("/me", h.getCurrentUser)
	r.Get("/{id}/stats", h.getStats)

	return r
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user_id").(string)

	var user models.User
	if err := h.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("user.getStats", err, "user not found")
			utils.Error(w, http.StatusNotFound, err)
			return
		}

		logger.Error("user.getCurrentUser", err, "failed to get user")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) getStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "me" {
		userID = r.Context().Value("user_id").(string)
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("user.getStats", err, "user not found")
			utils.Error(w, http.StatusNotFound, err)
			return
		}

		logger.Error("user.getStats", err, "failed to get user")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	var results []models.ScoreResponse
	if err := h.db.Model(&models.Score{}).
		Preload("Chart").
		Preload("Song").
		Raw(`
			WITH ranked_scores AS (
				SELECT
					*,
					ROW_NUMBER() OVER (PARTITION BY song_id ORDER BY over_power DESC) as rn
				FROM
					scores
				WHERE
					user_id = ?
			)
			SELECT
				*
			FROM
				ranked_scores
			WHERE
				rn = 1;
		`, userID).
		Find(&results).Error; err != nil {
		logger.Error("user.getStats", err, "failed to get user's OP stats")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	allOP := decimal.Zero

	genre := make(map[string][]models.ScoreResponse)
	genreOP := make(map[string]decimal.Decimal)

	version := make(map[string][]models.ScoreResponse)
	versionOP := make(map[string]decimal.Decimal)

	for _, item := range results {
		allOP = allOP.Add(item.OverPower)

		genre[item.Song.Genre] = append(genre[item.Song.Genre], item)
		genreOP[item.Song.Genre] = genreOP[item.Song.Genre].Add(item.OverPower)

		version[item.Song.Version] = append(version[item.Song.Version], item)
		versionOP[item.Song.Version] = versionOP[item.Song.Version].Add(item.OverPower)
	}

	response := models.OverPowerResponse{
		Stats: models.OverPowerResponseStats{
			All:     allOP,
			Genre:   genreOP,
			Version: versionOP,
		},
		Genre:   genre,
		Version: version,
	}

	utils.JSON(w, http.StatusOK, response)
}
