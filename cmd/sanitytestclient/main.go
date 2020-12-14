package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// tls sanitytestclient
func tlsClient() (error, *http.Client) {
	caCert, err := ioutil.ReadFile("conf/tls/cacrtto.pem")
	if err != nil {
		log.Fatal(err)
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
	return nil, client
}

// base64 encoded string
func Base64EncodedString(email, password string) string {
	auth := email + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// health
func withoutTlsGetRequestHealth() {
	fmt.Println("executing withoutTLSGetRequestHealth()")
	resp, err := http.Get("http://localhost:3000/v1/health")
	if err != nil {
		log.Fatalf("err in Get %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("err reading response %v", err)
	}
	fmt.Println("body is ", string(body))
}

func tlsGetRequestHealth() {
	fmt.Println("executing tLSGetRequestHealth()")
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls sanitytestclient err is %v", err)
	}
	resp, err := client.Get("https://localhost:3000/v1/health")
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

// create usermgmt
func withoutTlsPostRequestCreateUser() {
	fmt.Println("executing withoutTlsPostRequestCreateUser()")
	client := &http.Client{}
	url := "http://localhost:3000/v1/users"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	e := Base64EncodedString("someone-e@drinnovations.us", "jmdjmd")
	fmt.Printf("encoded string is %s\n", e)
	authHeader := fmt.Sprintf("Basic %s", e)
	req.Header.Add("Authorization", authHeader)

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

func tlsPostRequestCreateUser() {
	fmt.Println("executing tlsPostRequestCreateUser()")
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls sanitytestclient err is %v", err)
	}
	url := "https://localhost:3000/v1/users"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	e := Base64EncodedString("someone2-ijmd@drinnovations.us", "jmdjmd")
	fmt.Printf("encoded string is %s\n", e)
	authHeader := fmt.Sprintf("Basic %s", e)
	req.Header.Add("Authorization", authHeader)

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

// interestcal
func withoutTlsPostRequestCreateInvest() {
	fmt.Println("executing withoutTlsPostRequestCreateInvest()")
	client := &http.Client{}
	url := "http://localhost:3000/v1/interests"
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
				  "amount": 1
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

func tlsPostRequestCreateInvest() {
	fmt.Println("executing tlsPostRequestCreateInvest()")
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls sanitytestclient err is %v", err)
	}

	url := "https://localhost:3000/v1/interests"
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
				  "amount": 1
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
	// withoutTlsGetRequestHealth()
	tlsGetRequestHealth()
	//withoutTlsPostRequestCreateUser()
	tlsPostRequestCreateUser()
	// withoutTlsPostRequestCreateInvest()
	tlsPostRequestCreateInvest()
}
