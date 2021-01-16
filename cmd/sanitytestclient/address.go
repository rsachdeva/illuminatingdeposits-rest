package main

import (
	"log"
	"os"
)

func svcAddress() string {
	address := os.Getenv("DEPOSITS_REST_SERVICE_ADDRESS")
	if address != "" {
		log.Println("DEPOSITS_REST_SERVICE_ADDRESS:address to connect to service (if service in kubernetes should match for ingress) from env is", address)
		return address
	}
	return "localhost:3000"
}
