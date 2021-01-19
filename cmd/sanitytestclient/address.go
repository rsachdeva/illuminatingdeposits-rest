package main

import (
	"fmt"
	"os"
)

func svcAddress() string {
	address := "localhost"
	addr := os.Getenv("DEPOSITS_REST_SERVICE_ADDRESS")
	if !(addr == "" || addr == "localhost") {
		return addr
	}
	return fmt.Sprintf("%v:%v", address, 3000)
}
