package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rsachdeva/illuminatingdeposits-rest/readenv"
)

func tlsClient() (*http.Client, error) {
	caCert, err := ioutil.ReadFile("conf/tls/cacrtto.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	// AppendCertsFromPEM attempts to parse a series of PEM encoded certificates.
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return client, nil
}

func requestGetDbHealth(client *http.Client, prefix string) {
	fmt.Println("executing tLSGetRequestHealth()")
	resp, err := client.Get(fmt.Sprintf("%vlocalhost:3000/v1/health", prefix))
	if err != nil {
		log.Fatalf("err in get is %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("err reading response %v", err)
	}
	fmt.Println("body is ", string(body))
}

func requestPostCreateUser(client *http.Client, prefix string) {
	fmt.Println("executing requestPostCreateUser()")
	url := fmt.Sprintf("%vlocalhost:3000/v1/users", prefix)
	method := "POST"
	payload := strings.NewReader(`{
           "name":            "Rohit Sachdeva",
		   "email":           "growth-b1@drinnovations.us",
		   "roles":           ["USER"],
           "password":        "kubernetes",
           "password_confirm": "kubernetes"
    }`)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func requestPostCreateInterest(client *http.Client, prefix string) {
	fmt.Println("executing requestPostCreateInterest()")
	url := fmt.Sprintf("%vlocalhost:3000/v1/interests", prefix)
	method := "POST"
	payload := strings.NewReader(`{
		  "banks": [
			{
			  "name": "HAPPIEST",
			  "deposits": [
				{
				  "account": "1234",
				  "annualType": "Checking",
				  "annualRate%": 0,
				  "years": 1,
				  "amount": 100
				},
				{
				  "account": "1256",
				  "annualType": "CD",
				  "annualRate%": 24,
				  "years": 2,
				  "amount": 10700
				},
				{
				  "account": "1111",
				  "annualType": "CD",
				  "annualRate%": 1.01,
				  "years": 10,
				  "amount": 27000
				}
			  ]
			},
			{
			  "name": "NICE",
			  "deposits": [
				{
				  "account": "1234",
				  "annualType": "Brokered CD",
				  "annualRate%": 2.4,
				  "years": 7,
				  "amount": 10990
				}
			  ]
			},
			{
			  "name": "ANGRY",
			  "deposits": [
				{
				  "account": "1234",
				  "annualType": "Brokered CD",
				  "annualRate%": 5,
				  "years": 7,
				  "amount": 10990
				},
				{
				  "account": "9898",
				  "annualType": "CD",
				  "annualRate%": 2.22,
				  "years": 1,
				  "amount": 5500
				}
			  ]
			}
		  ]
		}`)
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	tlsEnabled := readenv.TlsEnabled()
	fmt.Println("tls is", tlsEnabled)

	var client *http.Client
	var err error
	var prefix string

	client = http.DefaultClient
	prefix = "http://"
	if tlsEnabled {
		client, err = tlsClient()
		if err != nil {
			log.Fatalf("tls client err is %v", err)
		}
		prefix = "https://"
	}

	requestGetDbHealth(client, prefix)
	requestPostCreateUser(client, prefix)
	requestPostCreateInterest(client, prefix)
}
