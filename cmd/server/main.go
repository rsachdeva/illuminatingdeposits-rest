package main

import (
	_ "expvar" // Register the expvar interestsvc
	"log"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"

	"github.com/rsachdeva/illuminatingdeposits-rest/appserver"
)

func main() {
	if err := appserver.ConfigureAndServe(); err != nil {
		log.Println("shutting down, error:", err)
		os.Exit(1)
	}
}
