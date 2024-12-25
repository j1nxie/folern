package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	oauth2Config *oauth2.Config
	db           *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		oauth2Config: &oauth2.Config{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  "https://localhost:3000/auth/callback",
			Scopes: []string{
				"identify",
				"email",
			},
			Endpoint: discord.Endpoint,
		},
		db: db,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/url", h.getAuthURL)
	r.Get("/callback", h.handleCallback)
	r.Get("/logout", h.logout)

	return r
}

func (h *AuthHandler) getAuthURL(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateState()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, models.FolernError{Message: "failed to generate state"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	url := h.oauth2Config.AuthCodeURL(state)

	utils.JSON(w, http.StatusOK, models.AuthURLResponse{URL: template.URL(url)})
}

func (h *AuthHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, models.FolernError{Message: "missing state cookie"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	if state != cookie.Value {
		logger.Error("auth.callback", err, "invalid state", "expected", cookie.Value, "actual", r.URL.Query().Get("state"))
		utils.Error(w, http.StatusBadRequest, models.FolernError{Message: "invalid state"})
		return
	}

	if code == "" {
		logger.Error("auth.callback", err, "invalid code")
		utils.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		logger.Error("auth.callback", err, "failed to exchange code")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	client := h.oauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		logger.Error("auth.callback", err, "failed to get user data")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	var discordUser models.DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		logger.Error("auth.callback", err, "failed to decode user data")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	var dbUser models.User
	result := h.db.Where("id = ?", discordUser.ID).First(&dbUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			dbUser = models.User{
				ID:       discordUser.ID,
				Email:    discordUser.Email,
				Username: discordUser.Username,
				Avatar:   discordUser.Avatar,
			}

			if err := h.db.Create(&dbUser).Error; err != nil {
				logger.Error("auth.callback", err, "failed to create user")
				utils.Error(w, http.StatusInternalServerError, err)
				return
			}
		} else {
			logger.Error("auth.callback", result.Error, "failed to query user")
			utils.Error(w, http.StatusInternalServerError, result.Error)
			return
		}
	} else {
		dbUser.Email = discordUser.Email
		dbUser.Username = discordUser.Username
		dbUser.Avatar = discordUser.Avatar
		if err := h.db.Save(&dbUser).Error; err != nil {
			logger.Error("auth.callback", err, "failed to update user")
			utils.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	jwtToken, err := utils.GenerateJWT(dbUser)
	if err != nil {
		logger.Error("auth.callback", err, "failed to generate jwt")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSON(w, http.StatusCreated, models.AuthResponse{Token: jwtToken, User: &dbUser})
}

func (h *AuthHandler) logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSON(w, http.StatusOK, "successfully logged out.")
}
