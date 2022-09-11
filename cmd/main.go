package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var addr, sqlAddr string

	fs := flag.NewFlagSet("guwu", flag.ExitOnError)
	fs.StringVar(&addr, "addr", ":8080", "HTTP server address")
	fs.StringVar(&sqlAddr, "sql", "", "SQL connection address")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error on parsing arguments: %v", err)
		os.Exit(1)
	}

	
}