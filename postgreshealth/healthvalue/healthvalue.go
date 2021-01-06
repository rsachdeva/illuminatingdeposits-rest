package healthvalue

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
)

type Postgres struct {
	Status string `json:"status"`
}

// StatusCheck returns nil if it can successfully talk to the postgresconn. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	ctx, span := trace.StartSpan(ctx, "platform.DB.StatusCheck")
	defer span.End()

	// Run a simple query to determine connectivity. The Db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// postgresconn but the postgresconn has since gone away. Running this query forces a
	// round trip to the postgresconn.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
