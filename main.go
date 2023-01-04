package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=guwu sslmode=disable timezone=UTC connect_timeout=5"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := config.ConnectPostgres(dsn)

	if b := flag.Bool("fresh", false, "drop all table and migrate"); *b {
		if err := config.LoadPostgresExtension(db); err != nil {
			log.Fatalf("error on load postgres extension: %v", err)
		}

		if err := config.MigrateAll(db); err != nil {
			log.Fatalf("error on migrate: %v", err)
		}
	}

	server := Server{
		Router: chi.NewRouter(),
		Dependencies: BusinessDeps{
			DB: db,
		},
	}

	server.MountHandlers()

	err := server.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutdown server: %v", err)
		return
	}

	log.Println("Server stopped")
}
