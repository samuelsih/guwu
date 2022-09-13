package config

import (
	"log"
	"os"

	_ "embed"

	"github.com/jmoiron/sqlx"
)

//go:embed init.up.sql
var sqlSchema string

func ConnectAndInitCockroach() *sqlx.DB {
	db, err := sqlx.Open("pgx", os.Getenv("COCKROACH_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(sqlSchema)

	return db
}