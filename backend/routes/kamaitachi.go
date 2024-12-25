package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"gorm.io/gorm"
)

type KamaitachiHandler struct {
	db *gorm.DB
}

func NewKamaitachiHandler(db *gorm.DB) *KamaitachiHandler {
	return &KamaitachiHandler{db: db}
}

func (h *KamaitachiHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/sync", h.syncScores)

	return r
}

var client = &http.Client{}

const kamaitachiURL string = "https://kamai.tachi.ac/api/v1/"

func (h *KamaitachiHandler) syncScores(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	kamaitachiURL := kamaitachiURL + "users/me/games/chunithm/Single/pbs/all"

	var creds models.UserAPIKey
	if err := h.db.Where("user_id = ?", userID).First(&creds).Error; err != nil {
		logger.Error("kt.sync", err, "failed to fetch user API key")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	apiKey := utils.DecryptAPIKey(creds.EncryptedAPIKey)

	req, _ := http.NewRequest("GET", kamaitachiURL, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("kt.sync", err, "failed to fetch PB data from Kamaitachi")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer res.Body.Close()

	var ktRes models.KamaitachiResponse[models.KamaitachiPBResponse]
	if err := json.NewDecoder(res.Body).Decode(&ktRes); err != nil {
		logger.Error("kt.sync", err, "failed to decode PB data")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusOK, ktRes)
}
