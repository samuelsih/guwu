package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/pkg/logger"
	"github.com/samuelsih/guwu/pkg/mail"
	"github.com/samuelsih/guwu/pkg/securer"
)

var (
	debug     = flag.Bool("debug", false, "set log level to debug")
	remigrate = flag.Bool("fresh", false, "drop all table and migrate")
)

func main() {
	flag.Parse()
	logger.SetMode(*debug)

	logger.SysInfof("Debug Mode: %v", *debug)

	secretKey := os.Getenv("SECURER_SECRET_KEY")
	if secretKey == "" {
		secretKey = "0f5297b6f0114171e9de547801b1e8bb929fe1d091e63c6377a392ec1baa3d0b"
	}

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

	mailHost := os.Getenv("MAIL_HOST")
	if mailHost == "" {
		mailHost = "localhost"
	}

	envMailPort := os.Getenv("MAIL_PORT")
	if envMailPort == "" {
		envMailPort = "1025"
	}

	mailPort, _ := strconv.Atoi(envMailPort)

	mailUser := os.Getenv("MAIL_USERNAME")
	if mailUser == "" {
		mailUser = "debuggerMail"
	}

	mailPassword := os.Getenv("MAIL_PASSWORD")

	mailEmail := os.Getenv("MAIL_EMAIL")
	if mailEmail == "" {
		mailEmail = "info@company.com"
	}

	totpSecretKey := os.Getenv("TOTP_SECRET")
	if totpSecretKey == "" {
		totpSecretKey = "4S62BZNFXXSZLCRO"
	}

	db := config.ConnectPostgres(dsn)
	securer.SetSecret(secretKey)
	redisDB := config.NewRedis(redisHost, redisPassword)
	
	mailer, err := mail.NewClient(mailHost, mailPort, mailEmail, mailPassword, mailUser, mailEmail)
	if err != nil {
		logger.SysFatal("error mailer: " + err.Error())
	}

	router := chi.NewRouter()

	if *remigrate {
		logger.SysInfo("Drop the table and remigrate")
		if err := config.LoadPostgresExtension(db); err != nil {
			logger.SysFatal("error on load postgres extension: %v", err)
		}

		if err := config.MigrateAll(db); err != nil {
			logger.SysFatal("error on migrate: %v", err)
		}
	}

	deps := Dependencies{
		DB:    db,
		Redis: redisDB,
		Mailer: mailer,
	}

	RunServer(router, ":"+port, deps)
}
