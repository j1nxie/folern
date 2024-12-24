package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/j1nxie/folern/database"
	"github.com/j1nxie/folern/logger"
	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/routes"
	"github.com/j1nxie/folern/utils"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var oauth2Config *oauth2.Config

func main() {
	database.InitDB()
	logger.InitLogger()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.Error(w, http.StatusNotFound, models.FolernError{Message: "route not found"})
	})

	authHandler := routes.NewAuthHandler(database.DB)
	userHandler := routes.NewUserHandler(database.DB)
	statusHandler := routes.NewStatusHandler()

	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/users", userHandler.Routes())
		r.Mount("/status", statusHandler.Routes())
	})

	logger.Operation("main.startup", "folern listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
