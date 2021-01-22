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

func requestPostCreateToken(client *http.Client, prefix string, email string, useExpired bool) string {
	fmt.Println("=============executing requestPostCreateToken()=============")
	var token string
	if useExpired {
		storedTk, err := ioutil.ReadFile("cmd/sanitytestclient/expiredtoken.data")
		token = string(storedTk)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Expired JWT for testing", token)
		return token
	}
	token = createToken(client, prefix, email)
	log.Println("New JWT that can be used in next requests for Authentication\n", token)
	return token
}

func createToken(client *http.Client, prefix string, email string) string {
	url := fmt.Sprintf("%v%v/v1/users/token", prefix, svcAddress())
	method := "POST"
	vusr := fmt.Sprintf(`{
		"verify_user": {
			"email": "%v",
            "password": "kubernetes"
		}
	}`, email)
	log.Println("vusr data is", vusr)
	payload := strings.NewReader(vusr)

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
	fmt.Printf("ctresp.VerifiedUser is:%#v\n", ctresp.VerifiedUser)
	token := ctresp.VerifiedUser.AccessToken
	log.Println("New JWT that can be used in next requests for Authentication", token)
	return token
}
