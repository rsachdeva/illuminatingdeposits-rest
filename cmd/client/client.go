package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
)

// https://stackoverflow.com/questions/47591342/request-context-wont-close-connection
type ResponseHandler func(resp *http.Response, err error) error

func main() {
	url := runTestServer() //just for testing
	err := DoGetWithTimeout(url, stdOutHandler, time.Duration(6*time.Second))
	if err != nil {
		fmt.Printf("Error in GET request to server %s, error = %s", url, err)
	}

}

//DoGetWithTimeout - Makes a request with context that timesput
func DoGetWithTimeout(getURL string, respHandler ResponseHandler, timeout time.Duration) error {
	fmt.Println("running now client..")
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	ctx = context.WithValue(ctx, "jmd", 5)
	v := ctx.Value("jmd")
	fmt.Println("v client is", v)
	defer cancelFunc()
	request, _ := http.NewRequestWithContext(ctx, "GET", getURL, nil)

	v = request.Context().Value("jmd")
	fmt.Println("v client request.Context() is", v)

	fResult := make(chan error, 0)
	go func() {
		fResult <- RoundTrip(request, respHandler)
	}()
	select {
	case <-ctx.Done():
		<-fResult //let the go routine end too
		return ctx.Err()
	case err := <-fResult:
		return err //any other errors in response
	}
}

//RoundTrip makes an http request and processes the response through a Response Handler func
func RoundTrip(request *http.Request, respHandler ResponseHandler) error {
	return respHandler(http.DefaultClient.Do(request))

}

func stdOutHandler(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//Handle the response
	//In this case we just print the body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body from response %s\n", string(body))
	return nil
}

func runTestServer() string {
	fmt.Println("starting run test server...")
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := ctx.Value("jmd")
		fmt.Println("v server  is", v)
		incoming, _ := httputil.DumpRequest(r, false)
		fmt.Printf("Server: Incoming Request %s", string(incoming))
		time.Sleep(15 * time.Second) // Do difficult Job
		w.Write([]byte("Hello There!"))
	}))
	fmt.Println("running test server...")
	return slowServer.URL
}
