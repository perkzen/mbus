package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/perkzen/mbus/bus-service/internal/common/config"
	"log"
	"net/http"
)

type Application struct {
	Handler *chi.Mux
	Logger  *log.Logger
	Env     *config.Environment
}

func NewApplication() (*Application, error) {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return &Application{
		Handler: r,
	}, nil
}
