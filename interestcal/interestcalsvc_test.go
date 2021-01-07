package interestcal_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
	"github.com/stretchr/testify/require"
)

func TestServiceServer_CreateInterest(t *testing.T) {
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
	verifyCredentials := `{"verify_user": { "email": "growth@drinnovations.us", "password": "kubernetes"}}`
	payload = strings.NewReader(verifyCredentials)

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

	var ctresp userauthnvalue.CreateTokenResponse
	decoder := json.NewDecoder(authnRes.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ctresp)
	require.Nil(t, err, "token response decoding should not be give error")
	fmt.Printf("ctresp is %v", ctresp)
	require.NotNil(t, ctresp.VerifiedUser.AccessToken, "access token should not be nil")

	accessToken := ctresp.VerifiedUser.AccessToken
	t.Logf("access accessToken is %v", accessToken)

	method = "POST"
	payload = strings.NewReader(`{
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

	url = fmt.Sprintf("%v/v1/interests", address)
	req, err = http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	istCalRes, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var ciresp interestvalue.CreateInterestResponse
	decoder = json.NewDecoder(istCalRes.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ciresp)
	require.Nil(t, err)
	require.Equal(t, 23.46, ciresp.Banks[0].Deposits[2].Delta, "delta for a deposit in a bank must match")
	require.Equal(t, 259.86, ciresp.Banks[0].Delta, "delta for a bank must match")
	require.Equal(t, 336.74, ciresp.Delta, "overall delta for all deposists in all banks must match")
}
