package testserver

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/errconv"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-rest/metriccnt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth"
	"github.com/rsachdeva/illuminatingdeposits-rest/recoverpanic"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt"
)

type ClientResult struct {
	URL            string
	TestClient     *http.Client
	PostgresClient *sqlx.DB
}

func InitRestHttpTLS(t *testing.T, allowPurge bool) ClientResult {
	log := log.New(os.Stdout, "DEPOSITSTESTS: ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)
	log.Println("Starting ServiceServer...")
	db, pool, resource := PostgresConnect()

	shutdownCh := make(chan os.Signal, 1)
	m := muxhttp.NewRouter(shutdownCh, log,
		reqlog.NewMiddleware(log),
		errconv.NewMiddleware(log),
		metriccnt.NewMiddleware(),
		recoverpanic.NewMiddleware(log))

	log.Println("Registering REST json PostgresHealthService...")
	postgreshealth.RegisterSvc(db, m)
	log.Println("Registering REST json UserMgmtService...")
	usermgmt.RegisterSvc(db, m)
	log.Println("Registering REST json UserAuthenticationService...")
	userauthn.RegisterSvc(db, m)
	log.Println("Registering REST json InterestCalService...")
	interestcal.RegisterSvc(log, m)

	s := httptest.NewTLSServer(m)

	t.Cleanup(func() {
		log.Println("Shutting down / stopping  the server...")
		log.Println("Stopping the server...")
		s.Close()
		t.Logf("Purge allowed is %v", allowPurge)
		if allowPurge {
			t.Log("Purging dockertest for postgres")
			err := pool.Purge(resource)
			if err != nil {
				t.Fatalf("Could not purge container: %v", err)
			}
		}
		log.Println("End of program")
	})

	cr := ClientResult{
		URL:            s.URL,
		TestClient:     s.Client(),
		PostgresClient: db,
	}

	log.Printf("cr is %v", cr)
	return cr
}
