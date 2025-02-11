package routes

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/middleware"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthHandler struct {
	discordOAuth2Config    *oauth2.Config
	kamaitachiOAuth2Config utils.KamaitachiOAuth2Config
	db                     *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		discordOAuth2Config: &oauth2.Config{
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURL:  "https://localhost:3000/auth/callback",
			Scopes: []string{
				"identify",
				"email",
			},
			Endpoint: discord.Endpoint,
		},
		kamaitachiOAuth2Config: utils.KamaitachiOAuth2Config{
			Config: &oauth2.Config{
				ClientID:     os.Getenv("KAMAITACHI_CLIENT_ID"),
				ClientSecret: os.Getenv("KAMAITACHI_CLIENT_SECRET"),
				RedirectURL:  "https://localhost:3000/auth/kt-callback",
				Scopes: []string{
					"customise_scores",
				},
				Endpoint: oauth2.Endpoint{
					AuthURL:   "https://kamai.tachi.ac/api/v1/oauth/request-auth",
					TokenURL:  "https://kamai.tachi.ac/api/v1/oauth/token",
					AuthStyle: oauth2.AuthStyleInParams,
				},
			},
		},
		db: db,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/url", h.getAuthURL)
	r.Get("/callback", h.handleDiscordCallback)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/kt-callback", h.handleKamaitachiCallback)
	})
	r.Get("/logout", h.logout)

	return r
}

func (h *AuthHandler) getAuthURL(w http.ResponseWriter, r *http.Request) {
	state, err := utils.GenerateState()
	if err != nil {
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_AUTH_URL")
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

	url := h.discordOAuth2Config.AuthCodeURL(state)

	models.SuccessResponse(w, http.StatusOK, "SUCCESSFULLY_CREATED_AUTH_URL", models.AuthURLResponse{URL: template.URL(url)})
}

func (h *AuthHandler) handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_MISSING_STATE")
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
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_INVALID_STATE")
		return
	}

	if code == "" {
		logger.Error("auth.callback", err, "invalid code")
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_INVALID_CODE")
		return
	}

	token, err := h.discordOAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		logger.Error("auth.callback", err, "failed to exchange code")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_EXCHANGE_CODE")
		return
	}

	client := h.discordOAuth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		logger.Error("auth.callback", err, "failed to get user data")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_DISCORD_USER_DATA")
		return
	}
	defer resp.Body.Close()

	var discordUser models.DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		logger.Error("auth.callback", err, "failed to decode user data")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_DISCORD_USER_DATA")
		return
	}

	dbUser := models.User{
		ID:       discordUser.ID,
		Email:    discordUser.Email,
		Username: discordUser.Username,
		Avatar:   discordUser.Avatar,
	}

	if err := h.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"email", "username", "avatar"}),
	}).Create(&dbUser).Error; err != nil {
		logger.Error("auth.callback", err, "failed to process user")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_UPDATE_USER_DATA")
		return
	}

	jwtToken, err := utils.GenerateJWT(dbUser)
	if err != nil {
		logger.Error("auth.callback", err, "failed to generate jwt")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GENERATE_JWT")
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

	models.SuccessResponse(w, http.StatusCreated, "SUCCESSFULLY_LOGGED_IN", models.AuthResponse{Token: jwtToken, User: &dbUser})
}

func (h *AuthHandler) handleKamaitachiCallback(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	code := r.URL.Query().Get("code")

	// TODO: neater errors here
	if code == "" {
		logger.Error("auth.kt-callback", models.FolernError{Message: "invalid code"}, "invalid code")
		models.ErrorResponse[any](w, http.StatusBadRequest, "ERROR_INVALID_CODE")
		return
	}

	token, err := h.kamaitachiOAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		logger.Error("auth.kt-callback", err, "failed to exchange code")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_EXCHANGE_CODE")
		return
	}

	client = h.kamaitachiOAuth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://kamai.tachi.ac/api/v1/me")
	if err != nil {
		logger.Error("auth.kt-callback", err, "failed to get user data")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_GET_KT_USER_DATA")
		return
	}
	defer resp.Body.Close()

	dbUserAPIKey := models.UserAPIKey{
		UserID:          userID,
		EncryptedAPIKey: utils.EncryptAPIKey(token.AccessToken),
	}

	if err := h.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"encrypted_api_key"}),
	}).Create(&dbUserAPIKey).Error; err != nil {
		logger.Error("auth.kt-callback", err, "failed to update user API key")
		models.ErrorResponse[any](w, http.StatusInternalServerError, "ERROR_FAILED_TO_UPDATE_API_KEY")
		return
	}

	models.SuccessResponse[any](w, http.StatusOK, "SUCCESSFULLY_AUTHENTICATED_WITH_KT", nil)
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

	models.SuccessResponse[any](w, http.StatusOK, "SUCCESSFULLY_LOGGED_OUT", nil)
}
