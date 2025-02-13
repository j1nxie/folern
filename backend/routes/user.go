package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/middleware"
	"github.com/j1nxie/folern/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	r.Get("/{id}", h.getCurrentUser)
	r.Get("/{id}/stats", h.getStats)
	r.Get("/{id}/scores", h.getScores)

	return r
}

func (h *UserHandler) retrieveScoresFromDB(userID string) ([]models.Score, error) {
	var results []models.Score
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
		logger.Error("user.retrieveScoresFromDB", err.Error(), "failed to get user's OP stats")
		return nil, err
	}

	return results, nil
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "me" {
		userID = r.Context().Value("user_id").(string)
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("user.getCurrentUser", err.Error(), "user not found")
			models.ErrorResponse[any](w, http.StatusNotFound, "ERROR_USER_NOT_FOUND")
			return
		}

		logger.Error("user.getCurrentUser", err.Error(), "failed to get user")
		models.ErrorResponse[any](w, http.StatusNotFound, "ERROR_FAILED_TO_GET_USER")
		return
	}

	models.SuccessResponse(w, http.StatusOK, "SUCCESSFULLY_RETURNED_USER", user)
}

func (h *UserHandler) getStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "me" {
		userID = r.Context().Value("user_id").(string)
	}
	category := r.URL.Query().Get("category")
	type_ := r.URL.Query().Get("type")

	if category == "" {
		category = "genres"
	}

	if type_ == "" {
		type_ = "possession"
	}

	if category != "versions" && category != "genres" {
		logger.Error("overpower.total", "invalid category", "invalid category")
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_INVALID_CATEGORY")
		return
	}

	if type_ != "possession" {
		logger.Error("overpower.total", "invalid type", "invalid type")
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_INVALID_TYPE")
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("user.getStats", err.Error(), "user not found")
			models.ErrorResponse[any](w, http.StatusNotFound, "ERROR_USER_NOT_FOUND")
			return
		}

		logger.Error("user.getStats", err.Error(), "failed to get user")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_USER")
		return
	}

	userScores, err := h.retrieveScoresFromDB(userID)
	if err != nil {
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_SCORES")
		return
	}

	var totalOP []models.TotalOverPower
	db := h.db.Model(&models.TotalOverPower{})

	versions := []string{
		"all",
		"chuni",
		"chuniplus",
		"air",
		"airplus",
		"star",
		"starplus",
		"amazon",
		"amazonplus",
		"crystal",
		"crystalplus",
		"paradise",
		"paradiselost",
		"new",
		"newplus",
		"sun",
		"sunplus",
		"luminous",
		"luminousplus",
		"verse",
	}

	genres := []string{
		"all",
		"POPS & ANIME",
		"niconico",
		"VARIETY",
		"東方Project",
		"イロドリミドリ",
		"ゲキマイ",
		"ORIGINAL",
	}

	var targetCategories []string
	switch category {
	case "genres":
		targetCategories = genres
	case "version":
		targetCategories = versions
	}

	sortingString := "ARRAY['" + strings.Join(targetCategories, "','") + "']"
	if err := db.Where("category IN ? AND type = ?", targetCategories, type_).
		Order(clause.OrderBy{
			Expression: gorm.Expr("array_position(" + sortingString + "::text[], category)"),
		}).
		Find(&totalOP).Error; err != nil {
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_OVERPOWER")
		return
	}

	if err := h.db.Model(&models.Score{}).
		Preload("Chart").
		Preload("Song").
		Raw(`
            WITH ranked_scores AS (
                SELECT *,
                ROW_NUMBER() OVER (PARTITION BY song_id ORDER BY over_power DESC) as rn
                FROM scores
                WHERE user_id = ?
            )
            SELECT * FROM ranked_scores WHERE rn = 1;
        `, userID).
		Find(&userScores).Error; err != nil {
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_OVERPOWER")
		return
	}

	var response []models.OverPowerStatsResponse

	currentAllOP := decimal.Zero
	currentOP := make(map[string]decimal.Decimal)
	for _, score := range userScores {
		switch category {
		case "genres":
			currentAllOP = currentAllOP.Add(score.OverPower)
			currentOP[score.Song.Genre] = currentOP[score.Song.Genre].Add(score.OverPower)
		case "versions":
			currentAllOP = currentAllOP.Add(score.OverPower)
			currentOP[score.Song.Version] = currentOP[score.Song.Version].Add(score.OverPower)
		}
	}

	for _, total := range totalOP {
		if total.Category == "all" {
			response = append(response, models.OverPowerStatsResponse{
				Category: "all",
				Type:     total.Type,
				Current:  currentAllOP,
				Maximum:  total.Value,
			})
		} else {
			response = append(response, models.OverPowerStatsResponse{
				Category: total.Category,
				Type:     total.Type,
				Current:  currentOP[total.Category],
				Maximum:  total.Value,
			})
		}
	}

	models.SuccessResponse(w, http.StatusOK, "SUCCESSFULLY_RETURNED_OVERPOWER", response)
}

func (h *UserHandler) getScores(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "me" {
		userID = r.Context().Value("user_id").(string)
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("user.getScores", err.Error(), "user not found")
			models.ErrorResponse[any](w, http.StatusNotFound, "ERROR_USER_NOT_FOUND")
			return
		}

		logger.Error("user.getScores", err.Error(), "failed to get user")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_USER")
		return
	}

	results, err := h.retrieveScoresFromDB(userID)
	if err != nil {
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_SCORES")
		return
	}

	genre := make(map[string][]models.Score)
	version := make(map[string][]models.Score)

	for _, item := range results {
		genre[item.Song.Genre] = append(genre[item.Song.Genre], item)
		version[item.Song.Version] = append(version[item.Song.Version], item)
	}

	response := models.OverPowerScoresResponse{
		Genre:   genre,
		Version: version,
	}

	models.SuccessResponse(w, http.StatusOK, "SUCCESSFULLY_RETURNED_SCORES", response)
}
