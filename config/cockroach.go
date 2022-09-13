package config

import (
	"log"
	"os"

	_ "embed"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jmoiron/sqlx"
)

//go:embed init.up.sql
var sqlInit string

//go:embed init.down.sql
var sqlDown string

func ConnectAndInitCockroach() *sqlx.DB {
	db, err := sqlx.Open("pgx", os.Getenv("COCKROACH_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Cant ping:", err)
	}
	
	db.MustExec(sqlDown)
	db.MustExec(sqlInit)

	return db
}