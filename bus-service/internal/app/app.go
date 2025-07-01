package app

import (
	"context"
	"database/sql"
	"github.com/perkzen/mbus/bus-service/internal/api"
	"github.com/perkzen/mbus/bus-service/internal/config"
	"github.com/perkzen/mbus/bus-service/internal/db"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"github.com/perkzen/mbus/bus-service/migrations"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type Application struct {
	Logger            *log.Logger
	Env               *config.Environment
	DB                *sql.DB
	BusStationHandler *api.BusStationHandler
	BusLineHandler    *api.BusLineHandler
	DepartureHandler  *api.DepartureHandler
	Cache             *redis.Client
}

func NewApplication(env *config.Environment) (*Application, error) {

	pgDb, err := db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		return nil, err
	}

	err = pgDb.Ping()
	if err != nil {
		panic(err)
	}

	err = db.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	busStationStore := store.NewPostgresBusStationStore(pgDb)
	busStationHandler := api.NewBusStationHandler(busStationStore)

	departureHandler := api.NewDepartureHandler(rdb)

	busLineStore := store.NewPostgresBusLineStore(pgDb)
	busLineHandler := api.NewBusLineHandler(busLineStore)

	return &Application{
		Env:               env,
		DB:                pgDb,
		Cache:             rdb,
		BusStationHandler: busStationHandler,
		BusLineHandler:    busLineHandler,
		DepartureHandler:  departureHandler,
	}, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return api.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
