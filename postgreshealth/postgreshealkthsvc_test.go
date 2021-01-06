package postgreshealth_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth/healthvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/stretchr/testify/require"
)

func TestServiceServer_HealthOk(t *testing.T) {
	t.Parallel()

	cr := testserver.InitRestHttpTLS(t, true)
	client := cr.TestClient
	address := cr.URL
	fmt.Printf("address is %v\n", address)
	url := fmt.Sprintf("%v/v1/health", address)

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("err in get is %v", err)
	}
	defer resp.Body.Close()

	var ht healthvalue.Postgres
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ht)
	require.Nil(t, err, "should be able to decode the error response")
	require.Equal(t, "Postgres Db Ok", ht.Status)

}

func TestServiceServer_HealthNotOk(t *testing.T) {
	t.Parallel()

	cr := testserver.InitRestHttpTLS(t, true)

	// disconnect db
	err := cr.PostgresClient.Close()
	if err != nil {
		t.Fatalf("Could not disconnect postgres db: %v", err)
	}

	client := cr.TestClient
	address := cr.URL
	fmt.Printf("address is %v\n", address)
	url := fmt.Sprintf("%v/v1/health", address)

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("err in get is %v", err)
	}
	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("err reading response %v", err)
	// }
	// fmt.Println("body is ", string(body))

	var ht healthvalue.Postgres
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ht)
	require.Nil(t, err, "should be able to decode the error response")
	require.Equal(t, "Postgres Db Not Ready", ht.Status)

}
