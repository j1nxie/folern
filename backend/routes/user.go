package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/middleware"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
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
	r.Get("/me/stats", h.getStats)

	return r
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user_id").(string)

	var user models.User
	if err := h.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		logger.Error("user.getCurrentUser", err, "failed to get user")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) getStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var results []models.ScoreResponse
	if err := h.db.Raw(`
		WITH ranked_scores AS (
			SELECT
				*,
				ROW_NUMBER() OVER (PARTITION BY song_id ORDER BY over_power DESC) AS rn
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
	`, userID).Scan(&results).Error; err != nil {
		logger.Error("user.getStats", err, "failed to get user's OP stats")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, results)
}
