package main

import (
	"encoding/base64"
	"fmt"
)

func Base64EncodedString(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	e := Base64EncodedString("someone@drinnovations.us", "grocery")
	fmt.Printf("encoded string is %s", e)
}
