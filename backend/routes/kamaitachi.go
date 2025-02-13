package routes

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/alitto/pond/v2"
	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/middleware"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KamaitachiHandler struct {
	db         *gorm.DB
	syncQueue  map[string]bool
	queueMutex sync.RWMutex
	pool       pond.Pool
}

func NewKamaitachiHandler(db *gorm.DB) *KamaitachiHandler {
	return &KamaitachiHandler{
		db:        db,
		syncQueue: make(map[string]bool),
		pool:      pond.NewPool(8),
	}
}

func (h *KamaitachiHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequireAuth)
	r.Get("/sync", h.syncScores)
	r.Get("/sync/status", h.syncStatus)

	return r
}

func (h *KamaitachiHandler) syncStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	h.queueMutex.RLock()
	inQueue := h.syncQueue[userID]
	h.queueMutex.RUnlock()

	models.SuccessResponse[any](w, http.StatusOK, "SUCCESSFULLY_RETURNED_SYNC_STATUS", map[string]bool{"syncing": inQueue})
}

func (h *KamaitachiHandler) syncScores(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	h.queueMutex.RLock()
	inQueue := h.syncQueue[userID]
	h.queueMutex.RUnlock()

	if inQueue {
		models.SuccessResponse[any](w, http.StatusAccepted, "SYNC_ALREADY_IN_PROGRESS", nil)
		return
	}

	h.queueMutex.Lock()
	h.syncQueue[userID] = true
	h.queueMutex.Unlock()

	h.pool.Submit(func() {
		defer func() {
			h.queueMutex.Lock()
			delete(h.syncQueue, userID)
			h.queueMutex.Unlock()
		}()

		h.processSyncRequest(userID)
	})

	models.SuccessResponse[any](w, http.StatusAccepted, "SUCCESSFULLY_STARTED_SYNC", nil)
}

var client = &http.Client{}

const kamaitachiURL string = "https://kamai.tachi.ac/api/v1/"

func (h *KamaitachiHandler) processSyncRequest(userID string) (int, int, error) {
	errorCount := 0
	kamaitachiURL := kamaitachiURL + "users/me/games/chunithm/Single/pbs/all"

	var creds models.UserAPIKey
	if err := h.db.Where("user_id = ?", userID).First(&creds).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("kt.sync", err.Error(), "user not authenticated with Kamaitachi")
			return 0, 0, err
		}

		logger.Error("kt.sync", err.Error(), "failed to fetch user API key")
		return 0, 0, err
	}

	apiKey := utils.DecryptAPIKey(creds.EncryptedAPIKey)

	req, _ := http.NewRequest("GET", kamaitachiURL, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("kt.sync", err.Error(), "failed to fetch PB data from Kamaitachi")
		return 0, 0, err
	}
	defer res.Body.Close()

	var ktRes models.KamaitachiResponse[models.KamaitachiPBResponse]
	if err := json.NewDecoder(res.Body).Decode(&ktRes); err != nil {
		logger.Error("kt.sync", err.Error(), "failed to decode PB data")
		return 0, 0, err
	}

	for _, pb := range ktRes.Body.PBs {
		var chart models.Chart
		result := h.db.Where("id = ?", pb.ChartID).First(&chart)

		if result.Error != nil {
			logger.Error("kt.sync", result.Error.Error(), "error when fetching chart for score")
			errorCount++
			continue
		}

		folernScore := models.Score{
			Score:     int64(pb.ScoreData.Score),
			ChartID:   pb.ChartID,
			SongID:    pb.SongID,
			Lamp:      pb.ScoreData.Lamp,
			UserID:    userID,
			OverPower: utils.CalculateOverpower(int64(pb.ScoreData.Score), chart.Level, string(pb.ScoreData.Lamp)),
		}

		if err := h.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "chart_id"},
				{Name: "song_id"},
				{Name: "user_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"score", "lamp", "over_power"}),
		}).Create(&folernScore).Error; err != nil {
			logger.Error("kt.sync", err.Error(), "failed to process score")
			errorCount++
			continue
		}
	}

	return len(ktRes.Body.PBs), errorCount, nil
}
