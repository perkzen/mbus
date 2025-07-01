package app

import (
	"database/sql"
	"github.com/perkzen/mbus/bus-service/internal/api"
	"github.com/perkzen/mbus/bus-service/internal/config"
	"github.com/perkzen/mbus/bus-service/internal/db"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"github.com/perkzen/mbus/bus-service/migrations"
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

	pgDb, err := db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		return nil, err
	}

	err = db.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	busStationStore := store.NewPostgresBusStationStore(pgDb)
	busStationHandler := api.NewBusStationHandler(busStationStore)

	busLineStore := store.NewPostgresBusLineStore(pgDb)
	busLineHandler := api.NewBusLineHandler(busLineStore)

	return &Application{
		Env:               env,
		DB:                pgDb,
		BusStationHandler: busStationHandler,
		BusLineHandler:    busLineHandler,
	}, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
