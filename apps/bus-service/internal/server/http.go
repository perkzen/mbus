package server

import (
	"context"
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/app"
	"github.com/perkzen/mbus/bus-service/internal/routes"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func NewHttpServer(app *app.Application) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Env.Port),
		Handler: routes.RegisterRoutes(app),
	}
}

func GracefulShutdown(server *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
