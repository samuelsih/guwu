package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/pkg/securer"
)

func main() {
	var secretKeyBytes [32]byte
	secretKey := os.Getenv("SECURER_SECRET_KEY")
	if secretKey == "" {
		secretKey = "0f5297b6f0114171e9de547801b1e8bb929fe1d091e63c6377a392ec1baa3d0b"
	}

	copy(secretKeyBytes[:], secretKey)
	securer.SetSecret(secretKeyBytes)

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=guwu sslmode=disable timezone=UTC connect_timeout=5"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = ""
	}

	db := config.ConnectPostgres(dsn)

	remigrate := flag.Bool("fresh", false, "drop all table and migrate")
	flag.Parse()

	log.Println("Remigrate:", *remigrate)

	if *remigrate {
		log.Println("Drop the table and remigrate")
		if err := config.LoadPostgresExtension(db); err != nil {
			log.Fatalf("error on load postgres extension: %v", err)
		}

		if err := config.MigrateAll(db); err != nil {
			log.Fatalf("error on migrate: %v", err)
		}
	}

	redisDB := config.NewRedis(redisHost, redisPassword)
	if redisDB == nil {
		log.Fatalf("redisDB is nil")
	}

	server := Server{
		Router: chi.NewRouter(),
		Dependencies: BusinessDeps{
			DB:      db,
			RedisDB: redisDB,
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
