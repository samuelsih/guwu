package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/samuelsih/guwu"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	debug := flag.Bool("debug", false, "sets log level to debug")
	production := flag.Bool("production", false, "is in prod.")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	server := guwu.NewServer(*production)

	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, os.Interrupt)

	server.Run(stop)
}
