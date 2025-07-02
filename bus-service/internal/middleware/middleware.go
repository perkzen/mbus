package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Init configures the provided chi.Mux router with CORS, JSON response compression, request logging, and panic recovery middleware.
func Init(r *chi.Mux) {

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Compress(5, "application/json"))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

}
