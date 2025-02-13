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
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database.InitDB()
	logger.InitLogger()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		models.ErrorResponse[any](w, http.StatusNotFound, "ERROR_NOT_FOUND")
	})

	authHandler := routes.NewAuthHandler(database.DB)
	kamaitachiHandler := routes.NewKamaitachiHandler(database.DB)
	statusHandler := routes.NewStatusHandler()
	userHandler := routes.NewUserHandler(database.DB)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", authHandler.Routes())
		r.Mount("/kamaitachi", kamaitachiHandler.Routes())
		r.Mount("/status", statusHandler.Routes())
		r.Mount("/users", userHandler.Routes())
	})

	logger.Operation("main.startup", "folern listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
