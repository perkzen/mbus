package main

import (
	"errors"
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/app"
	"github.com/perkzen/mbus/bus-service/internal/config"
	"github.com/perkzen/mbus/bus-service/internal/server"
	"log"
	"net/http"
)

func main() {
	env, _ := config.LoadEnvironment()
	restApp, _ := app.NewApplication(env)
	httpServer := server.NewHttpServer(restApp)

	defer restApp.DB.Close()

	done := make(chan bool, 1)

	go server.GracefulShutdown(httpServer, done)

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")

}
