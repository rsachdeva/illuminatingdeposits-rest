// Provides sanity test client with all http json REST requests with TLS and when required JWT Authentication.
// This is to help with quick check of overall system.
// It is useful when doing refactoring as well.
// Uses unique email every time to allow new user creation and uses access token for the newly created user.
// Replace already persisted email in requestPostCreateToken, if requestCreateUser is not desired,
// otherwise user not found error will happen.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/rsachdeva/illuminatingdeposits-rest/readenv"
)

const (
	address  = "localhost:3000"
	emailFmt = "growth-%v@drinnovations.us"
)

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
	email := fmt.Sprintf(emailFmt, uuid.New().String())

	nonAccessTokenRequests(client, prefix, email)
	accessToken := requestPostCreateToken(client, prefix, email, false)
	accessTokenRequiredRequests(accessToken, client, prefix)
}

func nonAccessTokenRequests(client *http.Client, prefix string, email string) {
	requestGetDbHealth(client, prefix)
	requestPostCreateUser(client, prefix, email)
}

func accessTokenRequiredRequests(accessToken string, client *http.Client, prefix string) {
	fmt.Println("accessToken to be sent for accessTokenRequiredRequests...", accessToken)
	requestPostCreateInterest(accessToken, client, prefix)
}

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
	resp, err := client.Get(fmt.Sprintf("%v%v/v1/health", prefix, address))
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

func requestPostCreateUser(client *http.Client, prefix string, email string) {
	fmt.Println("executing requestPostCreateUser()")
	url := fmt.Sprintf("%v%v/v1/users", prefix, address)
	method := "POST"
	usr := fmt.Sprintf(`{
           "name":            "Rohit Sachdeva",
		   "email":           "%v",
		   "roles":           ["USER"],
           "password":        "kubernetes",
           "password_confirm": "kubernetes"}`,
		email)
	fmt.Println("usr is ", usr)
	payload := strings.NewReader(usr)

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

func requestPostCreateInterest(accessToken string, client *http.Client, prefix string) {
	fmt.Println("executing requestPostCreateInterest()")
	url := fmt.Sprintf("%v%v/v1/interests", prefix, address)
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
	req.Header.Set("Authorization", "Bearer "+accessToken)

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
	fmt.Println("string(body) response: ")
	fmt.Println(string(body))
	fmt.Println("res.Status is", res.Status)
}
