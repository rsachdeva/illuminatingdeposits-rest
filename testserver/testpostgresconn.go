package testserver

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/rsachdeva/illuminatingdeposits-rest/tools/dbcli/schema"
)

func PostgresConnect() (*sqlx.DB, *dockertest.Pool, *dockertest.Resource) {
	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	// Construct url.
	pgURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pw, _ := pgURL.User.Password()
	env := []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + pgURL.Path,
	}

	resource, err := pool.Run("postgres", "13-alpine", env)

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *sqlx.DB
	pool.MaxWait = 115 * time.Second
	if err = pool.Retry(func() error {
		var err error
		// sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "postgres"))
		// postgres://postgres:postgres@db/postgres?sslmode=disable&timezone=utc
		connStr := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres?sslmode=disable&timezone=utc", resource.GetPort("5432/tcp"))
		log.Println("connStr is ", connStr)
		db, err = sqlx.Open("postgres", connStr)
		if err != nil {
			log.Printf(" sql open connection err is %v", err)
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = schema.Migrate(db)
	if err != nil {
		log.Fatalf("Could not bring the schema for db up to date with schema migrations: %v", err)
	}

	return db, pool, resource
}
