package main

import (
	"github.com/samuelsih/guwu/config"
)

func main() {
	db := config.ConnectAndInitCockroach()

	

	db.Close()
}