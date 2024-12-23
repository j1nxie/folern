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
)

type AuthHandler struct {
	oauth2Config *oauth2.Config
	// TODO: db
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		oauth2Config: &oauth2.Config{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  "http://localhost:3000/api/auth/callback",
			Scopes: []string{
				"identify",
				"email",
			},
			Endpoint: discord.Endpoint,
		},
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/url", h.getAuthURL)
	r.Post("/callback", h.handleCallback)

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
	var req models.CallbackRequest

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

	if r.URL.Query().Get("state") != cookie.Value {
		logger.Error("auth.callback", err, "invalid state", "expected", cookie.Value, "actual", r.URL.Query().Get("state"))
		utils.Error(w, http.StatusBadRequest, models.FolernError{Message: "invalid state"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("auth.callback", err, "failed to decode request")
		utils.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.oauth2Config.Exchange(r.Context(), req.Code)
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

	user := models.UserInfo{
		ID:       discordUser.ID,
		Email:    discordUser.Email,
		Username: discordUser.Username,
		Avatar:   discordUser.Avatar,
	}

	jwtToken, err := utils.GenerateJWT(user)
	if err != nil {
		logger.Error("auth.callback", err, "failed to generate jwt")
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.JSON(w, http.StatusCreated, models.AuthResponse{Token: jwtToken, User: &user})
}
