package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
)

func requestPostCreateToken(client *http.Client, prefix string) string {
	fmt.Println("executing requestPostCreateToken()")
	url := fmt.Sprintf("%vlocalhost:3000/v1/users/token", prefix)
	method := "POST"
	payload := strings.NewReader(`{
			"verify_user": {
				"email": "growth-a91@drinnovations.us",
				"password": "kubernetes"
			}
    }`)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("string(body) response: ")
	fmt.Println(string(body))
	fmt.Println("res.Status is", res.Status)

	var ctresp userauthnvalue.CreateTokenResponse
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	// https://stackoverflow.com/questions/45122496/why-does-json-unmarshal-need-a-pointer-to-a-map-if-a-map-is-a-reference-type
	if err := decoder.Decode(&ctresp); err != nil {
		log.Fatalln("decode error: ", err)
	}
	fmt.Println("vusr.AccessToken is: ")
	return ctresp.VerifiedUser.AccessToken
}
