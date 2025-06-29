package main

import (
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/app"
	"github.com/perkzen/mbus/bus-service/internal/server"
	"log"
	"net/http"
)

func main() {

	restApp, _ := app.NewApplication()
	httpServer := server.NewHttpServer(restApp)

	done := make(chan bool, 1)

	go server.GracefulShutdown(httpServer, done)

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")

}
