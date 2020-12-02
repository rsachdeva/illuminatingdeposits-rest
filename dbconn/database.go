package dbconn

import (
	"context"
	"log"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The dbconn driver in use.
	"go.opencensus.io/trace"
)

// Config is the required properties to use the dbconn.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

// Open knows how to open a dbconn connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {

	// Define SSL mode.
	sslMode := "require"
	// Transport Layer Security, and its now-deprecated predecessor, Secure Sockets Layer, are cryptographic protocols designed to provide communications security over a computer network.
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	log.Println("u.String() is ", u.String())
	// postgres://postgres:postgres@Db/postgres?sslmode=disable&timezone=utc when connecting (for debugging)
	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to the dbconn. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	ctx, span := trace.StartSpan(ctx, "platform.DB.StatusCheck")
	defer span.End()

	// Run a simple query to determine connectivity. The Db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// dbconn but the dbconn has since gone away. Running this query forces a
	// round trip to the dbconn.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
