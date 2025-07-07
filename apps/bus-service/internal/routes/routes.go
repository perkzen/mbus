package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/bus-service/internal/api"
	"github.com/perkzen/mbus/bus-service/internal/app"
	"github.com/perkzen/mbus/bus-service/internal/middleware"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	middleware.Init(r)

	r.Get("/health", api.MakeHandlerFunc(app.HealthCheck))

	r.Route("/bus-stations", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusStationHandler.GetBusStations))
		r.Get("/{code}", api.MakeHandlerFunc(app.BusStationHandler.GetBusStationByCode))
	})

	r.Route("/bus-lines", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusLineHandler.GetBusLines))
	})

	r.Route("/departures", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.DepartureHandler.GetDepartures))
	})

	return r
}
