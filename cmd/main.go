package main

import (
	"github.com/rs/zerolog"
	"github.com/samuelsih/guwu"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	guwu.Run()
}
