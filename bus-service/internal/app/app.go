package app

import (
	"database/sql"
	"github.com/perkzen/mbus/bus-service/internal/config"
	databasepackage "github.com/perkzen/mbus/bus-service/internal/database"
	"log"
)

type Application struct {
	Logger *log.Logger
	Env    *config.Environment
	DB     *sql.DB
}

func NewApplication(env *config.Environment) (*Application, error) {

	pgDb := databasepackage.NewPostgresDB(env.PostgresURL)
	dbConn, err := pgDb.Open()
	if err != nil {
		return nil, err
	}

	return &Application{
		Env: env,
		DB:  dbConn,
	}, nil
}
