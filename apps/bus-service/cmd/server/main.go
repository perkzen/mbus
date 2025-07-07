package main

import (
	"errors"
	"github.com/perkzen/mbus/apps/bus-service/internal/app"
	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/server"
	"log"
	"net/http"
)

func main() {
	env, err := config.LoadEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to load environment variables: %v", err)
	}

	restApp, err := app.NewApplication(env)
	if err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}

	httpServer := server.NewHttpServer(restApp)

	defer restApp.DB.Close()
	defer restApp.Cache.Close()

	done := make(chan bool, 1)

	go server.GracefulShutdown(httpServer, done)

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("❌ HTTP server error: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")

}
