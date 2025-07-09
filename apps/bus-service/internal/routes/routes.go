package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/apps/bus-service/internal/api"
	"github.com/perkzen/mbus/apps/bus-service/internal/app"
	"github.com/perkzen/mbus/apps/bus-service/internal/middleware"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	middleware.Init(r)

	r.Get("/api/health", api.MakeHandlerFunc(app.HealthCheck))

	r.Route("/api/bus-stations", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusStationHandler.GetBusStations))
		r.Get("/{code}", api.MakeHandlerFunc(app.BusStationHandler.GetBusStationByCode))
	})

	r.Route("/api/bus-lines", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.BusLineHandler.GetBusLines))
	})

	r.Route("/api/departures", func(r chi.Router) {
		r.Get("/", api.MakeHandlerFunc(app.DepartureHandler.GetDepartures))
	})

	return r
}
