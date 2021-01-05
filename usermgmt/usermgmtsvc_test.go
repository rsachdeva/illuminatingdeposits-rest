// Adds test that starts a Http server and client tests the user mgmt service with http routing
package usermgmt_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
)

func TestServiceServer_CreateUser(t *testing.T) {
	t.Parallel()

	cr := testserver.InitRestHttp(t, true)
	client := cr.TestClient
	address := cr.URL
	fmt.Printf("address is %v\n", address)
	url := fmt.Sprintf("%v/v1/users", address)
	method := "POST"
	usr := `{
           "name":            "Rohit Sachdeva",
		   "email":           "growth@drinnovations.us",
		   "roles":           ["USER"],
           "password":        "kubernetes",
           "password_confirm": "kubernetes"}`

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
	fmt.Println("in string the body is", string(body))
}
