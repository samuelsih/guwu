package config

import (
	"github.com/jmoiron/sqlx"
	"log"

	_ "embed"
	_ "github.com/lib/pq"
)

//go:embed init.up.sql
var sqlUp string

func ConnectPostgres(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to postgres: %v", err)
		return nil
	}

	log.Println("Postgres connect")

	return db
}

func LoadPostgresExtension(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	return err
}

func MigrateAll(db *sqlx.DB) error {
	_, err := db.Exec(sqlUp)

	return err
}
