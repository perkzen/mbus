package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/perkzen/mbus/apps/bus-service/internal/api"
	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/db"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/openrouteservice"
	"github.com/perkzen/mbus/apps/bus-service/internal/service/departure"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/migrations"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"os"
)

type Application struct {
	Logger            *slog.Logger
	Env               *config.Environment
	DB                *sql.DB
	BusStationHandler *api.BusStationHandler
	BusLineHandler    *api.BusLineHandler
	DepartureHandler  *api.DepartureHandler
	Cache             *redis.Client
}

func NewApplication(env *config.Environment) (*Application, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	pgDb, err := db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres DB: %w", err)
	}

	if err := pgDb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres DB: %w", err)
	}

	if err := db.MigrateFS(pgDb, migrations.FS, "."); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: env.RedisPassword,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	busStationStore := store.NewPostgresBusStationStore(pgDb)
	busStationHandler := api.NewBusStationHandler(busStationStore, logger)

	busLineStore := store.NewPostgresBusLineStore(pgDb)
	busLineHandler := api.NewBusLineHandler(busLineStore, logger)

	orsApiClient := openrouteservice.NewAPIClient(env.ORSApiKey, rdb)
	departureStore := store.NewPostgresDepartureStore(pgDb)
	directionStore := store.NewPostgresDirectionStore(pgDb)
	departureService := departure.NewService(orsApiClient, rdb, busStationStore, departureStore, busLineStore, directionStore)
	departureHandler := api.NewDepartureHandler(departureService, logger)

	return &Application{
		Logger:            logger,
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
