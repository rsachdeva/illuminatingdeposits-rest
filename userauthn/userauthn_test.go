package userauthn_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
	"github.com/stretchr/testify/require"
)

func TestServiceServer_CreateToken(t *testing.T) {
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

	vusr := `{
		"verify_user": {
			"email": "growth@drinnovations.us",
            "password": "kubernetes"
		}
	}`

	payload = strings.NewReader(vusr)

	url = fmt.Sprintf("%v/v1/users/token", address)
	req, err = http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resToken, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resToken.Body.Close()

	log.Printf("resToken is %v", resToken)

	var ctres userauthnvalue.CreateTokenResponse
	decoder := json.NewDecoder(resToken.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ctres)
	require.Nil(t, err, "token response decoding should not be give error")
	fmt.Printf("ctres is %v", ctres)
	require.NotNil(t, ctres.VerifiedUser.AccessToken, "access token should not be nil")
}

func TestServiceServer_CreateTokenNotAllowed(t *testing.T) {
	t.Parallel()

}
