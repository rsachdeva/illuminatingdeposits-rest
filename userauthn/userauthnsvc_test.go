package userauthn_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
	"github.com/stretchr/testify/require"
)

func TestServer_AccessTokenCreation(t *testing.T) {
	tt := []struct {
		name              string
		verifyCredentails string
		authnTestFunc     func(r io.ReadCloser)
	}{
		{
			name:              "Allowed",
			verifyCredentails: `{"verify_user": { "email": "growth@drinnovations.us", "password": "kubernetes"}}`,
			authnTestFunc: func(r io.ReadCloser) {
				var ctres userauthnvalue.CreateTokenResponse
				decoder := json.NewDecoder(r)
				decoder.DisallowUnknownFields()
				err := decoder.Decode(&ctres)
				require.Nil(t, err, "token response decoding should not be give error")
				fmt.Printf("ctres is %v", ctres)
				require.NotNil(t, ctres.VerifiedUser.AccessToken, "access token should not be nil")
			},
		},
		{
			name:              "NotAllowed",
			verifyCredentails: `{"verify_user": { "email": "growth@drinnovationsus", "password": "kubernete"}}`,
			authnTestFunc: func(r io.ReadCloser) {
				var errResp struct {
					Error string `json:"error"`
				}
				decoder := json.NewDecoder(r)
				decoder.DisallowUnknownFields()
				err := decoder.Decode(&errResp)
				require.Nil(t, err, "should be able to decode the error response")
				require.Equal(t, "NotFound Error: User not found for email growth@drinnovationsus: sql: no rows in result set", errResp.Error)
			},
		},
	}

	for _, tc := range tt {
		tc := tc // capture range variable https://golang.org/pkg/testing/#hdr-Subtests_and_Sub_benchmarks
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cr := testserver.InitRestHttpTLS(t, true)
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

			payload = strings.NewReader(tc.verifyCredentails)

			url = fmt.Sprintf("%v/v1/users/token", address)
			req, err = http.NewRequest(method, url, payload)

			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Add("Content-Type", "application/json")

			authnRes, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			defer authnRes.Body.Close()

			log.Printf("authnRes is %v", authnRes)

			// body, err := ioutil.ReadAll(authnRes.Body)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// fmt.Println("in string the body is", string(body))
			//
			// var resErr struct {
			// 	Error string `json:"error"`
			// }
			// err = json.Unmarshal(body, &resErr)
			// require.Equal(t, "NotFound Error: User not found for email growth@drinnovationsus: sql: no rows in result set", resErr.Error)

			// var errResp struct {
			// 	Error string `json:"error"`
			// }
			// decoder := json.NewDecoder(authnRes.Body)
			// decoder.DisallowUnknownFields()
			// err = decoder.Decode(&errResp)
			// require.Nil(t, err, "should be able to decode the error response")
			// require.Equal(t, "NotFound Error: User not found for email growth@drinnovationsus: sql: no rows in result set", errResp.Error)

			tc.authnTestFunc(authnRes.Body)
		})
	}
}
