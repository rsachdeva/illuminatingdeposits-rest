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

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/http/httputil"
// 	"time"
// )
//
// // https://stackoverflow.com/questions/47591342/request-context-wont-close-connection
// type ResponseHandler func(resp *http.Response, err error) error
//
// func main() {
// 	url := runTestServer() //just for testing
// 	err := DoGetWithTimeout(url, stdOutHandler, time.Duration(6*time.Second))
// 	if err != nil {
// 		fmt.Printf("Error in GET request to server %s, error = %s", url, err)
// 	}
//
// }
//
// //DoGetWithTimeout - Makes a request with context that timesput
// func DoGetWithTimeout(getURL string, respHandler ResponseHandler, timeout time.Duration) error {
// 	fmt.Println("running now client..")
// 	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
// 	ctx = context.WithValue(ctx, "jmd", 5)
// 	v := ctx.Value("jmd")
// 	fmt.Println("v client is", v)
// 	defer cancelFunc()
// 	request, _ := http.NewRequestWithContext(ctx, "GET", getURL, nil)
//
// 	v = request.Context().Value("jmd")
// 	fmt.Println("v client request.Context() is", v)
//
// 	fResult := make(chan error, 0)
// 	go func() {
// 		fResult <- RoundTrip(request, respHandler)
// 	}()
// 	select {
// 	case <-ctx.Done():
// 		<-fResult //let the go routine end too
// 		return ctx.Err()
// 	case err := <-fResult:
// 		return err //any other errors in response
// 	}
// }
//
// //RoundTrip makes an http request and processes the response through a Response Handler func
// func RoundTrip(request *http.Request, respHandler ResponseHandler) error {
// 	return respHandler(http.DefaultClient.Do(request))
//
// }
//
// func stdOutHandler(resp *http.Response, err error) error {
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	//Handle the response
// 	//In this case we just print the body
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Printf("Body from response %s\n", string(body))
// 	return nil
// }
//
// func runTestServer() string {
// 	fmt.Println("starting run test server...")
// 	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		v := ctx.Value("jmd")
// 		fmt.Println("v server  is", v)
// 		incoming, _ := httputil.DumpRequest(r, false)
// 		fmt.Printf("Server: Incoming Request %s", string(incoming))
// 		time.Sleep(15 * time.Second) // Do difficult Job
// 		w.Write([]byte("Hello There!"))
// 	}))
// 	fmt.Println("running test server...")
// 	return slowServer.URL
// }

// func tlsOpts() grpc.DialOption {
// 	certFile := "tlsdocker/ca.crt"
// 	creds, err := credentials.NewClientTLSFromFile(certFile, "")
// 	if err != nil {
// 		log.Fatalf("loading certificate error is %v", err)
// 	}
// 	opts := grpc.WithTransportCredentials(creds)
// 	return opts
// }

// b, err := ioutil.ReadFile(certFile)
// if err != nil {
// return nil, err
// }
// cp := x509.NewCertPool()
// if !cp.AppendCertsFromPEM(b) {
// return nil, fmt.Errorf("credentials: failed to append certificates")
// }
// return NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil

// Buy Security with Go book Kindle for reference
// https://github.com/PacktPublishing/Security-with-Go/blob/master/Chapter09/https_client/https_client.go

/**
package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	// Load cert
	cert, err := tls.LoadX509KeyPair("cert.pem", "privKey.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Configure TLS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()
	responder := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: responder}

	// Use client to make request.
	// Ignoring response, just verifying connection accepted.
	_, err = client.Get("https://example.com")
	if err != nil {
		log.Println("Error making request. ", err)
	}
}
*/
// NewClientTLSFromFile constructs TLS credentials from the provided root
// certificate authority certificate file(s) to validate server connections. If
// certificates to establish the identity of the client need to be included in
// the credentials (eg: for mTLS), use NewTLS instead, where a complete
// tls.Config can be specified.
// serverNameOverride is for testing only. If set to a non empty string,
// it will override the virtual host name of authority (e.g. :authority header
// field) in requests.

// from search 'golang capool to tls -mTLS' to https://play.golang.org/p/WUgzKP5Jvh
// google search keywords got from https://github.com/cf-routing/golang-app 'Using one-way TLS (server cert only)'

// func NewClientTLSFromFile(certFile, serverNameOverride string) (*x509.CertPool, error) {
// 	b, err := ioutil.ReadFile(certFile)
// 	if err != nil {
// 		log.Printf("Error in reading file %v", certFile)
// 		return nil, err
// 	}
// 	cp := x509.NewCertPool()
// 	if !cp.AppendCertsFromPEM(b) {
// 		return nil, fmt.Errorf("credentials: failed to append certificates")
// 	}
// 	return cp, nil
// }
// certFile := "config/tls/ca.crt"
// cp, err := NewClientTLSFromFile(certFile, "")
// if err != nil {
// 	log.Fatalf("loading certificate error is %v", err)
// }
// fmt.Printf("CertPool cp is %v\n", cp)
//
// tlsConfig := &tls.Config{RootCAs: cp, InsecureSkipVerify: false}
//
// responder := &http.Transport{TLSClientConfig: tlsConfig, DisableKeepAlives: true}
// client := &http.Client{Transport: responder}
// _, err = client.Get("https://localhost:3000/v1/health")
// if err != nil {
// 	log.Println("Error making request. ", err)
// }

// tls client
// https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang
func tlsClient() (error, *http.Client) {
	caCert, err := ioutil.ReadFile("conf/tls/cacrtto.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	// // AppendCertsFromPEM attempts to parse a series of PEM encoded certificates.
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
func Base64EncodedString(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// health
func withoutTlsGetRequestHealth() {
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
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls client err is %v", err)
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

// create user
func withoutTlsPostRequestCreateUser() {
	client := &http.Client{}
	url := "http://localhost:3000/v1/users"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	e := Base64EncodedString("someone-e@drinnovations.us", "jmdjmd")
	fmt.Println("encoded string is %s", e)
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
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls client err is %v", err)
	}
	url := "https://localhost:3000/v1/users"
	method := "POST"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	e := Base64EncodedString("someone-dtlsjmd@drinnovations.us", "jmdjmd")
	fmt.Println("encoded string is %s", e)
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

// invest
func withoutTlsPostRequestCreateInvest() {
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
	err, client := tlsClient()
	if err != nil {
		log.Fatalf("tls client err is %v", err)
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
	//tlsGetRequestHealth()
	//withoutTlsPostRequestCreateUser()
	// tlsPostRequestCreateUser()
	// withoutTlsPostRequestCreateInvest()
	tlsPostRequestCreateInvest()
}
