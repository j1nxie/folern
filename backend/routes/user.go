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

	return r
}

func (h *UserHandler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		logger.Error("user.getCurrentUser", err, "failed to get user")
		utils.Error(w, http.StatusInternalServerError, err)
	}

	utils.JSON(w, http.StatusOK, user)
}
