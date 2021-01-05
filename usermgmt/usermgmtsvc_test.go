// Adds test that starts a Http server and client tests the user mgmt service with http routing
package usermgmt_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
	"github.com/stretchr/testify/require"
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

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("in string the body is", string(body))

	var nu uservalue.User
	decoder := json.NewDecoder(res.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&nu)
	require.Nil(t, err, "user decording should not give error")
	fmt.Printf("nu is %v", nu)
	require.NotNil(t, nu.Uuid, "UUID should not be nil")
	require.Equal(t, "growth@drinnovations.us", nu.Email)
}
