package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/repository"

	"github.com/spitfy/gofermart/internal/app"

	"github.com/spitfy/gofermart/internal/handler"
)

func main() {
	cfg := config.GetConfig()
	db, err := repository.NewDBStore(cfg)
	if err != nil {
		log.Fatalf("Error database: %w", err)
	}
	defer db.Close()
	a, err := app.NewApp(cfg, db)
	if err != nil {
		log.Fatal(err)
	}
	router := handler.NewRouter(a.S)
	srv := app.NewServer(cfg, router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %s", err)
		}
	}()

	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
