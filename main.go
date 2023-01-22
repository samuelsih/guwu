package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/pkg/logger"
	"github.com/samuelsih/guwu/pkg/securer"
)

var (
	debug = flag.Bool("debug", false, "set log level to debug")
	remigrate = flag.Bool("fresh", false, "drop all table and migrate")
)

func main() {
	flag.Parse()
	logger.SetMode(*debug)

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

	if *remigrate {
		logger.SysInfo("Drop the table and remigrate")
		if err := config.LoadPostgresExtension(db); err != nil {
			logger.SysFatal("error on load postgres extension: %v", err)
		}

		if err := config.MigrateAll(db); err != nil {
			logger.SysFatal("error on migrate: %v", err)
		}
	}

	redisDB := config.NewRedis(redisHost, redisPassword)
	if redisDB == nil {
		logger.SysFatal("redisDB is nil")
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
		logger.SysFatal("cannot run server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutdown server: %v", err)
		return
	}

	logger.SysInfo("Server stopped")
}
