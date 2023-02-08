package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/samuelsih/guwu/config"
	"github.com/samuelsih/guwu/pkg/env"
	"github.com/samuelsih/guwu/pkg/logger"
	"github.com/samuelsih/guwu/pkg/mail"
	"github.com/samuelsih/guwu/pkg/securer"
)

var (
	debug     = flag.Bool("debug", false, "set log level to debug")
	remigrate = flag.Bool("fresh", false, "drop all table and migrate")
)

type EnvConfig struct {
	SecretKey     string `env:"SECURER_SECRET_KEY" default:"0f5297b6f0114171e9de547801b1e8bb929fe1d091e63c6377a392ec1baa3d0b"`
	Dsn           string `env:"DB_DSN" default:"host=localhost port=5432 user=postgres password=postgres dbname=guwu sslmode=disable timezone=UTC connect_timeout=5"`
	Port          string `env:"PORT" default:"8080"`
	RedisHost     string `env:"REDIS_HOST" default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" default:""`
	MailHost      string `env:"MAIL_HOST" default:"localhost"`
	MailPort      int    `env:"MAIL_PORT" default:"1025"`
	MailUsername  string `env:"MAIL_USERNAME" default:"debuggerMail"`
	MailPassword  string `env:"MAIL_PASSWORD" default:""`
	MailEmail     string `env:"MAIL_EMAIL" default:"info@company.com"`
	TOTPSecret    string `env:"TOTP_SECRET" default:"4S62BZNFXXSZLCRO"`
}

func main() {
	flag.Parse()
	logger.SetMode(*debug)

	logger.SysInfof("Debug Mode: %v", *debug)

	var e EnvConfig

	if err := env.Fill(&e); err != nil {
		logger.SysFatal("Error getting from .env: " + err.Error())
	}

	db := config.ConnectPostgres(e.Dsn)
	securer.SetSecret(e.SecretKey)
	redisDB := config.NewRedis(e.RedisHost, e.RedisPassword)

	mailer, err := mail.NewClient(e.MailHost, e.MailPort, e.MailEmail, e.MailPassword, e.MailUsername, e.MailEmail)
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
		DB:     db,
		Redis:  redisDB,
		Mailer: mailer,
	}

	RunServer(router, ":"+e.Port, deps)
}
