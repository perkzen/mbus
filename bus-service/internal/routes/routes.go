package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/bus-service/internal/api"
	"github.com/perkzen/mbus/bus-service/internal/app"
	"github.com/perkzen/mbus/bus-service/internal/middlewares"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	middlewares.Init(r)

	r.Get("/health", api.MakeHandlerFunc(app.HealthCheck))

	r.Route("/bus-stations", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusStationHandler.ListBusStations))
		r.Get("/{code}", api.MakeHandlerFunc(app.BusStationHandler.FindBusStationByCode))
	})

	r.Route("/bus-lines", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusLineHandler.ListBusLines))
	})

	return r
}
