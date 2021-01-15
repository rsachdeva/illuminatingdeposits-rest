package main

import (
	"log"
	"os"
)

func svcAddress() string {
	address := "localhost:3000"
	address, ok := os.LookupEnv("DEPOSITS_REST_SERVICE_ADDRESS")
	if ok {
		log.Println("DEPOSITS_REST_SERVICE_ADDRESS:address to connect to service (if service in kubernetes should match for ingress) from env is", address)
	}
	return address
}
