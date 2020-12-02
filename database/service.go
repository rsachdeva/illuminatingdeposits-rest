package database

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/router"
	"go.opencensus.io/trace"
)

// Service provides support for orchestration health checks.
type Service struct {
	Db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Health validates the router is healthy and ready to accept requests.
func (c *Service) Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "database.Service.Health")
	defer span.End()

	var health struct {
		Status string `json:"status"`
	}

	// Service if the database is ready.
	if err := StatusCheck(ctx, c.Db); err != nil {

		// If the database is not ready we will tell the cli and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "Db not ready"
		return router.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return router.Respond(ctx, w, health, http.StatusOK)
}
