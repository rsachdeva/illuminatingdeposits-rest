package main

import (
	_ "expvar" // Register the expvar interestsvc
	"log"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"
)

func main() {
	if err := ConfigureAndServe(); err != nil {
		log.Println("shutting down, error:", err)
		os.Exit(1)
	}
}
