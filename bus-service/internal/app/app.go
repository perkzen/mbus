package app

import (
	"database/sql"
	"github.com/perkzen/mbus/bus-service/internal/api"
	"github.com/perkzen/mbus/bus-service/internal/config"
	databasepackage "github.com/perkzen/mbus/bus-service/internal/database"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"log"
	"net/http"
)

type Application struct {
	Logger            *log.Logger
	Env               *config.Environment
	DB                *sql.DB
	BusStationHandler *api.BusStationHandler
	BusLineHandler    *api.BusLineHandler
}

func NewApplication(env *config.Environment) (*Application, error) {

	pgDb := databasepackage.NewPostgresDB(env.PostgresURL)
	dbConn, err := pgDb.Open()
	if err != nil {
		return nil, err
	}

	busStationStore := store.NewPostgresBusStationStore(dbConn)
	busStationHandler := api.NewBusStationHandler(busStationStore)

	busLineStore := store.NewPostgresBusLineStore(dbConn)
	busLineHandler := api.NewBusLineHandler(busLineStore)

	return &Application{
		Env:               env,
		DB:                dbConn,
		BusStationHandler: busStationHandler,
		BusLineHandler:    busLineHandler,
	}, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
