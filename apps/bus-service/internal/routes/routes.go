package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/apps/bus-service/internal/api"
	"github.com/perkzen/mbus/apps/bus-service/internal/app"
	"github.com/perkzen/mbus/apps/bus-service/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	middleware.Init(r)

	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Get("/health", api.MakeHandlerFunc(app.HealthCheck))

	r.Route("/api", func(r chi.Router) {
		r.Route("/bus-stations", func(r chi.Router) {
			r.Get("/", api.MakeHandlerFunc(app.BusStationHandler.GetBusStations))
			r.Get("/{id}", api.MakeHandlerFunc(app.BusStationHandler.GetBusStationByID))

		})

		r.Route("/bus-lines", func(r chi.Router) {
			r.Get("/", api.MakeHandlerFunc(app.BusLineHandler.GetBusLines))
		})

		r.Route("/departures", func(r chi.Router) {
			r.Get("/", api.MakeHandlerFunc(app.DepartureHandler.GetDepartures))
		})
	})

	return r
}
