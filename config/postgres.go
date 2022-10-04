package config

import (
	"github.com/rs/zerolog/log"

	_ "embed"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed init.up.sql
var sqlUp string

func ConnectPostgres(dsn string) *sqlx.DB {	
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable timezone=UTC connect_timeout=5"
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msg("Postgres serve on : " + dsn)

	return db
}

func MigrateAll(db *sqlx.DB) error {
	_, err := db.Exec(sqlUp)

	return err
}
