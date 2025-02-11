package routes

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/models"
)

type StatusHandler struct{}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (h *StatusHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.getStatus)

	return r
}

func (h *StatusHandler) getStatus(w http.ResponseWriter, r *http.Request) {
	serverTime := time.Now().Format(time.RFC3339)
	buildInfo, _ := debug.ReadBuildInfo()
	version := buildInfo.Main.Version

	models.SuccessResponse(w, http.StatusOK, "SUCCESSFULLY_RETURNED_STATUS", models.Status{ServerTime: serverTime, Version: version})
}
