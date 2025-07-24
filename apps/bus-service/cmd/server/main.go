package main

import (
	"errors"
	_ "github.com/perkzen/mbus/apps/bus-service/docs"
	"github.com/perkzen/mbus/apps/bus-service/internal/app"
	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/server"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"log"
	"net/http"
	"time"
)

// @title mubs Bus Service API
// @version 1.0
// @description This is the API documentation for the mubs Bus Service.
func main() {
	env, err := config.LoadEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to load environment variables: %v", err)
	}

	var restApp *app.Application

	err = utils.Retry("Initialize application (DB + Redis)", 10, 3*time.Second, func() error {
		var initErr error
		restApp, initErr = app.NewApplication(env)
		return initErr
	})

	if err != nil {
		log.Fatalf("❌ Could not initialize application after retries: %v", err)
	}

	httpServer := server.NewHttpServer(restApp)
	log.Printf("✅ Server is running at http://localhost%s\n", httpServer.Addr)
	log.Printf("Swagger documentation is available at http://localhost%s/swagger/index.html\n", httpServer.Addr)

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
