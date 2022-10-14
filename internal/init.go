package internal

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal().Msg(err.Error())
	}
}