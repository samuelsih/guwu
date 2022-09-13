package main

import (
	"os"
	"os/signal"

	"github.com/samuelsih/guwu"
)

func main() {
	server := guwu.NewServer()

	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, os.Interrupt)

	server.Run(stop)
}