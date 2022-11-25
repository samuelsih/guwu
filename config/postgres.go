package config

import (
	"github.com/rs/zerolog/log"

	_ "embed"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed init.up.sql
var sqlUp string

func ConnectPostgres(dsn string) *sqlx.DB {	
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable timezone=UTC connect_timeout=5 sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msg("Postgres serve on : " + dsn)

	return db
}

func MigrateAll(db *sqlx.DB) error {
	_, err := db.Exec(sqlUp)

	return err
}
